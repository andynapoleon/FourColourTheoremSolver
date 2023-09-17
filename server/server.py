from flask import Flask, jsonify
from flask_cors import CORS;    # allows interacting with other servers

from pymongo.mongo_client import MongoClient

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

@app.route("/api/home/models")
def solve_graph():
    import clingo
    with open('server/asp program/program.lp', 'r') as file:
        program = file.read()
    with open('server/asp program/colors.lp', 'r') as file:
        colors = file.read()
    with open('server/asp program/graph.lp', 'r') as file:
        graph = file.read()

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

    return(jsonify(graph))

@app.route("/api/analyze")
def analyze_map():
    import skimage as ski
    from skimage import io, segmentation, color
    import matplotlib.pyplot as plt
    import numpy as np
    from scipy.signal import convolve2d
    from skimage.morphology import binary_dilation, square

    # this needs to be an argument now
    image = io.imread("img01.png")

    # remove alpha channel
    if image.shape[2] == 4:
        image = color.rgba2rgb(image)
    # make greyscale
    image = color.rgb2gray(image)

    # find all uncolored chunks of the map
    vertexs = []
    num = 0
    seed_point = (0,0)
    # find size of image
    height, width = image.shape[:2]
    for x in range(width):
        for y in range(height):
            if (image[y, x] == 1):
                # find the chunk associated with a vertex
                vertex = segmentation.flood(image, (y,x))
                vertexs.append(vertex)
                # remove the chunk from the map
                image = segmentation.flood_fill(image, (y,x), 0)

    # find all adjacecies between countries
    edges = list()
    num_vertexs = len(vertexs)
    # fuzzyness is how far countries are allowed to be apart to still be considered as bordering each other
    fuzzyness = 10
    # find neighbours for each country
    for i in range(num_vertexs):
        # expand countries size to check overlapping
        dilated_image = binary_dilation(vertexs[i], footprint=square(fuzzyness))
        # check each possible other country
        for j in range((i + 1), num_vertexs):
            # take only the overlap between the enlarged country and its neighbour
            overlap = np.minimum(dilated_image, vertexs[j])
            # check if there is any overlap
            all_zeros = np.all(overlap == 0)
            # if adjacent add an edge
            if  (all_zeros == False):
                edges.append((i, j))

    return("placeholder")
        


if __name__ == "__main__":
    app.run(debug=True, port=8080)