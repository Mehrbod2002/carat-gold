#!/usr/bin/env python3
from bson import ObjectId
from datetime import datetime
import threading
import websockets
from websockets import WebSocketServerProtocol
from urllib.parse import urlparse, parse_qs
import asyncio
import traceback
import json
from datetime import datetime, timedelta
import pymongo

db_user = "financial_data"
db_pass = "sgQRVCeGCJ29WK4zsgQRVCeGCJ2k6A"
db_port = 27017
db_host = "localhost"
url = f"mongodb://{db_user}:{db_pass}@{db_host}:{db_port}/?authSource=admin"
mongo_client = pymongo.MongoClient(url)
db = mongo_client["caratGold"]
collection = db["history"]
server_address = '0.0.0.0'
clients = set()
lock = threading.Lock()


def printer(*message):
    print(message, flush=True)


printer("Start python code ....")


async def handle_client_wss(websocket, path):
    if path == "/feed":
        with lock:
            clients.add(websocket)
        try:
            while True:
                async for data in websocket:
                    if data == "ping":
                        await websocket.send("pong".encode("utf-8"))
        except (websockets.ConnectionClosedOK, websockets.ConnectionClosedError):
            with lock:
                try:
                    clients.remove(websocket)
                except KeyError:
                    pass
    elif "history" in path:
        try:
            parsed_path = urlparse(path)
            params = parse_qs(parsed_path.query)
            to = params.get('to', None)
            _from = params.get('from', None)
            symbol = params.get("symbol", None)[0]
            result = collection.find_one({"symbol": symbol})
            if result:
                await websocket.send(json.dumps(result["data_array"]))
            await websocket.close()
        except Exception as e:
            await websocket.send([])
            await websocket.close()


async def handle_client_8080(reader, writer):
    buffer = b""
    last_bars_dict = {}  # Dictionary to store last bars for each symbol

    while True:
        try:
            data = await reader.read(500)
            if not data:
                break
            buffer += data

            while b'}' in buffer:
                json_start = buffer.find(b'{')
                json_end = buffer.find(b'}')
                if json_start != -1 and json_end != -1 and json_start < json_end:
                    complete_json = buffer[json_start:json_end +
                                           1].decode('utf-8')
                    buffer = buffer[json_end + 1:]
                    complete_json = json.loads(complete_json)

                    current_time_utc = datetime.utcnow()
                    current_timestamp = int(current_time_utc.timestamp())
                    complete_json["time"] = current_timestamp

                    try:
                        with lock:
                            for client in set(clients):
                                try:
                                    await client.send(json.dumps(complete_json))
                                except (websockets.ConnectionClosedOK, websockets.ConnectionClosedError):
                                    clients.discard(client)
                    except RuntimeError as e:
                        traceback.print_exc()
                        print(f"Error sending data to clients: {e}")

                    symbol = complete_json.get("symbol", None)
                    trade_price = float(complete_json["bid"])

                    if symbol not in last_bars_dict:
                        last_bars_dict[symbol] = {
                            "symbol": symbol,
                            "time": complete_json["time"],
                            "open": trade_price,
                            "high": trade_price,
                            "low": trade_price,
                            "close": trade_price,
                        }
                    else:
                        last_bars_dict[symbol]["high"] = max(
                            last_bars_dict[symbol]["high"], trade_price)
                        last_bars_dict[symbol]["low"] = min(
                            last_bars_dict[symbol]["low"], trade_price)
                        last_bars_dict[symbol]["close"] = trade_price

                        query = {
                            "symbol": symbol,
                            "data_array.time": complete_json["time"]
                        }

                        existing_document = collection.find_one(query)

                        if existing_document:
                            update = {
                                "$set": {
                                    "data_array.$.high": last_bars_dict[symbol]["high"],
                                    "data_array.$.low": last_bars_dict[symbol]["low"],
                                    "data_array.$.close": trade_price
                                }
                            }
                            collection.update_one(query, update)
                        else:
                            bar = last_bars_dict[symbol]
                            update = {
                                "$push": {"data_array": bar},
                            }
                            collection.update_one(
                                {"symbol": symbol}, update, upsert=True)

        except Exception as e:
            traceback.print_exc()
            print(f"Error processing data: {e}")
            continue


async def start_server_wss():
    server = await websockets.serve(handle_client_wss, server_address, 5050, ping_interval=60, create_protocol=WebSocketServerProtocol,  ping_timeout=30 * 24 * 3600)

    async with server:
        await server.wait_closed()


async def start_server_8080():
    server = await asyncio.start_server(handle_client_8080, server_address, 5741)

    async with server:
        await server.serve_forever()


def get_next_bar_time(bar_time):
    if bar_time > 10**12:
        date = datetime.utcfromtimestamp(bar_time / 1000)
    else:
        date = datetime.utcfromtimestamp(bar_time)
    date += timedelta(seconds=60)

    return int(date.timestamp())


async def main():
    tasks = [start_server_wss(), start_server_8080()]

    try:
        await asyncio.gather(*tasks)
    except KeyboardInterrupt as Output:
        printer("Output : ", Output)
    finally:
        printer("Loop closed")

loop = asyncio.new_event_loop()
asyncio.set_event_loop(loop)

try:
    loop.run_until_complete(main())
except KeyboardInterrupt:
    pass
finally:
    loop.close()
