from typing import Callable
from lib.controllers.controller import MAX, Controller
import readkeys

from lib.runtime import Runtime

class KeyboardController(Controller):
    def wait_for_signal(self, runtime: Runtime):
        while True:
            s = readkeys.getch()
            if not s:
                continue
            if s == ' ':
                runtime.toggle_function()
            else:
                try:
                    i = int(s)
                    if 0 <= i < MAX:
                        runtime.switch_mode(i)
                except:
                    pass
