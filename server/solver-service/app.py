from flask import Flask, jsonify, request
from flask_cors import CORS
from dotenv import load_dotenv
load_dotenv()
import os
import clingo
import skimage as ski
from skimage import io, segmentation, color
import matplotlib.pyplot as plt
import numpy as np
from skimage.morphology import binary_dilation, square
import time
import json
import grpc
from google.protobuf.timestamp_pb2 import Timestamp
from datetime import datetime
import proto.logs.logger_pb2 as logger_pb2
import proto.logs.logger_pb2_grpc as logger_pb2_grpc

# app instance
app = Flask(__name__)
CORS(app)

# Add logging configuration
LOGGER_URL = "logger-service:50001"  # gRPC service address

def log_event(user_id, event_type, description, severity=1, metadata=None):
    try:
        # Create gRPC channel
        channel = grpc.insecure_channel(LOGGER_URL)
        stub = logger_pb2_grpc.LoggerServiceStub(channel)

        # Create timestamp
        timestamp = datetime.utcnow().strftime('%Y-%m-%dT%H:%M:%SZ')

        # Create log request
        request = logger_pb2.LogRequest(
            service_name="map_coloring_service",
            event_type=event_type,
            user_id=user_id,
            description=description,
            severity=severity,
            timestamp=timestamp,
            metadata=metadata or {}
        )

        # Make gRPC call
        response = stub.LogEvent(request)
        print(f"Log event response: {response}")
        
        if not response.success:
            print(f"Failed to log event: {response.message}")
            
    except Exception as e:
        print(f"Error logging event: {str(e)}")


@app.route("/", methods=["GET"])
def return_home():
    return jsonify(
        {"message": "Hello World!!!", "people": ["Sheikh", "Riley", "Peter", "Andy"]}
    )


@app.route('/api/solve', methods=['POST'])
def solve():
    try:
        data = request.get_json()
        if not data:
            log_event(
                "unknown",
                "map_coloring_failed",
                "No JSON data received",
                2,
                {"error": "Missing request data"}
            )
            return jsonify({"error": "No JSON data received"}), 400

        user_id = data.get("userId", "unknown")
        log_event(
            user_id,
            "map_coloring_started",
            "Starting map coloring process",
            1,
            {"width": str(data.get("width")), "height": str(data.get("height"))}
        )

        print("Received request data:", data.keys())
        
        if 'image' not in data or 'width' not in data or 'height' not in data:
            log_event(
                user_id,
                "map_coloring_failed",
                "Missing required fields",
                2,
                {"missing_fields": str([f for f in ['image', 'width', 'height'] if f not in data])}
            )
            return jsonify({"error": "Missing required fields"}), 400

        # Convert image data to integers if they're strings
        image_data = data["image"]
        if isinstance(image_data[0], str):
            image_data = [int(x) for x in image_data]
        
        width = int(data["width"])
        height = int(data["height"])
        
        # Convert image data to numpy array
        index = 0
        array = np.zeros((height, width))

        for y in range(height):
            for x in range(width):
                try:
                    pixel_value = int(image_data[index])
                    array[y][x] = 1 if pixel_value > 128 else 0
                except (ValueError, TypeError) as e:
                    log_event(
                        user_id,
                        "map_coloring_failed",
                        f"Invalid pixel data at index {index}",
                        2,
                        {"error": str(e), "index": str(index)}
                    )
                    return jsonify({"error": f"Invalid pixel data at index {index}"}), 400
                index += 4

        begin = time.time()
        
        log_event(
            user_id,
            "map_processing_step",
            "Finding countries",
            1,
            {"step": "get_vertices"}
        )
        vertices, black, vertice_matrix = get_vertices(array)
        
        log_event(
            user_id,
            "map_processing_step",
            "Finding neighbours",
            1,
            {"step": "find_edges"}
        )
        edges = find_edges(array, vertices, vertice_matrix)
        
        log_event(
            user_id,
            "map_processing_step",
            "Creating problem instance",
            1,
            {"step": "generate_program", "vertices": str(len(vertices)), "edges": str(len(edges))}
        )
        program = generate_program(len(vertices), edges)
        
        log_event(
            user_id,
            "map_processing_step",
            "Selecting colors",
            1,
            {"step": "solve_graph"}
        )
        solution = solve_graph(program)
        
        log_event(
            user_id,
            "map_processing_step",
            "Coloring map",
            1,
            {"step": "color_map"}
        )
        colored_map = color_map(vertices, solution, black)
        
        end = time.time()
        processing_time = end - begin

        log_event(
            user_id,
            "map_coloring_completed",
            "Successfully colored map",
            1,
            {
                "processing_time": f"{processing_time:.2f}",
                "vertices": str(len(vertices)),
                "edges": str(len(edges))
            }
        )

        # Convert numpy array to list for JSON serialization
        result = colored_map.tolist()
        return jsonify(result)

    except Exception as e:
        user_id = data.get("userId", "unknown") if 'data' in locals() else "unknown"
        log_event(
            user_id,
            "map_coloring_failed",
            "Unexpected error during map coloring",
            3,
            {"error": str(e)}
        )
        print(f"Error processing request: {str(e)}")
        import traceback
        traceback.print_exc()
        return jsonify({"error": str(e)}), 500


