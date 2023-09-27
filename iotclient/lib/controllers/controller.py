from abc import ABC, abstractmethod
from typing import Callable

MAX = 2

class Controller(ABC):
    @abstractmethod
    def wait_for_signal(self, callback: Callable[[int], None]):
        pass
