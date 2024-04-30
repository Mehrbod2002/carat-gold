import asyncio
import websockets
import requests
import datetime

token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI2NjI5NDRmMWMwODE4ZWMwYmI5ZWJhYzciLCJuYW1lIjoiIiwicGhvbmUiOiI5ODkxMzg3ODAyNzUiLCJlbWFpbCI6IiIsImNyZWF0ZWRfYXQiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiIsImV4cCI6MTcxNjU3MjcwNH0.Ej1devNLhmZCH8nMASnv5WutbDOEyf8-oFO_pkEyihE"

# edit_user = requests.post("http://127.0.0.1:3000/api/user/edit_user", json={
#     "name": "Mehrbod man",
#     "email": "m9.akhlaghpoor@gmail.com",
#     "phone": "989138780275",
#     "address": [
#         {"address": "King bernam", "label": "my home"},
#         {"address": "my queen nam bernam", "label": "my work"},
#     ]
# }, headers={
#     "Authorization": token
# }).json()

# data = {
#     "products_ids": [],
#     "payment_method": "CRYPTO",
#     "delivery_method": "Hold The Gold",
#     "total_price": 168333.84,
#     "status_deliery":"",
# }

# data = requests.post("http://127.0.0.1:3000/api/user/create_transaction", json=data, headers={
#     "Authorization": token
# })

# print(data.json())

# print(data.json())
# async def connect_to_websocket():
#     uri = "wss://goldshop24.co/feed"
#     async with websockets.connect(uri) as websocket:
#         while True:
#             try:
#                 message = await websocket.recv()
#                 print("Received message:", message)
#             except websockets.exceptions.ConnectionClosed:
#                 print("Connection closed")
#                 break

# # asyncio.get_event_loop().run_until_complete(connect_to_websocket())
# data = {
#     'symbol': 'FX:XAUUSD',
#     'timeframe': '1',
#     'until': 1713360757,
#     'to': 1713378757,
#     'count': 1
# }

# data = requests.post("http://127.0.0.1:3000/history",json=data)

# print(data.json())

import requests

url = "https://test.bitpay.com/tokens"

payload = { "facade": "pos" }
headers = {
    "accept": "application/json",
    "Content-Type": "application/json",
    "X-Accept-Version": "2.0.0"
}

response = requests.post(url, json=payload, headers=headers)

print(response.text)