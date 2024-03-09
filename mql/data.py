from flask import Flask, request, jsonify
from datetime import datetime
import MetaTrader5 as mt5
import requests

url = 'https://goldshop24.co/api/window/get_account'
secret = "MysticalDragon$7392&WhisperingWinds&SunsetHaven$AuroraBorealis"

headers = {
    'Content-Type': 'application/json',
    'Authorization': secret,
}

app = Flask(__name__)

path = "c:\\Program Files\\MetaTrader 5\\terminal64.exe"

def initialize_mt5():
    try:
        account = requests.get(url, headers=headers).json()["accounts"]
        if not mt5.initialize():
            return

        if mt5.login(account["login"], password=account["password"], server=account["server"]):
            return True
 
        return False
    except:
        pass

@app.before_request
def before_request():
    if 'X-Secret-Header' not in request.headers or request.headers['X-Secret-Header'] != secret:
        return jsonify({"status":False,"message": "Unauthorized"}), 401

    if not hasattr(request, 'mt5_initialized'):
        request.mt5_initialized = True
        if not mt5.initialize():
            return jsonify({"status":False,"message": "Unauthorized"}), 401

@app.route('/get_history', methods=['GET'])
def get_history():
    from_date = datetime(2020, 1, 1)
    history_orders = mt5.history_orders_get(from_date, datetime.now())

    formatted_orders = []
    for order in history_orders:
        formatted_order = {
            "ticket": order[0],
            "time_setup": order[1],
            "time_setup_msc": order[2],
            "time_done": order[3],
            "time_done_msc": order[4],
            "time_expiration": order[5],
            "type": order[6],
            "type_time": order[7],
            "type_filling": order[8],
            "state": order[9],
            "magic": order[10],
            "position_id": order[11],
            "position_by_id": order[12],
            "reason": order[13],
            "volume_initial": order[14],
            "volume_current": order[15],
            "price_open": order[16],
            "sl": order[17],
            "tp": order[18],
            "price_current": order[19],
            "price_stoplimit": order[19],
            "symbol": order[21],
            "comment": order[22]
        }
        formatted_orders.append(formatted_order)

    return jsonify({"data": formatted_orders, "status": True}), 200

@app.route('/reinitialize', methods=['POST'])
def trigger_reinitialize():
    if request.method == 'POST':
        if initialize_mt5():
            return jsonify({"status": True, "message": "MetaTrader 5 reinitialized successfully"}), 200
        else:
            return jsonify({"status": False, "message": "Failed to reinitialize MetaTrader 5"}), 500
    else:
        return jsonify({"status": False, "message": "Invalid request method"}), 405


@app.route('/positions', methods=['GET'])
def positions():
    orders = mt5.positions_get()
    
    parameter_names = [
        "ticket", "magic_number", "order_id", "position_id", "position_by_id", 
        "volume", "position_time", "position_time_msc", "type", "volume_initial",
        "price_open", "sl", "tp", "price_current", "swap", "profit", "symbol",
        "comment", "external_id"
    ]

    positions_list = []
    for order in orders:
        position_dict = {}
        for i, parameter_name in enumerate(parameter_names):
            position_dict[parameter_name] = order[i]
        positions_list.append(position_dict)

    return jsonify({"data": positions_list, "status": True}), 200

@app.route('/send_order', methods=['POST'])
def send_order():
    symbol = request.json.get('symbol')
    lot = request.json.get('volume')
    deviation = request.json.get('deviation', 20)
    type = request.json.get('type')

    if mt5.symbol_info_tick(symbol) == None:
        return jsonify({"status":False,"message": "Invalid symbol", "data": "invalid symbol"}), 400
    request_data = {
        "action": mt5.TRADE_ACTION_DEAL,
        "symbol": symbol,
        "volume": float(lot),
        "type": type,
        "price": mt5.symbol_info_tick(symbol).ask,
        "deviation": deviation,
        "magic": 12345,
        "comment": "carat",
        "type_time": mt5.ORDER_TIME_GTC,
        "type_filling": mt5.ORDER_FILLING_FOK,
    }

    result = mt5.order_send(request_data)
    if result == None:
        return jsonify({"status":False,"message": "Order send failed", "data": str(mt5.last_error()[1])}), 400
    else:
        if result.retcode != 1009 and result.retcode != 1008:
            return jsonify({"status":False,"message": "Order send failed", "data": result.comment}), 400
        return jsonify({"status":True,"message": "Order placed successfully", "data": result.order}), 200

