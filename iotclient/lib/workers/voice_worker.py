import logging
from threading import Thread
from queue import Queue
import time
import pyaudio
from lib.workers.worker import Worker
import webrtcvad
import wave
from io import BytesIO

FORMAT = pyaudio.paInt16
RATE_PROCESS = 16000
CHANNELS = 1
BLOCKS_PER_SECOND = 50
BLOCK_SIZE = int(RATE_PROCESS / float(BLOCKS_PER_SECOND))

class VoiceWorker(Worker):
    def __init__(self, vox: bool) -> None:
        super().__init__()
        self.intermediate_queue = Queue()
        self.pa = pyaudio.PyAudio()
        self.stream = self.pa.open(
            format=FORMAT,
            channels=CHANNELS,
            rate=RATE_PROCESS,
            input=True,
            frames_per_buffer=BLOCK_SIZE,
            stream_callback=self.pyaudio_callback
        )
        self.vox = vox
        if vox:
            self.vad = webrtcvad.Vad(3)
            self.is_purging = False

    def purge(self):
        self.is_purging = True
        if not self.vox:
            # TODO bugs here
            buffer = bytearray()
            while not self.intermediate_queue.empty():
                frame = self.intermediate_queue.get()
                buffer.extend(frame)
            b = buffer.copy()
            if len(b) < 16000:
                return
            bytesio = self.make_wav(b)
            if self.queue:
                self.queue.put(bytesio.getbuffer())
      
    def start(self):
        self.is_purging = False
        self.stream.start_stream()
        if self.vox:
            self.vad_thread = Thread(target=self.vad_worker, name="vad worker")
            self.vad_thread.daemon = True
            self.vad_thread.start()

    def pyaudio_callback(self, in_data, frame_count, time_info, status):
        if not self.queue or self.is_purging:
            return (None, pyaudio.paContinue)
        self.intermediate_queue.put(in_data)
        return (None, pyaudio.paContinue)
    
    def is_voice(self, frame):
        return self.vad.is_speech(frame, RATE_PROCESS)
    
    def make_wav(self, b: bytearray):
        bytesio = BytesIO()
        wf = wave.open(bytesio, 'wb')
        wf.setnchannels(CHANNELS)
        wf.setsampwidth(2)
        wf.setframerate(RATE_PROCESS)
        wf.writeframes(b)
        wf.close()
        return bytesio
    
    def vad_worker(self):
        buffer = bytearray()
        last_frame = time.time()
        activated = False
        i = 0
        while True:
            if not self.queue:
                buffer = bytearray()
            while not self.intermediate_queue.empty():
                frame = self.intermediate_queue.get()
                is_v_frame = len(frame) >= 640 and self.is_voice(frame)
                if not activated:
                    activated = is_v_frame
                    if activated:
                        logging.debug("VOX activated")
                if not is_v_frame and time.time() - last_frame > 2:
                    break
                if is_v_frame:
                    last_frame = time.time()
                if activated:
                    buffer.extend(frame)
                    
            if self.is_purging or (len(buffer) > 16000 and time.time() - last_frame > 2):
                logging.debug("VOX complete")
                b = buffer.copy()
                bytesio = self.make_wav(b)
                if self.queue:
                    self.queue.put(bytesio.getbuffer())
                buffer = bytearray()
                last_frame = time.time()
                activated = False
            if self.is_purging:
                return
