int socket;

ENUM_TIMEFRAMES timeframes[] = {PERIOD_M1, PERIOD_M5, PERIOD_M15, PERIOD_H1, PERIOD_D1, PERIOD_W1, PERIOD_MN1};

string DataSymbol(string symbolName, string type) {
    double ask = SymbolInfoDouble(symbolName, SYMBOL_ASK);
    double bid = SymbolInfoDouble(symbolName, SYMBOL_BID);
    // double high = iHigh(symbolName, PERIOD_CURRENT, 0);
    // double low = iLow(symbolName, PERIOD_CURRENT, 0);
    // double open = iOpen(symbolName, PERIOD_CURRENT, 0);
    // double close = iClose(symbolName, PERIOD_CURRENT, 0);
    // ENUM_TIMEFRAMES timeframe = Period();
    // long currentDateTimeString = TimeCurrent() * 1000;

    return StringFormat(
        "{\"symbol\":\"%s\",\"ask\":%f,\"bid\":%f}",
        symbolName, ask, bid
    );
    // return StringFormat(
    //    "{\"time\":\"%d\",\"symbol\":\"%s\",\"ask\":%f,\"bid\":%f,\"high\":%f,\"low\":%f,\"open\":%f,\"close\":%f,\"timeframe\":\"%s\",\"type\":\"%s\"}",
    //    currentDateTimeString, symbolName, ask, bid, high, low, open, close, EnumToString(timeframe), type
    //);
}

void GetSymbolList1() {
    string jsonData = DataSymbol("XAUUSD", "commodity");
    uchar byteArray[];
    StringToCharArray(jsonData, byteArray);
    SocketSend(socket, byteArray, StringLen(jsonData));
}

string serverAddress = "52.0.87.119";
int serverPort = 5741;

int OnInit() {
    socket = SocketCreate();
    if (SocketConnect(socket, serverAddress, serverPort, 1000)) {
        Print("Connected to the server");
        return(INIT_SUCCEEDED);
    } else {
        Print("Connection failed with error: ", GetLastError());
        return(INIT_FAILED);
    }
}

void OnDeinit(const int reason) {
    SocketClose(socket);
}

void ConnectionAndTick() {
    if (!SocketIsConnected(socket)) {
        Print("Disconnected from the server");
        socket = SocketCreate();
        if (SocketConnect(socket, serverAddress, serverPort, 1000)) {
            Print("Reconnected to the server");
        } else {
            Print("Reconnection failed with error: ", GetLastError());
        }
    }
    GetSymbolList1();
}

void OnTick() {
    ConnectionAndTick();
}