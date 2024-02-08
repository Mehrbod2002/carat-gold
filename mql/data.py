from flask import Flask, request, jsonify
from datetime import datetime
import MetaTrader5 as mt5

app = Flask(__name__)


login = 51981488
password = "bogxs0lz"
server = "Alpari-MT5-Demo"
secret = "MysticalDragon$7392&WhisperingWinds&SunsetHaven$AuroraBorealis"

def initialize_mt5():
    if not mt5.initialize():
        if mt5.login(login, password=password, server=server):
            return False
    return True

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
            "magic_number": order[1],
            "order_type": order[2],
            "position_id": order[3],
            "position_by_id": order[4],
            "volume": order[5],
            "time_setup": order[6],
            "time_done": order[7],
            "time_expiration": order[8],
            "type_time": order[9],
            "magic": order[10],
            "position": order[11],
            "reason": order[12],
            "type": order[13],
            "volume_initial": order[14],
            "sl": order[15],
            "tp": order[16],
            "price_open": order[17],
            "price_current": order[18],
            "price_stoplimit": order[19],
            "symbol": order[20],
            "comment": order[21],
            "external_id": order[22]
        }
        formatted_orders.append(formatted_order)

    return jsonify({"data": formatted_orders, "status": True}), 200

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
    lot = request.json.get('lot', 0.1)
    deviation = request.json.get('deviation', 20)
    type = request.json.get('type')

    request_data = {
        "action": mt5.TRADE_ACTION_DEAL,
        "symbol": symbol,
        "volume": lot,
        "type": type,
        "price": mt5.symbol_info_tick(symbol).ask,
        "deviation": deviation,
        "magic": "magic carat",
        "comment": "carat",
        "type_time": mt5.ORDER_TIME_GTC,
        "type_filling": mt5.ORDER_FILLING_FOK,
    }

    result = mt5.order_send(request_data)
    if result.retcode != mt5.TRADE_RETCODE_DONE:
        return jsonify({"status":False,"message": "Order send failed", "result": str(result)}), 400
    else:
        return jsonify({"status":True,"message": "Order placed successfully", "order": result.order, "trade_result": str(result)}), 200

@app.route('/cancel_order/<int:ticket_id>', methods=['DELETE'])
def cancel_order(ticket_id):
    position = None
    for p in mt5.positions_get():
        if p.ticket == ticket_id:
            position = p
            break

    if position:
        side = mt5.ORDER_TYPE_BUY if position.type == 0 else mt5.ORDER_TYPE_SELL
        request = {
            "action": mt5.TRADE_ACTION_CLOSE_BY,
            "position": ticket_id,
            "symbol": position.symbol,
            "volume": position.volume,
            "type": side,
            "price": mt5.symbol_info_tick(position.symbol).bid,
        }
        result = mt5.order_send(request)

        if result.retcode == mt5.TRADE_RETCODE_DONE:
            return jsonify({"status":False,"message": "Position closed successfully", "ticket_id": ticket_id}), 200
        else:
            return jsonify({"status":False,"message": "Failed to close position", "result": str(result)}), 400
    else:
        return jsonify({"status":False,"message": "Position not found"}), 404

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

if __name__ == '__main__':
    initialize_mt5()
    app.run(debug=True ,port=80, host="172.31.24.144")
