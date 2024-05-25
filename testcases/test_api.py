import asyncio
import websockets
import requests
import datetime

token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI2NjMyMmNjNWFiN2QzMDU1MWZmMTg5OTgiLCJuYW1lIjoiIiwicGhvbmUiOiI5ODkxMzg3ODAyNzUiLCJlbWFpbCI6IiIsImNyZWF0ZWRfYXQiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiIsImV4cCI6MTcxNzQyMzMzNX0.qGJYVLEji-CRofRfSxrQTNDeyY2Uvc0sc_OPn9RvMpg"
# register = requests.post("https://goldshop24.co/api/auth/user/register", json={
#     # "name": "Mehrbod man",
#     "phone": "989138780275",
#     "otp_code": 12345,
#     # "address": [
#     #     {"address": "King bernam", "label": "my home"},
#     #     {"address": "my queen nam bernam", "label": "my work"},
#     # ]
# }, headers={
#     "Authorization": token
# })
# print(register.text)
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

data = requests.get("https://goldshop24.co/api/user/payment_methods", json={
    # "name": "Mehrbod man",
    "phone": "989138780275",
    "otp_code": 12345,
    # "address": [
    #     {"address": "King bernam", "label": "my home"},
    #     {"address": "my queen nam bernam", "label": "my work"},
    # ]
}, headers={
    "Authorization": token
})
print(data.text)