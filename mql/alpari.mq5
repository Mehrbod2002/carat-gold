int socket;

ENUM_TIMEFRAMES timeframes[] = {PERIOD_M1, PERIOD_M5, PERIOD_M15, PERIOD_H1, PERIOD_D1, PERIOD_W1, PERIOD_MN1};
string forexSymbols[] = {
    "EURUSD",
    "GBPUSD",
    "USDJPY",
    "EURGBP",
    "EURJPY",
    "GBPJPY",
    "USDCAD",
    "CADCHF",
    "USDCHF",
    "CHFJPY",
    "EURCHF",
    "GBPCHF",
    "CADJPY",
    "EURCAD",
    "GBPCAD",
    "AUDUSD",
    "NZDUSD",
    "AUDCHF",
    "CHFPLN",
    "NZDCHF",
    "CHFSGD",
    "AUDCAD",
    "AUDJPY",
    "AUDNZD",
    "EURAUD",
    "EURNOK",
    "EURNZD",
    "EURSEK",
    "GBPAUD",
    "GBPNZD",
    "NZDCAD",
    "NZDJPY",
    "USDNOK",
    "USDSEK",
    "AUDDKK",
    "EURHUF",
    "EURMXN",
    "EURPLN",
    "EURTRK",
    "EURZAR",
    "GBPNOK",
    "GBPPLN",
    "GBPSEK",
    "NOKSEK",
    "PLNJPY",
    "USDMXN",
    "USDHUF",
    "USDPLN",
    "USDTRY",
    "USDZAR",
    "EURRUB",
    "USDRUB",
    "USDILS",
    "USDCNH",
    "GBPZAR",
    "AUDPLN",
    "AUDSGD",
    "EURSGD",
    "GBPSGD",
    "NZDSGD",
    "USDDKK",
    "SGDJPY",
    "USDSGD",
    "EURHKD",
    "GBPDKK",
    "USDHKD",
    "EURDKK",
    "EURCZK",
    "USDCZK",
    "USDTHB"
};

string cryptoSymbols[] = {
    "CARDANO", "BAT", "BITCOINCASH", "BITCOIN", "DOGECOIN", "POLKADOT", "DASH", "EOS",
    "ETHCLASSIC", "ETHEREUM", "FILECOIN", "IOTA", "CHAINLINK", "AAVE", "LITECOIN", "ZCASH",
    "POLYGON", "NEO", "SOLANA", "SUSHISWAP", "THETA", "TRON", "UNISWAP", "VECHAIN", "STELLAR",
    "MONERO", "XRP", "TEZOS"
};

string indexesSymbols[] = {
    "EURO50", "FRANCE40", "GERMANY40", "AUS200", "JAPAN225", "USSPX500", "UK100", "US30", 
    "USNDAQ100", "SPAIN35", "CHINA50", "CHINAHshar", "SWISS20", "FRANCE120", "HONGKONG50", 
    "GERMANY50", "GERTECH30", "HOLLAND25"
};

