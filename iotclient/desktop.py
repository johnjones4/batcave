import logging
import os
from lib import start
from lib.controllers.keyboard_controller import KeyboardController
from lib.displays.terminal_display import TerminalDisplay

logging.basicConfig(level=logging.DEBUG)
d = TerminalDisplay()
c = KeyboardController()
start(d, c)
