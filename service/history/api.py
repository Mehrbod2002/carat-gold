from flask import Flask, jsonify, request, make_response
from websocket import create_connection
from flask_cors import CORS
import json
import random
import string
import re

app = Flask(__name__)
CORS(app, resources={r"/history": {"origins": "*"}})


def generateSession():
    stringLength = 12
    letters = string.ascii_lowercase
    random_string = ''.join(random.choice(letters)
                            for i in range(stringLength))
    return "qs_" + random_string


def generateChartSession():
    stringLength = 12
    letters = string.ascii_lowercase
    random_string = ''.join(random.choice(letters)
                            for i in range(stringLength))
    return "cs_" + random_string


def prependHeader(st):
    return "~m~" + str(len(st)) + "~m~" + st


def constructMessage(func, paramList):
    return json.dumps({
        "m": func,
        "p": paramList
    }, separators=(',', ':'))


def createMessage(func, paramList):
    return prependHeader(constructMessage(func, paramList))


def sendRawMessage(ws, message):
    ws.send(prependHeader(message))


def sendMessage(ws, func, args):
    ws.send(createMessage(func, args))


@app.route('/history', methods=['POST'])
def get_data():
    data = request.get_json()

    try:
        if len(data["symbol"]) == 0 or len(data["symbol"]) == 50:
            response = jsonify({"status": False, "m": "invalid_parameters"})
            return make_response(response, 406)
        if data["count"] == 0 or data["to"] == 0 or data["until"] == 0:
            response = jsonify({"status": False, "m": "invalid_parameters"})
            return make_response(response, 406)
    except KeyError:
        response = jsonify({"status": False, "m": "invalid_parameters"})
        return make_response(response, 406)

    headers = json.dumps({
        'Origin': 'https://data.tradingview.com'
    })

    ws = create_connection(
        'wss://data.tradingview.com/socket.io/websocket', headers=headers)
    session = generateSession()
    chart_session = generateChartSession()

    data["to"] = int(data["to"])
    data["until"] = int(data["until"])

    sendMessage(ws, "set_auth_token", ["unauthorized_user_token"])
    sendMessage(ws, "chart_create_session", [chart_session, ""])
    sendMessage(ws, "quote_create_session", [session])
    sendMessage(ws, "quote_set_fields", [session, "ch", "chp", "current_session", "description", "local_description", "language", "exchange", "fractional", "is_tradable",
                "lp", "lp_time", "minmov", "minmove2", "original_name", "pricescale", "pro_name", "short_name", "type", "update_mode", "volume", "currency_code", "rchp", "rtc"])
    # sendMessage(ws, "quote_add_symbols", [
    #             session, data['symbol'], {}])
    sendMessage(ws, "resolve_symbol", [
                chart_session, "symbol_1", "={\"symbol\":\""+data['symbol']+"\",\"adjustment\":\"splits\"}"])
    sendMessage(ws, "create_series", [
                chart_session, "s1", "s1", "symbol_1", data['timeframe'], int(data['count']), "r,"+str(data['until'])+":"+str(data['to'])])

    while True:
        try:
            result = ws.recv()
            pattern = re.compile("~m~\d+~m~~h~\d+$")
            if pattern.match(result):
                ws.recv()
                ws.send(result)

            for i in result.split("~m~"):
                if "error" in i:
                    err = json.loads(i)
                    if err["m"] == "critical_error":
                        response = jsonify({"status": False, "m": err["p"][1]})
                        return make_response(response, 500)
                    if err["m"] == "symbol_error":
                        response = jsonify({"status": False, "m": err["p"][2]})
                        return make_response(response, 500)
                    if err["m"] == "series_error":
                        response = jsonify({"status": False, "m": err["p"][3]})
                        return make_response(response, 500)

            if "timescale_update" in str(result):
                for i in result.split("~m~"):
                    if "timescale_update" in i:
                        loadData = json.loads(i)
                        response = jsonify(loadData['p'][1]['s1']["s"])
                        return make_response(response, 200)
        except Exception:
            response = jsonify({"status": False, "m": "internal_error"})
            return make_response(response, 500)


@app.route('/market_status', methods=['GET'])
def get_status():
    data = {"symbol": "FX:XAUUSD"}

    headers = json.dumps({
        'Origin': 'https://data.tradingview.com'
    })

    ws = create_connection(
        'wss://data.tradingview.com/socket.io/websocket', headers=headers)
    session = generateSession()
    chart_session = generateChartSession()

    sendMessage(ws, "set_auth_token", ["unauthorized_user_token"])
    sendMessage(ws, "chart_create_session", [chart_session, ""])
    sendMessage(ws, "quote_create_session", [session])
    sendMessage(ws, "quote_set_fields", [session, "ch", "chp", "current_session", "description", "local_description", "language", "exchange", "fractional", "is_tradable",
                "lp", "lp_time", "minmov", "minmove2", "original_name", "pricescale", "pro_name", "short_name", "type", "update_mode", "volume", "currency_code", "rchp", "rtc"])
    sendMessage(ws, "resolve_symbol", [
                chart_session, "symbol_1", "={\"symbol\":\""+data['symbol']+"\",\"adjustment\":\"splits\"}"])

    while True:
        try:
            result = ws.recv()
            pattern = re.compile("~m~\d+~m~~h~\d+$")
            if pattern.match(result):
                ws.recv()
                ws.send(result)

            for i in result.split("~m~"):
                if "error" in i:
                    err = json.loads(i)
                    if err["m"] == "critical_error":
                        response = jsonify({"status": False, "m": err["p"][1]})
                        return make_response(response, 500)
                    if err["m"] == "symbol_error":
                        response = jsonify({"status": False, "m": err["p"][2]})
                        return make_response(response, 500)
                    if err["m"] == "series_error":
                        response = jsonify({"status": False, "m": err["p"][3]})
                        return make_response(response, 500)

            if "symbol_resolved" in str(result):
                for i in result.split("~m~"):
                    if "symbol_resolved" in i:
                        loadData = json.loads(i)
                        response = jsonify(loadData['p'][2])
                        return make_response(response, 200)
        except Exception:
            response = jsonify({"status": False, "m": "internal_error"})
            return make_response(response, 500)


if __name__ == '__main__':
    app.run(debug=True)
