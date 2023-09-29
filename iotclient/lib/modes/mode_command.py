import logging
from threading import Thread
import time
from lib.displays.display import Display
from lib.modes.mode import Mode
from lib.workers.command_worker import CommandWorker
from lib.workers.voice_worker import VoiceWorker

class ModeCommand(Mode):
    def __init__(self, display: Display, voice_worker: VoiceWorker, command_worker: CommandWorker):
        super().__init__(display)
        self.voice_worker = voice_worker
        self.command_worker = command_worker

    def start(self):
        logging.debug("Starting comand mode")
        self.display.reset()
        self.command_worker.play()
        if self.voice_worker.vox:
            self.voice_worker.play()
        self.running = True
        self.worker_thread = Thread(target=self.worker, name="command mode worker")
        self.worker_thread.daemon = True
        self.worker_thread.start()

    def stop(self):
        logging.debug("Stopping comand mode")
        self.running = False
        self.worker_thread = None
        self.voice_worker.pause()
        self.command_worker.pause()

    def toggle(self):
        if self.voice_worker.is_paused():
            logging.debug("Listening toggle on")
            self.display.set_status_light(0, True)
            self.voice_worker.play()
        else:
            logging.debug("Listening toggle off")
            self.display.set_status_light(0, False)
            self.voice_worker.purge()
            time.sleep(1)
            self.voice_worker.pause()

    def worker(self):
        while self.running:
            if self.voice_worker.queue and not self.voice_worker.queue.empty():
                logging.debug("Sending command audio")
                self.display.set_status_light(1, True)
                self.display.set_status_light(2, False)
                audio = self.voice_worker.queue.get()
                self.command_worker.send(audio)
            if self.command_worker.queue and not self.command_worker.queue.empty():
                logging.debug("Receiving command response data")
                incoming = self.command_worker.queue.get()
                if isinstance(incoming, Exception):
                    self.display.set_status_light(1, False)
                    self.display.set_status_light(2, True)
                else:
                    if incoming['type'] == "request":
                        self.display.write(f"You: {incoming['request']['message']['text']}")
                    elif incoming['type'] == "response":
                        self.display.write(f"HAL 9000: {incoming['response']['message']['text']}")
                        # TODO images?
                        self.display.set_status_light(1, False)
                        self.display.set_status_light(2, False)
                    elif incoming['type'] == "push":
                        self.display.write(f"HAL 9000: {incoming['push']['message']['text']}")
                        # TODO images?
                        self.display.set_status_light(1, False)
                        self.display.set_status_light(2, False)
