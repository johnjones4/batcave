import json
import os
import requests
import pathlib

URL = os.getenv("OLLAMA_URL")
MODEL = os.getenv("OLLAMA_MODEL", "jsonllama2")
file = pathlib.Path(pathlib.Path(__file__).parent.resolve()).joinpath('intent_descriptions.json').resolve()
with open(file, 'r') as jsonfile:
    intent_map = json.loads(jsonfile.read())
    for intent in intent_map:
        res = requests.post(URL, json={
            "model": MODEL,
            "prompt": intent_map[intent]
        })
        print(f"INSERT INTO intents (intent_label,embedding) VALUES ('{intent}','{json.dumps(res.json()['embedding'])}');")
        