from flask import Flask, jsonify
from flask_cors import CORS;    # allows interacting with other servers

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

if __name__ == "__main__":
    app.run(debug=True, port=8080)