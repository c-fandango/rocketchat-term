#!/usr/bin/python3

import asyncio
import pathlib
import websockets
import json
import time
import os

HOSTNAME = "localhost"
WSS_PORT = 5550
BIN_DIR = os.path.dirname(os.path.abspath(__file__))
DATA_DIR = os.path.join(BIN_DIR, '../data')
DATA_FILE = os.path.join(DATA_DIR, "ws_responses.json")


class WSResponse:
    def __init__(self, data_path):
        self.playbook = self.load_playbook(data_path)
        print("loaded data")

    def load_playbook(self, file_path):
        with open(file_path, "r") as fle:
            data = json.load(fle)

        return data

    async def serve(self, websocket):
        for response in self.playbook["init"]:
            data = await websocket.recv()

            message = json.loads(data)

            content = message.get("msg")

            if content == response["trigger_msg"]:
                send_msg = response["response"]

            send_msg["id"] = message.get("id")

            print(f"sending {send_msg}")
            await websocket.send(json.dumps(send_msg))

        for response in self.playbook["room_msgs"]:
            time.sleep(response["timing"])
            send_msg = response["response"]

            await websocket.send(json.dumps(send_msg))


async def main():
    ws_response = WSResponse(DATA_FILE)

    async with websockets.serve(ws_response.serve, HOSTNAME, WSS_PORT):
        await asyncio.Future()


if __name__ == "__main__":
    asyncio.run(main())
