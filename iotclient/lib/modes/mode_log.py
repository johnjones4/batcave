import logging
from queue import Queue
from threading import Thread
from lib.displays.display import Display
from lib.modes.mode import Mode
from lib.workers.log_worker import LogWorker

class ModeLog(Mode):
    def __init__(self, display: Display, log_worker: LogWorker):
        super().__init__(display)
        self.log_worker = log_worker

    def start(self):
        logging.debug("Starting log mode")
        self.display.reset()
        self.running = True
        self.log_worker.play()
        self.worker_thread = Thread(target=self.worker, name="log mode worker")
        self.worker_thread.daemon = True
        self.worker_thread.start()

    def worker(self):
        while self.running:
            while self.running and self.log_worker.queue and not self.log_worker.queue.empty():
                self.display.write(self.log_worker.queue.get())

    def stop(self):
        logging.debug("Stopping log mode")
        self.log_worker.pause()
        self.running = False
        self.worker_thread = None
