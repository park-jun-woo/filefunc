class Server:
    def __init__(self, host, port):
        self.host = host
        self.port = port

    def start(self):
        print(f"Starting {self.host}:{self.port}")

    def stop(self):
        print("Stopping")
