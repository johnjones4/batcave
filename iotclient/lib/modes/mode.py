from abc import ABC, abstractmethod

from lib.displays.display import Display

class Mode(ABC):
    def __init__(self, display: Display):
        self.display = display
        
    @abstractmethod
    def start(self):
        pass

    @abstractmethod
    def stop(self):
        pass

    def toggle(self):
        pass
