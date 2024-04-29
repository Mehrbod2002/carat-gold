int socket;

ENUM_TIMEFRAMES timeframes[] = {PERIOD_M1, PERIOD_M5, PERIOD_M15, PERIOD_H1, PERIOD_D1, PERIOD_W1, PERIOD_MN1};

string commoditiesSymbols[] = {
    "XAUUSD",
};

string DataSymbol(string symbolName,string type) {
    double ask = SymbolInfoDouble(symbolName, SYMBOL_BID);
    double bid = SymbolInfoDouble(symbolName, SYMBOL_ASK);
    double high = iHigh(symbolName, PERIOD_CURRENT, 0);
    double low = iLow(symbolName, PERIOD_CURRENT, 0);
    double open = iOpen(symbolName, PERIOD_CURRENT, 0);
    double close = iClose(symbolName, PERIOD_CURRENT, 0);
    ENUM_TIMEFRAMES timeframe = Period();
    long currentDateTimeString = TimeCurrent() * 1000;

    return StringFormat(
        "{\"time\":\"%d\",\"symbol\":\"%s\",\"ask\":%f,\"bid\":%f,\"high\":%f,\"low\":%f,\"open\":%f,\"close\":%f,\"timeframe\":\"%s\",\"type\":\"%s\"}",
        currentDateTimeString, symbolName, bid, ask, high, low, open, close, EnumToString(timeframe), type
    );
}

void GetSymbolList() {
    for (int i = 0; i < ArraySize(commoditiesSymbols); i++) {
        string jsonData = DataSymbol(commoditiesSymbols[i], "commodity");
        uchar byteArray[];
        StringToCharArray(jsonData, byteArray);
        SocketSend(socket, byteArray, StringLen(jsonData));
    }
}

string serverAddress = "52.0.87.119";
int serverPort = 5741;

int OnInit() {
    int maxRetries = 5000;
    int retryCount = 0;

    while (retryCount < maxRetries) {
        socket = SocketCreate();

        if (SocketConnect(socket, serverAddress, serverPort, 1000)) {
            EventSetMillisecondTimer(800);
            Print("Connected to the server");
            return(INIT_SUCCEEDED);
        }

        int error = GetLastError();
        Print("Retry attempt: ", retryCount + 1," error : ", error);
        retryCount++;

        Sleep(1000);
    }

    Print("Connection failed after ", maxRetries, " retries");
    return(INIT_FAILED);
}

void OnDeinit(const int reason) {
    SocketClose(socket);
}

void ConnectionAndTick() {
    if (!SocketIsConnected(socket)) {
        Print("Disconnected from the server");
        socket = SocketCreate();
        if (SocketConnect(socket, serverAddress, serverPort, 1000)) {
            Print("Connected to the server");
        } else {
            int error = GetLastError();
            Print(" error : ", error);
            Sleep(1000);
        }
    }

    GetSymbolList();
}

void OnTick() {
    ConnectionAndTick();
}

void OnTimer() {
    ConnectionAndTick();
}

