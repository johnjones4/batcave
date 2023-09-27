from lib.displays.display import Display

class TerminalDisplay(Display):
    def reset(self):
        print("RESET", end='\r\n')

    def write(self, s: str):
        print(s, end='\r\n')

    def set_status_light(self, n: int, o: bool):
        print(f"Status Light {n}: {o}", end='\r\n')
