from flask import Flask, request, jsonify
from datetime import datetime
import MetaTrader5 as mt5

app = Flask(__name__)


login = 51981488
password = "bogxs0lz"
server = "Alpari-MT5-Demo"
secret = "secret"

def initialize_mt5():
    if not mt5.initialize():
        if mt5.login(login, password=password, server=server)
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
    history_orders = mt5.history_orders_total(from_date, datetime.now())
    if history_orders > 0:
        return jsonify({"status":False,"message": "Total history orders={}".format(history_orders)}), 200
    else:
        return jsonify({"status":False,"message": "Orders not found in history"}), 404

@app.route('/positions', methods=['GET'])
def positions():
    orders = mt5.positions_get()
    return jsonify(orders), 200

@app.route('/send_order', methods=['POST'])
def send_order():
    symbol = request.json.get('symbol')
    lot = request.json.get('lot', 0.1)
    deviation = request.json.get('deviation', 20)

    request_data = {
        "action": mt5.TRADE_ACTION_DEAL,
        "symbol": symbol,
        "volume": lot,
        "type": mt5.ORDER_TYPE_BUY,
        "price": mt5.symbol_info_tick(symbol).ask,
        "deviation": deviation,
        "magic": 234000,
        "comment": "Python Script Order",
        "type_time": mt5.ORDER_TIME_GTC,
        "type_filling": mt5.ORDER_FILLING_FOK,
    }

    result = mt5.order_send(request_data)
    if result.retcode != mt5.TRADE_RETCODE_DONE:
        return jsonify({"status":False,"message": "Order send failed", "result": str(result)}), 400
    else:
        return jsonify({"status":False,"message": "Order placed successfully", "order": result.order, "trade_result": str(result)}), 200

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

if __name__ == '__main__':
    initialize_mt5()
    app.run(debug=True)
