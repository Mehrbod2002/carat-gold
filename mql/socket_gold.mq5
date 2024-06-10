int socket;

string DataSymbol(string symbolName) {
    double ask = SymbolInfoDouble(symbolName, SYMBOL_ASK);
    double bid = SymbolInfoDouble(symbolName, SYMBOL_BID);

    return StringFormat(
        "{\"symbol\":\"%s\",\"ask\":%f,\"bid\":%f}",
        symbolName, ask, bid
    );
}

string serverAddress = "52.0.87.119";
int serverPort = 5741;

bool EstablishSocketConnection() {
    if (!SocketIsConnected(socket)) {
        socket = SocketCreate();
        if (SocketConnect(socket, serverAddress, serverPort, 5000)) {
            Print("Connected to the server");
            return true;
        } else {
            Print("Connection failed with error: ", GetLastError());
            return false;
        }
    }
    return true;
}

void SendSymbolData(string symbolData) {
    uchar byteArray[];
    StringToCharArray(symbolData, byteArray);
    SocketSend(socket, byteArray, StringLen(symbolData));
}

int OnInit() {
    if (EstablishSocketConnection()) {
        return INIT_SUCCEEDED;
    } else {
        return INIT_FAILED;
    }
}

void OnDeinit(const int reason) {
    SocketClose(socket);
}

void ConnectionAndTick() {
    if (!EstablishSocketConnection()) {
        return;
    }
    string symbolData = DataSymbol("XAUUSD");
    SendSymbolData(symbolData);
}

void OnTick() {
    ConnectionAndTick();
}