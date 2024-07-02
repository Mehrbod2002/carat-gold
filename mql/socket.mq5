int socket;

ENUM_TIMEFRAMES timeframes[] = {PERIOD_M1, PERIOD_M5, PERIOD_M15, PERIOD_H1, PERIOD_D1, PERIOD_W1, PERIOD_MN1};

string forexSymbols[] = {
    "TESLA_CFD",
    "APPLE_CFD",
    "PAYPAL_HOLDINGS_CFD",
    "META_PLATFORMS_CFD",
    "AMAZON_CFD",
};

string cryptoSymbols[] = {
    "BITCOIN",
    "ETHEREUM",
    "BNB",
    "DOGECOIN",
};

string commoditiesSymbols[] = {
    "EURUSD",
    "GBPUSD",
};

string DataSymbol(string symbolName, string type) {
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
    for (int i = 0; i < ArraySize(forexSymbols); i++) {
        string jsonData1 = DataSymbol(forexSymbols[i], "forex");
        uchar byteArray1[];
        StringToCharArray(jsonData1, byteArray1);
        int sentBytes = SocketSend(socket, byteArray1, StringLen(jsonData1));
    }
    for (int i = 0; i < ArraySize(commoditiesSymbols); i++) {
        string jsonData3 = DataSymbol(commoditiesSymbols[i], "commodity");
        uchar byteArray3[];
        StringToCharArray(jsonData3, byteArray3);
        int sentBytes = SocketSend(socket, byteArray3, StringLen(jsonData3));
    }
    for (int i = 0; i < ArraySize(cryptoSymbols); i++) {
        string jsonData4 = DataSymbol(cryptoSymbols[i], "crypto");
        uchar byteArray4[];
        StringToCharArray(jsonData4, byteArray4);
        int sentBytes = SocketSend(socket, byteArray4, StringLen(jsonData4));
    }
}

string serverAddress = "52.0.87.119";
int serverPort = 5742;

int OnInit() {
    int maxRetries = 5000;
    int retryCount = 0;

    while (retryCount < maxRetries) {
        socket = SocketCreate();

        if (maxRetries == retryCount) {
            Print("Connection failed after ", maxRetries, " retries");
            return(INIT_FAILED);
        }

        if (SocketConnect(socket, serverAddress, serverPort, 1000)) {
            EventSetMillisecondTimer(3);
            Print("Connected to the first server");
            return(INIT_SUCCEEDED);
        }

        int error = GetLastError();
        Print("Retry attempt (Server 1): ", retryCount + 1," error : ", error);
        retryCount++;

        Sleep(1000);
    }

    return(INIT_FAILED);
}

void OnDeinit(const int reason) {
    SocketClose(socket);
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

    GetSymbolList();
}

void OnTimer() {
    ConnectionAndTick();
}

