int socket1;
int socket2;

ENUM_TIMEFRAMES timeframes[] = {PERIOD_M1, PERIOD_M5, PERIOD_M15, PERIOD_H1, PERIOD_D1, PERIOD_W1, PERIOD_MN1};

string forexSymbols[] = {
    "EURUSD",
    "GBPUSD",
    "AUDUSD",
    "GBPJPY",
    "USDJPY",
    "NZDUSD",
    "EURGBP",
    "AUDCHF",
    "AUDJPY",
    "AUDNZD",
    "CADCHF",
    "CADJPY",
    "CHFJPY",
    "EURAUD",
    "EURCAD",
    "EURCHF",
    "EYRJPY",
    "EURNOK",
    "EURNZD",
    "EURSEK",
    "GBPAUD",
    "GBPCAD",
    "GBPCHF",
    "GBPNZD",
    "NZDCAD",
    "NZDCHF",
    "NZDJPY",
    "AUDCAD",
    "USDTRY",
    "SGDJPY",
    "USDCAD",
    "USDCHF",
    "USDCNH",
    "USDHKD",
    "USDMXN",
    "USDNOK",
    "USDSEK",
    "USDSGD",
    "USDZAR",
    "EURCZK",
    "EURHUF",
    "EURPLN",
    "USDCZK",
    "USDHUF",
    "USDPLN",
    "ZARJPY",
    "TRYJPY",
    "EURTRY"
};

string cryptoSymbols[] = {
    "BTCUSD",
    "BNBUSD",
    "BCHUSD",
    "DOGEUSD",
    "DOTUSD",
    "EOSUSD",
    "ETHUSD",
    "LINKUSD",
    "ADAUSD",
    "LTCUSD",
    "MATICUSD",
    "UNIUSD",
    "XLMUSD",
    "XTZUSD",
    "AVAXUSD",
    "KSMUSD",
    "GLMRUSD",
    "SOLUSD"
};

string indexesSymbols[] = {
    "UK100",
    "US100",
    "US200",
    "US30",
    "US500",
    "CHINA50",
    "FRA40",
    "HK50",
    "JP225",
    "NL25",
    "SING30",
    "AUD200",
    "ES35",
    "CH20",
    "CHINAH",
    "GER30"
};

string commoditiesSymbols[] = {
    "XAUUSD",
    "XAGUSD",
    "UKOIL",
    "USOIL",
    "NATAGAS",
    "SOYBEANS",
    "COPPER",
    "SUGAR",
    "CORN",
    "WHEAT"
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

void GetSymbolList1() {
    string jsonData = DataSymbol("XAUUSD", "commodity");
    uchar byteArray[];
    StringToCharArray(jsonData, byteArray);
    SocketSend(socket2, byteArray, StringLen(jsonData));
}

void GetSymbolList2() {
    for (int i = 0; i < ArraySize(forexSymbols); i++) {
        string jsonData1 = DataSymbol(forexSymbols[i], "forex");
        uchar byteArray1[];
        StringToCharArray(jsonData1, byteArray1);
        int sentBytes = SocketSend(socket1, byteArray1, StringLen(jsonData1));
    }
    for (int i = 0; i < ArraySize(indexesSymbols); i++) {
        string jsonData2 = DataSymbol(indexesSymbols[i], "index");
        uchar byteArray2[];
        StringToCharArray(jsonData2, byteArray2);
        int sentBytes = SocketSend(socket1, byteArray2, StringLen(jsonData2));
    }
    for (int i = 0; i < ArraySize(commoditiesSymbols); i++) {
        string jsonData3 = DataSymbol(commoditiesSymbols[i], "commodity");
        uchar byteArray3[];
        StringToCharArray(jsonData3, byteArray3);
        int sentBytes = SocketSend(socket1, byteArray3, StringLen(jsonData3));
    }
    for (int i = 0; i < ArraySize(cryptoSymbols); i++) {
        string jsonData4 = DataSymbol(cryptoSymbols[i], "crypto");
        uchar byteArray4[];
        StringToCharArray(jsonData4, byteArray4);
        int sentBytes = SocketSend(socket1, byteArray4, StringLen(jsonData4));
    }
}

string serverAddress1 = "52.0.87.119";
int serverPort1 = 5742;

string serverAddress2 = "52.0.87.119";
int serverPort2 = 5741;

int OnInit() {
    int maxRetries = 5000;
    int retryCount = 0;

    while (retryCount < maxRetries) {
        socket1 = SocketCreate();

        if (SocketConnect(socket1, serverAddress1, serverPort1, 1000)) {
            EventSetMillisecondTimer(100);
            Print("Connected to the first server");
            break;
        }

        int error = GetLastError();
        Print("Retry attempt (Server 1): ", retryCount + 1," error : ", error);
        retryCount++;

        Sleep(1000);
    }

    retryCount = 0;

    while (retryCount < maxRetries) {
        socket2 = SocketCreate();

        if (SocketConnect(socket2, serverAddress2, serverPort2, 1000)) {
            EventSetMillisecondTimer(800);
            Print("Connected to the second server");
            return(INIT_SUCCEEDED);
        }

        int error = GetLastError();
        Print("Retry attempt (Server 2): ", retryCount + 1," error : ", error);
        retryCount++;

        Sleep(1000);
    }

    Print("Connection failed after ", maxRetries, " retries");
    return(INIT_FAILED);
}

void OnDeinit(const int reason) {
    SocketClose(socket1);
    SocketClose(socket2);
}

void ConnectionAndTick() {
    if (!SocketIsConnected(socket1)) {
        Print("Disconnected from the first server");
        socket1 = SocketCreate();
        if (SocketConnect(socket1, serverAddress1, serverPort1, 1000)) {
            Print("Connected to the first server");
        } else {
            int error = GetLastError();
            Print(" error : ", error);
            Sleep(1000);
        }
    }

    if (!SocketIsConnected(socket2)) {
        Print("Disconnected from the second server");
        socket2 = SocketCreate();
        if (SocketConnect(socket2, serverAddress2, serverPort2, 1000)) {
            Print("Connected to the second server");
        } else {
            int error = GetLastError();
            Print(" error : ", error);
            Sleep(1000);
        }
    }

    GetSymbolList1();
    GetSymbolList2();
}

void OnTick() {
    ConnectionAndTick();
}

void OnTimer() {
    ConnectionAndTick();
}

