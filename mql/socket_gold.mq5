int socket;

ENUM_TIMEFRAMES timeframes[] = {PERIOD_M1, PERIOD_M5, PERIOD_M15, PERIOD_H1, PERIOD_D1, PERIOD_W1, PERIOD_MN1};

string DataSymbol(string symbolName, string type) {
    double ask = SymbolInfoDouble(symbolName, SYMBOL_ASK);
    double bid = SymbolInfoDouble(symbolName, SYMBOL_BID);
    double high = iHigh(symbolName, PERIOD_CURRENT, 0);
    double low = iLow(symbolName, PERIOD_CURRENT, 0);
    double open = iOpen(symbolName, PERIOD_CURRENT, 0);
    double close = iClose(symbolName, PERIOD_CURRENT, 0);
    ENUM_TIMEFRAMES timeframe = Period();
    long currentDateTimeString = TimeCurrent() * 1000;

    if (ask == 0 || bid == 0 || high == 0 || low == 0 || open == 0 || close == 0) {
        PrintFormat("Error retrieving data for %s. Ask: %f, Bid: %f, High: %f, Low: %f, Open: %f, Close: %f", 
                    symbolName, ask, bid, high, low, open, close);
    }

    return StringFormat(
        "{\"time\":\"%d\",\"symbol\":\"%s\",\"ask\":%f,\"bid\":%f,\"high\":%f,\"low\":%f,\"open\":%f,\"close\":%f,\"timeframe\":\"%s\",\"type\":\"%s\"}",
        currentDateTimeString, symbolName, ask, bid, high, low, open, close, EnumToString(timeframe), type
    );
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
    int maxRetries = 5000;
    int retryCount = 0;

    while (retryCount < maxRetries) {
        socket = SocketCreate();

        if (retryCount == maxRetries) {
            Print("Connection failed after ", maxRetries, " retries");
            return(INIT_FAILED);
        }

        if (SocketConnect(socket, serverAddress, serverPort, 1000)) {
            EventSetTimer(5);
            Print("Connected to the first server");

            GetSymbolList1();
            return(INIT_SUCCEEDED);
        }

        int error = GetLastError();
        Print("Retry attempt (Server 1): ", retryCount + 1, " error : ", error);
        retryCount++;

        Sleep(1000);
    }

    return(INIT_FAILED);
}

void OnDeinit(const int reason) {
    SocketClose(socket);
    EventKillTimer();
}

void ConnectionAndTick() {
    if (!SocketIsConnected(socket)) {
        Print("Disconnected from the first server");
        socket = SocketCreate();
        if (SocketConnect(socket, serverAddress, serverPort, 1000)) {
            Print("Connected to the first server");
        } else {
            int error = GetLastError();
            Print(" error : ", error);
            Sleep(1000);
        }
    }

    GetSymbolList1();
}

void OnTick() {
    ConnectionAndTick();
}

void OnTimer() {
    ConnectionAndTick();
}
