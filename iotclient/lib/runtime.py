from lib.modes.mode import Mode
from lib.controllers.controller import Controller
from lib.workers.worker import Worker

class Runtime:
    def __init__(self, modes: list[Mode], workers: list[Worker], controller: Controller):
        self.modes = modes
        self.workers = workers
        self.controller = controller
        self.current_mode = None

    def start(self):
        for worker in self.workers:
            worker.start()
        self.switch_mode(0)
        self.controller.wait_for_signal(self)

    def switch_mode(self, new_mode: int):
        if self.current_mode:
            self.current_mode.stop()
        self.current_mode = self.modes[new_mode]
        self.current_mode.start()

    def toggle_function(self):
        if self.current_mode:
            self.current_mode.toggle()
