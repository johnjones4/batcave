import logging
from threading import Thread
import time
from lib.util import make_std_headers
from lib.workers.exception import WorkerStartupException
from lib.workers.worker import Worker
import websocket

class LogWorker(Worker):
    def __init__(self, hostname: str, secure_transport: bool, client_id: str, api_key: str):
        super().__init__()
        self.secure_transport = secure_transport
        self.hostname = hostname
        self.client_id = client_id
        self.api_key = api_key

    def worker(self):
        while True:
            try:
                msg = self.ws.recv()
                if self.queue:
                    self.queue.put(msg.strip())
                else:
                    time.sleep(1)
            except:
                logging.exception("Error receiving log messages")
                return #TODO reconnect?

    def start(self):
        self.ws = websocket.WebSocket()
        p = "s" if self.secure_transport else ""
        try:
            self.ws.connect(f"ws{p}://{self.hostname}/api/client/log",header= make_std_headers(self.client_id, self.api_key))
        except Exception as e:
            raise WorkerStartupException(e, "Error connecting to command endpoint")
        self.worker_thread = Thread(target=self.worker, name="log worker")
        self.worker_thread.daemon = True
        self.worker_thread.start()
