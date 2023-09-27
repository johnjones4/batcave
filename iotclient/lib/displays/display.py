from abc import ABC, abstractmethod
from typing import Callable

class Display(ABC):
    @abstractmethod
    def reset(self):
        pass

    @abstractmethod
    def write(self, s: str):
        pass

    @abstractmethod
    def set_status_light(self, n: int, o: bool):
        pass
