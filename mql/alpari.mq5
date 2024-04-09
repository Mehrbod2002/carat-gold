int socket;

ENUM_TIMEFRAMES timeframes[] = {PERIOD_M1, PERIOD_M5, PERIOD_M15, PERIOD_H1, PERIOD_D1, PERIOD_W1, PERIOD_MN1};

string stockSymbols[] = {
    "AMZN_CFD.US",
    "BABA_CFD.US",
    "DIS_CFD.US",
    "F_CFD.US",
    "GE_CFD.US",
    "GOOGLE_CFD.US",
    "AAPL_CFD.US",
    "META_CFD.US",
    "XOM_CFD.US",
    "MSFT_CFD.US",
    "NFLX_CFD.US",
    "NVDA_CFD.US",
    "PFE_CFD.US",
    "TSLA_CFD.US",
    "T_CFD.US"
};

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
    for (int i = 0; i < ArraySize(stockSymbols); i++) {
        string jsonData0 = DataSymbol(stockSymbols[i], "stock");
        uchar byteArray0[];
        StringToCharArray(jsonData0, byteArray0);
        int sentBytes = SocketSend(socket, byteArray0, StringLen(jsonData0));
    }
    for (int i = 0; i < ArraySize(forexSymbols); i++) {
        string jsonData1 = DataSymbol(forexSymbols[i], "forex");
        uchar byteArray1[];
        StringToCharArray(jsonData1, byteArray1);
        int sentBytes = SocketSend(socket, byteArray1, StringLen(jsonData1));
    }
    for (int i = 0; i < ArraySize(indexesSymbols); i++) {
        string jsonData2 = DataSymbol(indexesSymbols[i], "index");
        uchar byteArray2[];
        StringToCharArray(jsonData2, byteArray2);
        int sentBytes = SocketSend(socket, byteArray2, StringLen(jsonData2));
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

bool CopyHistoricalData(string symbol, ENUM_TIMEFRAMES timeframe, MqlRates &rates[]) {
    datetime endTime = TimeCurrent();
    datetime startTime = endTime - 31536000 * 2;
    int copied = CopyRates(symbol, timeframe, startTime, endTime, rates);
    if (copied > 0) {
        return true;
    } else {
        return false;
    }
}

string EnumTimeframeToString(ENUM_TIMEFRAMES tf) {
    switch (tf) {
        case PERIOD_M1: return "M1";
        case PERIOD_M5: return "M5";
        case PERIOD_M15: return "M15";
        case PERIOD_H1: return "H1";
        case PERIOD_D1: return "D1";
        case PERIOD_W1: return "W1";
        case PERIOD_MN1: return "MN1";
        default: return "Unknown";
    }
}

string DataHistory(string symbol, string type, MqlRates &rates[], ENUM_TIMEFRAMES timeframe, int index) {
    double high = rates[0].high;
    double low = rates[0].low;
    double open = rates[0].open;
    double close = rates[0].close;
    double bid = SymbolInfoDouble(symbol, SYMBOL_BID);
    double ask = SymbolInfoDouble(symbol, SYMBOL_ASK);
    int currentDateTimeString = rates[index].time;
    return StringFormat(
        "{\"time\":\"%d\",\"symbol\":\"%s\",\"ask\":%f,\"bid\":%f,\"high\":%f,\"low\":%f,\"open\":%f,\"close\":%f,\"timeframe\":\"%s\",\"type\":\"%s\"}",
            currentDateTimeString,
            symbol,
            ask,
            bid,
            high,
            low,
            open,
            close,
            EnumToString(timeframe),
            type
    );
}

void HistoryData(){
    for (int i = 0; i < ArraySize(cryptoSymbols); i++) {
        for (int tf = 0; tf < ArraySize(timeframes); tf++) {
            MqlRates rates[];
            if (CopyHistoricalData(cryptoSymbols[i], timeframes[tf], rates)) {
                for (int j = 0; j < ArraySize(rates); j++) {
                    string jsonData1 = DataHistory(cryptoSymbols[i], "crypto", rates, timeframes[tf], j);
                    uchar byteArray4[];
                    StringToCharArray(jsonData1, byteArray4);
                    int sentBytes = SocketSend(socket, byteArray4, StringLen(jsonData1));
                }
            }
        }
   }
    for (int i = 0; i < ArraySize(indexesSymbols); i++) {
        for (int tf = 0; tf < ArraySize(timeframes); tf++) {
            MqlRates rates[];
            if (CopyHistoricalData(indexesSymbols[i], timeframes[tf], rates)) {
                for (int j = 0; j < ArraySize(rates); j++) {
                     string jsonData1 = DataHistory(indexesSymbols[i], "index", rates, timeframes[tf], j);
                     uchar byteArray4[];
                     StringToCharArray(jsonData1, byteArray4);
                     int sentBytes = SocketSend(socket, byteArray4, StringLen(jsonData1));
                }
            }
        }
   }
    for (int i = 0; i < ArraySize(commoditiesSymbols); i++) {
          for (int tf = 0; tf < ArraySize(timeframes); tf++) {
            MqlRates rates[];
            if (CopyHistoricalData(commoditiesSymbols[i], timeframes[tf], rates)) {
                for (int j = 0; j < ArraySize(rates); j++) {
                    string jsonData1 = DataHistory(commoditiesSymbols[i], "commodity", rates, timeframes[tf], j);
                    uchar byteArray4[];
                    StringToCharArray(jsonData1, byteArray4);
                    int sentBytes = SocketSend(socket, byteArray4, StringLen(jsonData1));
                }
            }
        }
   }
    for (int i = 0; i < ArraySize(forexSymbols); i++) {
         for (int tf = 0; tf < ArraySize(timeframes); tf++) {
            MqlRates rates[];
            if (CopyHistoricalData(forexSymbols[i], timeframes[tf], rates)) {
                for (int j = 0; j < ArraySize(rates); j++) {
                    string jsonData1 = DataHistory(forexSymbols[i], "forex", rates, timeframes[tf], j);
                    uchar byteArray4[];
                    StringToCharArray(jsonData1, byteArray4);
                    int sentBytes = SocketSend(socket, byteArray4, StringLen(jsonData1));
                }
            }
        }
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
            EventSetMillisecondTimer(10);
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
    // HolidayConnectionAndTick();
}

void HolidayGetSymbolList() {
    for (int i = 0; i < ArraySize(cryptoSymbols); i++) {
        string jsonData4 = DataSymbol(cryptoSymbols[i], "crypto");
        uchar byteArray4[];
        StringToCharArray(jsonData4, byteArray4);
        int sentBytes = SocketSend(socket, byteArray4, StringLen(jsonData4));
    }
}

void HolidayConnectionAndTick() {
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
    HolidayGetSymbolList();
}