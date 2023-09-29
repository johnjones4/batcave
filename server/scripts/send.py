import json
import os
import sys
import uuid
import requests

ROOT = os.getenv('ROOT', 'http://localhost:8080')
CLIENT_ID = os.getenv('CLIENT_ID')
API_KEY = os.getenv('API_KEY')
EVENT_ID = str(uuid.uuid4())
MSG = sys.argv[1]

res = requests.post(f"{ROOT}/api/client/message", json={
  'eventId': EVENT_ID,
  'clientId': CLIENT_ID,
  'message': {
    'text': MSG
  }
}, headers={
  'X-Client-Id': CLIENT_ID,
  'X-Api-Key': API_KEY
})

print(json.dumps(res.json(), indent='  '))
