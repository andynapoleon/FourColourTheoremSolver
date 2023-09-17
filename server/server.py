from flask import Flask, jsonify
from flask_cors import CORS;    # allows interacting with other servers
import clingo
from pymongo.mongo_client import MongoClient
import skimage as ski
from skimage import io, segmentation, color
import matplotlib.pyplot as plt
import numpy as np
from scipy.signal import convolve2d
from skimage.morphology import binary_dilation, square

# app instance
app = Flask(__name__)
CORS(app)

# /api/home
@app.route("/api/home", methods=['GET'])
def return_home():
    return jsonify({
        "message": "Hello World!!!",
        "people": ["Sheikh", "Riley", "Peter", "Andy"]
    })


@app.route("/api/solve")
def solve():
    print("loading image")
    image = io.imread("testimg4.png")
    print("preprocessing image")
    image = preprocess_image(image)
    print("finding countries")
    vertices, black = get_vertices(image)
    print("finding neighbours")
    edges = find_edges(image, vertices)
    print("creating problem instance")
    program = generate_program(len(vertices), edges)
    print("selecting colors")
    solution = solve_graph(program)
    print("coloring map")
    color_map(vertices, solution, black)
    # use solution to color map
    # return colored map
    return "test"

def preprocess_image(image):
     # remove alpha channel
    if image.shape[2] == 4:
        image = color.rgba2rgb(image)
    # make greyscale
    image = color.rgb2gray(image)
    return image

def get_vertices(image):
    # find all uncolored chunks of the map
    vertices = []
    num = 0
    seed_point = (0,0)
    # find size of image
    height, width = image.shape[:2]
    for x in range(width):
        for y in range(height):
            if (image[y, x] == 1):
                # find the chunk associated with a vertex
                vertex = segmentation.flood(image, (y,x))
                vertices.append(vertex)
                # remove the chunk from the map
                image = segmentation.flood_fill(image, (y,x), 0)
    return vertices, image


def find_edges(image, vertices):
    # find all adjacecies between countries
    edges = list()
    num_vertices = len(vertices)
    # fuzzyness is how far countries are allowed to be apart to still be considered as bordering each other
    fuzzyness = 10
    # find neighbours for each country
    for i in range(num_vertices):
        # expand countries size to check overlapping
        dilated_image = binary_dilation(vertices[i], footprint=square(fuzzyness))
        # check each possible other country
        for j in range((i + 1), num_vertices):
            # take only the overlap between the enlarged country and its neighbour
            overlap = np.minimum(dilated_image, vertices[j])
            # check if there is any overlap
            all_zeros = np.all(overlap == 0)
            # if adjacent add an edge
            if  (all_zeros == False):
                edges.append((i, j))

    return edges

def generate_program(num_vertices, edges):
    program = ""
    for vertex in range(num_vertices):
        program += "vertex(" + str(vertex) +")."
    for edge in edges:
        program += "edge(" + str(edge[0]) + "," + str(edge[1]) + ")."
    return program

def solve_graph(graph):
    with open('server/asp program/program.lp', 'r') as file:
        program = file.read()
    with open('server/asp program/colors.lp', 'r') as file:
        colors = file.read()

    ctl = clingo.Control()
    ctl.add("pro", [], program + colors + graph)
    ctl.ground([("pro",[])])
    ctl.configuration.solve.models="1"      # max number of models to calculate, 0 for all
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
    return(graph)

def color_map(vertices, solution, black):
    image = black
    image = color.gray2rgb(image)
    for i in range(len(vertices)):
        mask = vertices[i]
        vertices[i] = color.gray2rgb(vertices[i])
        colored = solution[str(i)]
        if (colored == "green"):
            new_color = (0,255,0)
        elif (colored == "blue"):
            new_color = (0,0,255)
        elif (colored == "red"):
            new_color = (255,0,0)
        else:
            new_color = (255,255,0)
        vertices[i][mask] = new_color
        image = np.maximum(image, vertices[i])
    io.imsave('image.png', image)

if __name__ == "__main__":
    app.run(debug=True, port=8080)