import requests

url = 'https://goldshop24.co/api/window/get_account'
secret = "MysticalDragon$7392&WhisperingWinds&SunsetHaven$AuroraBorealis"

headers = {
    'Content-Type': 'application/json',
    'Authorization': secret,
}

account = requests.get(url, headers=headers).json()

print(account)