# Add a health check endpoint
@app.route('/health', methods=['GET'])
def health():
    return jsonify({"status": "healthy"})


def preprocess_image(image):
    # remove alpha channel
    # if image.shape[2] == 4:
    #   image = color.rgba2rgb(image)
    # make greyscale
    image = color.rgb2gray(image)
    return image


def get_vertices(image):
    # find all uncolored chunks of the map
    vertices = []
    num = 0
    seed_point = (0, 0)
    # find size of image
    height, width = image.shape[:2]
    vertice_matrix = image
    for x in range(width):
        for y in range(height):
            if image[y, x] == 1:
                num += 1
                # find the chunk associated with a vertex
                vertex = segmentation.flood(image, (y, x))
                vertice_matrix[vertex] = num
                vertices.append(vertex)
                # remove the chunk from the map
                image = segmentation.flood_fill(image, (y, x), 0)

    return vertices, image, vertice_matrix


def find_edges(image, vertices, vertice_matrix):
    # find all adjacecies between countries
    edges = list()
    num_vertices = len(vertices)
    # fuzzyness is how far countries are allowed to be apart to still be considered as bordering each other
    fuzzyness = 8
    # find neighbours for each country
    start = time.time()
    for i in range(num_vertices):
        # expand countries size to check overlapping
        dilated_image = binary_dilation(vertices[i], footprint=square(fuzzyness))
        height, width = dilated_image.shape[:2]

        vertice_matrix_copy = vertice_matrix.copy()
        vertice_matrix_copy[np.logical_not(dilated_image)] = 0
        adjacents = np.unique(vertice_matrix_copy)
        for adjacent in adjacents:
            if adjacent != 0 and adjacent != (i + 1):
                edges.append((i, int(adjacent - 1)))

        """
        # check each possible other country
        for j in range((i + 1), num_vertices):
            # take only the overlap between the enlarged country and its neighbour
            overlap = np.minimum(dilated_image, vertices[j])
            # check if there is any overlap
            all_zeros = np.all(overlap == 0)
            # if adjacent add an edge
            if  (all_zeros == False):
                edges.append((i, j))
        """

    end = time.time()

    print(end - start)
    return edges


def generate_program(num_vertices, edges):
    program = ""
    for vertex in range(num_vertices):
        program += "vertex(" + str(vertex) + ")."
    for edge in edges:
        program += "edge(" + str(edge[0]) + "," + str(edge[1]) + ")."
    return program


def solve_graph(graph):
    with open("./asp_program/program.lp", "r") as file:
        program = file.read()
    with open("./asp_program/colors.lp", "r") as file:
        colors = file.read()

    ctl = clingo.Control()
    ctl.add("pro", [], program + colors + graph)
    ctl.ground([("pro", [])])
    ctl.configuration.solve.models = "1"  # max number of models to calculate, 0 for all
    models = []
    with ctl.solve(yield_=True) as handle:
        for model in handle:
            # select all atoms which would be shown in program output
            models.append(model.symbols(shown=True))
    model = models[0]
    list_models = list()
    graph = dict()
    for atom in model:
        vertex = str(atom.arguments[0])
        color = str(atom.arguments[1])
        graph[vertex] = color
    return graph


def color_map(vertices, solution, black):
    image = black
    image = color.gray2rgb(image)
    for i in range(len(vertices)):
        mask = vertices[i]
        vertices[i] = color.gray2rgb(vertices[i])
        colored = solution[str(i)]
        if colored == "green":
            new_color = (0, 255, 0)
        elif colored == "blue":
            new_color = (0, 0, 255)
        elif colored == "red":
            new_color = (255, 0, 0)
        else:
            new_color = (255, 255, 0)
        vertices[i][mask] = new_color
        image = np.maximum(image, vertices[i])
    # io.imsave("image.png", image)
    return image


if __name__ == "__main__":
    port = os.getenv("PORT")
    app.run(port=(port or 1000))
