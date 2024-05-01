import requests
import json
import websocket

url = 'https://goldshop24.co/history'
headers = {'Content-Type': 'application/json'}

data = {
    'symbol': 'FX:XAUUSD',
    'timeframe': '1',
    'until': 1714321363.0,
    'to': 1714387995.0,
    'count': 300
}

response = requests.post(url, headers=headers, data=json.dumps(data))

print(response.json())
# import asyncio
# import websockets
# import json

# async def connect_to_server():
#     async with websockets.connect(url) as websocket:
#         while True:
#             message = await websocket.recv()
#             data = json.loads(message)
#             if (data["symbol"] == "BTCUSD"):
#                 print(data)

# async def main():
#     await connect_to_server()

# if __name__ == "__main__":
#     asyncio.run(main())
