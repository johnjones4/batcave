from abc import ABC, abstractmethod
from queue import Queue

class Worker(ABC):
    def __init__(self):
        self.queue = None

    def is_paused(self):
        return self.queue == None

    def play(self):
        self.queue = Queue()

    def pause(self):
        self.queue = None
    
    @abstractmethod
    def start(self):
        pass

