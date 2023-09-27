class WorkerStartupException(Exception):
    def __init__(self, root: Exception, msg: str):
        self.root = root
        self.msg = msg

    def __str__(self):
        return f"{self.msg}: {self.root}"
