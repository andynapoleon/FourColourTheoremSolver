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

if __name__ == "__main__":
    app.run(debug=True, port=8080)