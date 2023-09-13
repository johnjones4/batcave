import json
import os
import openai
import pathlib

openai.api_key = os.getenv("OPENAI_API_KEY")
file = pathlib.Path(pathlib.Path(__file__).parent.resolve()).joinpath('intent_descriptions.json').resolve()
with open(file, 'r') as jsonfile:
    intent_map = json.loads(jsonfile.read())
    for intent in intent_map:
        emb = openai.Embedding.create(
            model="text-embedding-ada-002",
            input=intent_map[intent]
        )
        print(f"INSERT INTO intents (intent_label,embedding) VALUES ('{intent}','{json.dumps(emb.data[0].embedding)}');")
        