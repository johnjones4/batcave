import logging
from queue import Queue
from threading import Thread
from lib.util import make_std_headers
from lib.workers.exception import WorkerStartupException
from lib.workers.worker import Worker
import base64
import websocket
from uuid import uuid4
import time
import json

class CommandWorker(Worker):
    def __init__(self, hostname: str, secure_transport: bool, client_id: str, api_key: str):
        super().__init__()
        self.secure_transport = secure_transport
        self.hostname = hostname
        self.client_id = client_id
        self.api_key = api_key
        self.send_queue = Queue()

    def send(self, audio):
        self.send_queue.put(audio)

    def send_worker(self):
        while True:
            while not self.send_queue.empty():
                try:
                    self.ws.send(json.dumps({
                        "clientId": self.client_id,
                        "eventId": str(uuid4()),
                        "message": {
                            "audio": {
                                "data": base64.b64encode(self.send_queue.get()).decode('ascii'),
                            }
                        }
                    }))
                except:
                    logging.exception("Error sending command")


    def receive_worker(self):
         while True:
            try:
                msg = self.ws.recv()
                if self.queue:
                    self.queue.put(json.loads(msg))
                else:
                    time.sleep(1)
            except Exception as e:
                logging.exception("Error receiving command response")
                if self.queue:
                    self.queue.put(e)
                return

    def start(self):
        self.ws = websocket.WebSocket()
        p = "s" if self.secure_transport else ""
        try:
            self.ws.connect(f"ws{p}://{self.hostname}/api/client/converse",header= make_std_headers(self.client_id, self.api_key))
        except Exception as e:
            raise WorkerStartupException(e, "Error connecting to command endpoint")
        self.send_worker_thread = Thread(target=self.send_worker, name="command send worker")
        self.send_worker_thread.daemon = True
        self.send_worker_thread.start()
        self.receive_worker_thread = Thread(target=self.receive_worker, name="command receive worker")
        self.receive_worker_thread.daemon = True
        self.receive_worker_thread.start()
