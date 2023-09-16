from flask import Flask, jsonify
from flask_cors import CORS;    # allows interacting with other servers
from dotenv import load_dotenv
import os

load_dotenv()

from pymongo.mongo_client import MongoClient

uri = os.getenv("URI")
# Create a new client and connect to the server
client = MongoClient(uri)

# Send a ping to confirm a successful connection
try:
    client.admin.command('ping')
    print("Pinged your deployment. You successfully connected to MongoDB!")
except Exception as e:
    print(e)

db = client['GraphDB']
collection = db['Nodes']

for document in collection.find():
    print(document)

client.close()



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