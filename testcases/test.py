import asyncio
import websockets

async def connect_to_websocket():
    url = "wss://feed.caratgold/feed"
    # url = "http://127.0.0.1:3903/feed"
    async with websockets.connect(url) as websocket:
        response = await websocket.recv()
        print(f"Received: {response}")

if __name__ == "__main__":
    asyncio.get_event_loop().run_until_complete(connect_to_websocket())
