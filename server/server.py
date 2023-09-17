from flask import Flask, jsonify
from flask_cors import CORS;    # allows interacting with other servers
from pymongo.mongo_client import MongoClient
from dotenv import load_dotenv
load_dotenv()
import os

# app instance
app = Flask(__name__)
CORS(app)

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

if __name__ == "__main__":
    app.run(debug=True, port=8080)