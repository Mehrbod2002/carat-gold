# import asyncio
# import websockets
import requests
# import datetime

# token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI2NjUxYjZlZmJhODA5NTkwNjIxMzkyMjYiLCJuYW1lIjoiIiwicGhvbmUiOiI5ODkxMzg3ODAyNzUiLCJlbWFpbCI6IiIsImNyZWF0ZWRfYXQiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiIsImV4cCI6MTcxOTIyMzQ2NX0.R_rJU7MhtVN9PvUfearTa8Ke2lGTa8S9B6b2KffyeVE"
# # edit_user = requests.post("https://goldshop24.co/api/auth/user/send_otp", json={
# #     "phone": "989138780275",
# # }, headers={
# #     "Authorization": token
# # })

# # edit_user = requests.post("http://127.0.0.1:3000/api/auth/user/register", json={
# #     "otp_code": 12345,
# #     "phone": "989138780275",
# # }, headers={
# #     "Authorization": token
# # })
# # print(edit_user)
# # print(edit_user.json())


# me = requests.get("http://127.0.0.1:3000/api/user/me", json={
#     "first_name": "Mehrbod man",
#     "last_name":"test",
#     "phone":"989138780275",
#     "email":"m9.akhlaghpour@gmail.com",
#     "address": [
#         {"address": "King bernam","label": "my home","city":"test","country":"iran","region":"T"},
#         {"address": "my queen nam bernam", "label": "my work","city":"test","country":"iran","region":"A"},
#     ]
# }, headers={
#     "Authorization": token
# })
# print(me)
# print(me.json())

# # edit_user = requests.post("http://127.0.0.1:3000/api/user/edit_user", json={
# #     "name": "Mehrbod man",
# #     "email": "m9.akhlaghpoor@gmail.com",
# #     "phone": "989138780275",
# #     "address": [
# #         {"address": "King bernam", "label": "my home"},
# #         {"address": "my queen nam bernam", "label": "my work"},
# #     ]
# # }, headers={
# #     "Authorization": token
# # }).json()

# # data = {
# #     # "products_ids": [],
# #     # "payment_method": "CRYPTO",
# #     # "delivery_method": "Hold The Gold",
# #     # "total_price": 168333.84,
# #     # "status_deliery":"",
# # }

# # data = requests.get("https://goldshop24.co/api/user/general_data", json=data, headers={
# #     "Authorization": token
# # })

# # print(data.json())

# # # print(data.json())
# # async def connect_to_websocket():
# #     uri = "wss://goldshop24.co/feed"
# #     async with websockets.connect(uri) as websocket:
# #         while True:
# #             try:
# #                 message = await websocket.recv()
# #                 print("Received message:", message)
# #             except websockets.exceptions.ConnectionClosed:
# #                 print("Connection closed")
# #                 break

# # # asyncio.get_event_loop().run_until_complete(connect_to_websocket())
# # data = {
# #     'symbol': 'FX:XAUUSD',
# #     # 'timeframe': '1',
# #     # 'until': 1713360757,
# #     # 'to': 1713378757,
# #     # 'count': 1
# # }

# # # data = requests.post("http://127.0.0.1:5000/history",json=data)
data = requests.get("https://goldshop24.co/market_status")
print(data.json())
# # print(data.json())
# # session_time = data.json()['session']
# # session_start, session_end = session_time.split('-')

# # current_time = datetime.datetime.now().time()

# # print(session_start,session_end)
# # session_start_utc = datetime.datetime.strptime(session_start, '%H%M').replace(tzinfo=datetime.timezone.utc)
# # session_end_utc = datetime.datetime.strptime(session_end, '%H%M').replace(tzinfo=datetime.timezone.utc)

# # current_time_utc = datetime.datetime.utcnow()

# # print(current_time_utc,session_end_utc)
# # if current_time_utc >= session_end_utc:
# #     # Calculate time until next session starts (assuming next session is tomorrow)
# #     next_session_start_utc = datetime.datetime.combine(datetime.date.today() + datetime.timedelta(days=1),
# #                                                        session_start_utc.time())
# #     time_until_next_session = next_session_start_utc - current_time_utc
# #     print("Market is currently closed. It will open in:", time_until_next_session)
# # else:
# #     print("Market is currently open.")

# # from bitpay.client import Client

# # # Initialize BitPay client with your API token
# # bitpay = Client.create_pos_client(api_token='YOUR_API_TOKEN')

# # # Create a payment invoice
# # invoice_data = {
# #     "price": 100,  # Amount in Tether
# #     "currency": "USDT",
# #     # Additional options can be added here, like "redirectURL" for redirection after payment
# # }
# # invoice = bitpay.create_invoice(invoice_data)

# # # Extract QR code URL from the invoice
# # qr_code_url = invoice['data']['url']
# # print("QR Code URL:", qr_code_url)


#!/usr/bin/env python3

# import requests
# from datetime import datetime, timedelta
# import pytz

# url = "https://api.tradinghours.com/v3/markets/status?fin_id=XNYS&timezone=utc"
# # url = "https://api.tradinghours.com/v3/markets?group=all"
# headers = {
#     "Content-Type": "application/json",
#     "Authorization": "Bearer 1UgpHDBNPiXbr9mrp873nf7NV5JVQHjOPiRlyWGRbbacfb84"
# }

# response = requests.get(url, headers=headers)

# data = response.json()
# print(data)
# # for i in data["data"]:
# #     try:
# #         if i["fin_id"] != None:
# #             if "AE." in (i["fin_id"]):
# #                 print(i)
# #     except:
# #         pass
# market_data = data['data']['US.NYSE']
# status = market_data['status']
# next_bell_utc = market_data['next_bell']

# next_bell_utc = datetime.strptime(next_bell_utc, "%Y-%m-%dT%H:%M:%S%z")
# utc_time = datetime.strptime(data['meta']['utc_time'], "%Y-%m-%dT%H:%M:%S%z")

# if utc_time >= next_bell_utc:
#     status = 'Open'
# else:
#     status = 'Closed'

# print(f"Market status based on UTC time: {status} {int(next_bell_utc.timestamp())}")
