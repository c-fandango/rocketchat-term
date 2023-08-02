#!/usr/bin/python3

from http.server import BaseHTTPRequestHandler, HTTPServer
import os


HOSTNAME = "localhost"
PORT = 5555
BIN_DIR = os.path.dirname(os.path.abspath(__file__))
DATA_DIR = os.path.join(BIN_DIR, '../data')
DATA_FILE = os.path.join(DATA_DIR, "http_response.json")


class ServeRooms(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/api/v1/rooms.get":
            self.send_response(200)
            self.send_header("Content-type", "application/json")
            self.end_headers()
            self.wfile.write(self.load_responseb(DATA_FILE))
        else:
            self.send_response(404)
            self.end_headers()

    def load_responseb(self, file_path):
        with open(file_path, "r") as fle:
            data = fle.read()

        return data.encode()


if __name__ == "__main__":
    webServer = HTTPServer((HOSTNAME, PORT), ServeRooms)

    try:
        webServer.serve_forever()
    except KeyboardInterrupt:
        pass

    webServer.server_close()
    print("Server stopped.")
