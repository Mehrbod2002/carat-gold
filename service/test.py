import asyncio
import websockets
import requests

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

# asyncio.get_event_loop().run_until_complete(connect_to_websocket())
data = {
    'symbol': 'FX:XAUUSD',
    'timeframe': '1',
    'until': 1713360757,
    'to': 1713378757,
    'count': 1
}

data = requests.post("http://127.0.0.1:3000/api/public/get_history",json=data)

print(data.json())