@app.route('/cancel_order', methods=['POST'])
def cancel_order():
    ticket_id = request.json.get('ticket_id')
    position = None
    for p in mt5.positions_get():
        if p.ticket == ticket_id:
            position = p
            break

    if position:
        side = 0
        if position.type == mt5.ORDER_TYPE_BUY:
            side = mt5.ORDER_TYPE_SELL
        elif position.type == mt5.ORDER_TYPE_SELL:
            side = mt5.ORDER_TYPE_BUY
        elif position.type == mt5.ORDER_TYPE_BUY_LIMIT:
            side = mt5.ORDER_TYPE_SELL_LIMIT
        elif position.type == mt5.ORDER_TYPE_SELL_LIMIT:
            side = mt5.ORDER_TYPE_BUY_LIMIT
        elif position.type == mt5.ORDER_TYPE_BUY_STOP:
            side = mt5.ORDER_TYPE_SELL_STOP
        elif position.type == mt5.ORDER_TYPE_SELL_STOP:
            side = mt5.ORDER_TYPE_BUY_STOP
        elif position.type == mt5.ORDER_TYPE_BUY_STOP_LIMIT:
            side = mt5.ORDER_TYPE_SELL_STOP_LIMIT
        elif position.type == mt5.ORDER_TYPE_SELL_STOP_LIMIT:
            side = mt5.ORDER_TYPE_BUY_STOP_LIMIT
        requestMt5 = {
            "action": mt5.TRADE_ACTION_CLOSE_BY,
            "position_by": ticket_id,
            "symbol": position.symbol,
            "volume": position.volume,
            "type": side,
            "tp": position.tp,
            "sl": position.sl,
            "price": mt5.symbol_info_tick(position.symbol).bid,
        }
        result = mt5.order_send(requestMt5)
        if result.retcode == mt5.TRADE_RETCODE_DONE:
            return jsonify({"status":False,"data": "Position closed successfully", "ticket_id": ticket_id}), 200
        else:
            return jsonify({"status":False,"data": result.comment, "result": str(result)}), 400
    else:
        return jsonify({"status":False,"data": "Position not found"}), 404

@app.route('/account_info', methods=['GET'])
def account_info():
    if not hasattr(request, 'mt5_initialized'):
        return jsonify({"status": False, "message": "MetaTrader5 not initialized"}), 500

    account_info = mt5.account_info()
    if account_info:
        account_info_dict = {
            "login": account_info[0],
            "balance": account_info[1],
            "credit": account_info[2],
            "company": account_info[3],
            "currency": account_info[4],
            "server": account_info[5],
            "stopout_level": account_info[6],
            "leverage": account_info[7],
            "agent_account": account_info[8],
            "margin_so_mode": account_info[9],
            "trade_allowed": account_info[10],
            "trade_expert": account_info[11],
            "margin_mode": account_info[12],
            "currency_digits": account_info[13],
            "fifo_close": account_info[14],
            "balance_status": account_info[15],
            "credit_status": account_info[16],
            "email": account_info[17]
        }
        return successful(account_info_dict), 200
    else:
        return jsonify({"status": False, "message": "Failed to retrieve account info"}), 500

def successful(data):
    return jsonify({"status": True, "data": data})

initialize_mt5()
app.run(debug=True ,port=80, host="172.31.24.144")
