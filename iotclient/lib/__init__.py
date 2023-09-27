import logging
import os
from lib.controllers.controller import Controller
from lib.displays.display import Display
from lib.modes.mode_command import ModeCommand
from lib.modes.mode_log import ModeLog
from lib.runtime import Runtime
from lib.workers.command_worker import CommandWorker
from lib.workers.log_worker import LogWorker
from lib.workers.voice_worker import VoiceWorker

def start(d: Display, c: Controller):
    try:
        hostname = os.getenv("HOSTNAME")
        secure_transport = True if os.getenv("SECURE_TRANSPORT") else False
        client_id = os.getenv("CLIENT_ID")
        api_key = os.getenv("API_KEY")
        voice_worker = VoiceWorker(False)
        log_worker = LogWorker(hostname, secure_transport, client_id, api_key)
        command_worker = CommandWorker(hostname, secure_transport, client_id, api_key)
        modes = [
            ModeCommand(d, voice_worker, command_worker),
            ModeLog(d, log_worker),
        ]
        rt = Runtime(modes, [voice_worker, log_worker, command_worker], c)
        rt.start()
    except:
        logging.exception("Exception during startup")