string commoditiesSymbols[] = {
    "GOLD", "GOLDoz", "GOLDgr", "GOLDEURO", "SILVER", "SILVEREURO", "PLATINUM", "PALLADIUM",
    "ALUMINUM", "COPPER", "LEAD", "ZINC", "WTI", "BRENT", "NAT.GAS"
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

void HolidayGetSymbolList() {
    for (int i = 0; i < ArraySize(cryptoSymbols); i++) {
        string jsonData4 = DataSymbol(cryptoSymbols[i], "crypto");
        uchar byteArray4[];
        StringToCharArray(jsonData4, byteArray4);
        int sentBytes = SocketSend(socket, byteArray4, StringLen(jsonData4));
    }
}

void GetSymbolList() {
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

void ReadFromSocket() {
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
   char   rsp[1024];
   string result;
   uint   timeout_check=GetTickCount()+100;

   do {
      uint len=SocketIsReadable(socket);
      if(len > 0) {
        int rsp_len;
        rsp_len=SocketRead(socket,rsp,len,100);
        if(rsp_len>0) {
            result+=CharArrayToString(rsp,0,rsp_len);

            ProcessData(result);
            return;
        }
      }
   } while(GetTickCount()<timeout_check && !IsStopped());
}

string serverAddress = "185.202.113.18";
int serverPort = 8081;

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

void ConnectionAndTick(bool holiday) {
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
    if (holiday) {
        HolidayGetSymbolList();
    } else {
        GetSymbolList();
    }
}

void OnTick() {
    //ConnectionAndTick(false);
}

void OnTimer() {
    //ConnectionAndTick(true);
    ReadFromSocket();
}

void ProcessData(const string &data) {
    string components[];
    int count = StringSplit(data, '|', components);

    string requestID = components[0];
    if (count > 1) {
        if (components[1] == "OPEN_TRADE") {
            string symbol = components[2];
            double price = StringToDouble(components[3]);
            
            if (count >= 12) {
                int option = StringToInteger(components[4]);
                double volume = StringToDouble(components[5]);
                int slipppage = StringToInteger(components[6]);
                int magic_number = StringToInteger(components[7]);
                double stoploss = StringToDouble(components[8]);
                double takeprofit = StringToDouble(components[9]);
                string comment = components[10];
                int magic = StringToInteger(components[11]);
                datetime expiration = StringToTime(components[12]);

                OpenTrade(requestID, option, magic_number, symbol, price, volume, slipppage, stoploss, takeprofit, comment, magic, expiration);
            } else {
                SendError("INVALID", requestID);
                return;
            }
        } else if (components[0] == "HISTORY_ORDERS") {
            CreateJsonString(requestID);
            return;
        } else if (components[0] == "CLOSE_TRADE") {
            string symbol = components[1];
            int ticket = StringToInteger(components[4]);
            CloseTrade(requestID, symbol, ticket);
        } else if (components[0] == "CURRENT_ORDERS") {
            CurrentCreateJsonString(requestID);
            return;
        } else {
            SendError("INVALID", requestID);
            return;
        }
    } else {
        SendError("INVALID", requestID);
        return;
    }
}

void OpenTrade(string requestID ,int opt, int magic_number,string symbol, double price, double volume, int slipppage, 
    double stoploss, double takeprofit, string comment, int magic, datetime expiration) {

    MqlTradeRequest request={};
    MqlTradeResult result={};
    
    request.action=TRADE_ACTION_DEAL;
    request.magic=magic_number;
    request.symbol=symbol;
    request.volume=volume;
    request.sl=stoploss;
    request.tp=takeprofit;
    request.type=opt;
    request.price=price;
    if (OrderSend(request,result) > 0) {
        SendData("", requestID);
    }

    SendError("", requestID);
}

void CloseTrade(string requestID, string symbol, int ticket) {

    MqlTradeRequest request={};
    MqlTradeResult result={};
    
    request.action = TRADE_ACTION_REMOVE;
    request.symbol = symbol;
    request.order = ticket;
    if (OrderSend(request,result) > 0) {
        SendData("", requestID);
    }

    SendError("", requestID);
}

bool IsTradeRunning(ulong ticket) {
    int orderSelectResult = OrderSelect(ticket);
    
    if (orderSelectResult) {
        int orderStatus = OrderGetInteger(ORDER_STATE);
        
        return (orderStatus == ORDER_STATE_FILLED || orderStatus == ORDER_STATE_PARTIAL);
    } else {
        return false;
    }
}

void CreateJsonString(string requestID) {
    string jsonString = "[";

    int totalOrders = HistoryOrdersTotal();
    for (int i = 0; i < totalOrders; i++) {
        ulong ticket = HistoryOrderGetTicket(i);
        double openPrice = HistoryOrderGetDouble(ticket, ORDER_PRICE_OPEN);
        double lots = HistoryOrderGetDouble(ticket, ORDER_VOLUME_CURRENT);
        int type = HistoryOrderGetInteger(ticket, ORDER_TYPE);

        string symbol = HistoryOrderGetString(ticket, ORDER_SYMBOL);

        string tradeType;
        if (type == ORDER_TYPE_BUY) {
            tradeType = "Buy";
        } else if (type == ORDER_TYPE_BUY_LIMIT) {
            tradeType = "Buy Limit";
        } else if (type == ORDER_TYPE_BUY_STOP) {
            tradeType = "Buy";
        } else if (type == ORDER_TYPE_SELL) {
            tradeType = "Sell";
        } else if (type == ORDER_TYPE_SELL_LIMIT) {
            tradeType = "Sell limit";
        } else if (type == ORDER_TYPE_SELL_STOP) {
            tradeType = "Sell Stop";
        }

        bool isOpen = IsTradeRunning(ticket);

        string tradeInfo = StringFormat(
            "{\"Ticket\":%d,\"Symbol\":\"%s\",\"Type\":\"%s\",\"Lots\":%f,\"OpenPrice\":%f,\"IsOpen\":%s}",
            ticket, symbol, tradeType, lots, openPrice, isOpen ? "true" : "false"
        );

        jsonString += tradeInfo;

        if (i < totalOrders - 1) {
            jsonString += ",";
        }
    }

    jsonString += "]";
    SendData(jsonString, requestID);
    return;
}

void SendData(string data, string requestID) {
    string response = StringFormat(
      "{\"status\":%s,\"data\":\"%s\",\"id\":\"%s\"}",
        "\"true\"", data, requestID
    );

    uchar byteArray[];
    StringToCharArray(response, byteArray);
    SocketSend(socket, byteArray, StringLen(response));
    return;
}

void SendError(string data, string requestID) {
    string errorDescription;
    
    if (data == "") {
        int lastError = GetLastError();
        errorDescription = IntegerToString(lastError);
    } else {
        errorDescription = data;
    }

    string response = StringFormat(
      "{\"status\":%s,\"data\":\"%s\",\"id\":\"%s\"}",
        "\"false\"", errorDescription, requestID
    );

    uchar byteArray[];
    StringToCharArray(response, byteArray);
    SocketSend(socket, byteArray, StringLen(response));
    Print(response);
    return;
}

void CurrentCreateJsonString(string requestID) {
    string jsonString = "[";

    int totalOrders = OrdersTotal();
    for (int i = 0; i < totalOrders; i++) {
        ulong ticket = HistoryOrderGetTicket(i);
        double openPrice = HistoryOrderGetDouble(ticket, ORDER_PRICE_OPEN);
        double lots = HistoryOrderGetDouble(ticket, ORDER_VOLUME_CURRENT);
        int type = HistoryOrderGetInteger(ticket, ORDER_TYPE);

        string symbol = HistoryOrderGetString(ticket, ORDER_SYMBOL);

        string tradeType;
        if (type == ORDER_TYPE_BUY) {
            tradeType = "ORDER_TYPE_BUY";
        } else if (type == ORDER_TYPE_BUY_LIMIT) {
            tradeType = "ORDER_TYPE_BUY_LIMIT";
        } else if (type == ORDER_TYPE_BUY_STOP) {
            tradeType = "ORDER_TYPE_BUY_STOP";
        } else if (type == ORDER_TYPE_SELL) {
            tradeType = "ORDER_TYPE_SELL";
        } else if (type == ORDER_TYPE_SELL_LIMIT) {
            tradeType = "ORDER_TYPE_SELL_LIMIT";
        } else if (type == ORDER_TYPE_SELL_STOP) {
            tradeType = "ORDER_TYPE_SELL_STOP";
        }

        bool isOpen = true;

        string tradeInfo = StringFormat(
            "{\"Ticket\":%d,\"Symbol\":\"%s\",\"Type\":\"%s\",\"Lots\":%f,\"OpenPrice\":%f,\"IsOpen\":%s}",
            ticket, symbol, tradeType, lots, openPrice, isOpen ? "true" : "false"
        );

        jsonString += tradeInfo;

        if (i < totalOrders - 1) {
            jsonString += ",";
        }
    }

    jsonString += "]";
    SendData(jsonString, requestID);
    return;
}
