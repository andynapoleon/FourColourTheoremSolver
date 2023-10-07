from flask import Flask, jsonify, request
from flask_cors import CORS;    # allows interacting with other servers
from pymongo.mongo_client import MongoClient
import os
import clingo
from pymongo.mongo_client import MongoClient
import skimage as ski
from skimage import io, segmentation, color
import matplotlib.pyplot as plt
import numpy as np
from scipy.signal import convolve2d
from skimage.morphology import binary_dilation, square
import time
import json

# app instance
app = Flask(__name__)
CORS(app)
   
'''
uri = os.environ.get("URI")
# Create a new client and connect to the server
client = MongoClient(uri)
# Send a ping to confirm a successful connection
try:
    client.admin.command('ping')
    print("Pinged your deployment. You successfully connected to MongoDB!")
except Exception as e:
    print(e)

@app.route("/api/database")
def read():
    db = client["GraphDB"]
    collection = db["Nodes"]

    # Document to insert
    document1 = {"name": "Alice", "email": "examplealice@example.com", "vertices": "[]", "edges": "[]"}
    document2 = {"name": "Bob", "email": "examplebob@example.com", "vertices": "[]", "edges": "[]"}

    # Insert the documents into the collection
    # Use a for loop to insert the list of documents if there are multiple
    collection.insert_one(document1)
    collection.insert_one(document2)

    # Delete a document
    delete_filter = {"name": "Alice"}
    # Delete the document that matches the filter
    result = collection.delete_one(delete_filter)
    print("Deleted count:", result.deleted_count)
    
    documents = collection.find()
    for document in documents:
        print(document)
    
    return return_home()

 
# /api/home
@app.route("/api/home", methods=['GET'])
def return_home():
    return jsonify({
        "message": "Hello World!!!",
        "people": ["Sheikh", "Riley", "Peter", "Andy"]
    })
'''

@app.route("/api/solve", methods=['POST'])
def solve():
    image = list(request.json["image"].values())
    width = request.json["width"]
    height = request.json["height"]
    index = 0
    array = np.zeros((height, width))
    
    for y in range(height):
        for x in range(width):
            if (image[index] > 128):
                array[y][x] = 1
            else:
                array[y][x] = 0
            index += 4
    image = array
    begin =time.time()
    print("loading image")
    print("preprocessing image")
    #image = preprocess_image(image)
    print("finding countries")
    vertices, black, vertice_matrix = get_vertices(image)
    print("finding neighbours")
    edges = find_edges(image, vertices, vertice_matrix)
    print("creating problem instance")
    program = generate_program(len(vertices), edges)
    print("selecting colors")
    solution = solve_graph(program)
    print("coloring map")
    colored_map = color_map(vertices, solution, black)
    end = time.time()
    print(end -begin)
    array = colored_map.tolist()
    json_img = jsonify(array)
    return json_img


def preprocess_image(image):
     # remove alpha channel
    #if image.shape[2] == 4:
     #   image = color.rgba2rgb(image)
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
    vertice_matrix = image
    for x in range(width):
        for y in range(height):
            if (image[y, x] == 1):
                num += 1
                # find the chunk associated with a vertex
                vertex = segmentation.flood(image, (y,x))
                vertice_matrix[vertex] = (num)
                vertices.append(vertex)
                # remove the chunk from the map
                image = segmentation.flood_fill(image, (y,x), 0)

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
            if (adjacent != 0 and adjacent != (i+1)):
                edges.append((i, int(adjacent-1)))

        '''
        # check each possible other country
        for j in range((i + 1), num_vertices):
            # take only the overlap between the enlarged country and its neighbour
            overlap = np.minimum(dilated_image, vertices[j])
            # check if there is any overlap
            all_zeros = np.all(overlap == 0)
            # if adjacent add an edge
            if  (all_zeros == False):
                edges.append((i, j))
        '''
        
        
    end = time.time()

    print(end - start)
    return edges

def generate_program(num_vertices, edges):
    program = ""
    for vertex in range(num_vertices):
        program += "vertex(" + str(vertex) +")."
    for edge in edges:
        program += "edge(" + str(edge[0]) + "," + str(edge[1]) + ")."
    return program

def solve_graph(graph):
    with open('asp program/program.lp', 'r') as file:
        program = file.read()
    with open('asp program/colors.lp', 'r') as file:
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
    return image

if __name__ == "__main__":
    app.run(debug=True, port=8080)