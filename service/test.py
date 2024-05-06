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
#     # "products_ids": [],
#     # "payment_method": "CRYPTO",
#     # "delivery_method": "Hold The Gold",
#     # "total_price": 168333.84,
#     # "status_deliery":"",
# }

# data = requests.get("https://goldshop24.co/api/user/general_data", json=data, headers={
#     "Authorization": token
# })

# print(data.json())

# # print(data.json())
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
data = {
    'symbol': 'FX:XAUUSD',
    # 'timeframe': '1',
    # 'until': 1713360757,
    # 'to': 1713378757,
    # 'count': 1
}

# data = requests.post("http://127.0.0.1:5000/history",json=data)
data = requests.get("http://127.0.0.1:5000/market_status",json=data)


session_time = data.json()['session']
session_start, session_end = session_time.split('-')

# Get current time
current_time = datetime.datetime.now().time()

print(session_start)
session_start_hour, session_start_minute = map(int, session_start.split(':'))

# If current time is after session end, the market has closed for today
if current_time >= datetime.time(hour=int(session_end[:2]), minute=int(session_end[-2:])):
    # Calculate time until next session starts (assuming next session is tomorrow)
    next_session_start = datetime.datetime.combine(datetime.date.today() + datetime.timedelta(days=1),
                                                   datetime.time(hour=session_start_hour, minute=session_start_minute))
    time_until_next_session = next_session_start - datetime.datetime.now()
    print("Market is currently closed. It will open in:", time_until_next_session)
else:
    print("Market is currently open.")


# from bitpay.client import Client

# # Initialize BitPay client with your API token
# bitpay = Client.create_pos_client(api_token='YOUR_API_TOKEN')

# # Create a payment invoice
# invoice_data = {
#     "price": 100,  # Amount in Tether
#     "currency": "USDT",
#     # Additional options can be added here, like "redirectURL" for redirection after payment
# }
# invoice = bitpay.create_invoice(invoice_data)

# # Extract QR code URL from the invoice
# qr_code_url = invoice['data']['url']
# print("QR Code URL:", qr_code_url)
