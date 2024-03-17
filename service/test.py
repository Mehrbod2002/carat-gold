import asyncio
import websockets

async def connect_to_websocket():
    uri = "wss://goldshop24.co/feed"
    async with websockets.connect(uri) as websocket:
        while True:
            try:
                message = await websocket.recv()
                print("Received message:", message)
            except websockets.exceptions.ConnectionClosed:
                print("Connection closed")
                break

asyncio.get_event_loop().run_until_complete(connect_to_websocket())
