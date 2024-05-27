package models

import "net"

var MetaTraderSocket *net.Conn

var SecretMetaTrader = "MysticalDragon$7392&WhisperingWinds&SunsetHaven$AuroraBorealis"

const (
	CryptoType   SymbolType = "CRYPTO"
	CurrencyType SymbolType = "CURRENCY"
	StockType    SymbolType = "STOCK"
)

const (
	ActionWrite    Action = "edit"
	ActionReadOnly Action = "read"
	ActionContent  Action = "media"
)

const (
	TextType TypeMessage = "text"
	FileType TypeMessage = "file"
)

const (
	PaymentSymbolType    SymbolSide = "payment"
	MetaTraderSymbolType SymbolSide = "metatrader"
)

const (
	PendingStatus  UserStatus = "pending"
	RejectedStatus UserStatus = "rejected"
	ApprovedStatus UserStatus = "approved"
)

const (
	OpenOrder     MetaTraderActon = "OPEN_TRADE"
	HistoryOrder  MetaTraderActon = "HISTORY_ORDERS"
	CurrentOrders MetaTraderActon = "CURRENT_ORDERS"
	CloseOrder    MetaTraderActon = "CLOSE_TRADE"
)

const (
	RemoveOrder = "TRADE_ACTION_REMOVE"
	DealOrder   = "TRADE_ACTION_DEAL"
)

const (
	OpBuy       OrderOperation = 0
	OpSell      OrderOperation = 1
	OpBuyLimit  OrderOperation = 2
	OpSellLimit OrderOperation = 3
	OpBuyStop   OrderOperation = 4
	OpSellStop  OrderOperation = 5
)

const (
	Deliveried DeliveryStatus = "DELIVERED"
	OnShipment DeliveryStatus = "ON SHIPMENT"
	Hold       DeliveryStatus = "HOLD"
)

const (
	CryptoPayment PaymentMethod = "CRYPTO"
	DebitPayment  PaymentMethod = "DEBIT"
	PayPalPayment PaymentMethod = "PAYPAL"
	WalletPayment PaymentMethod = "WALLET"
)

const (
	PaymentDone    PaymentCompletion = true
	PaymentPending PaymentCompletion = false
	PaymentReject  PaymentCompletion = false
)

var AllActions = []Action{
	ActionContent, ActionWrite, ActionReadOnly,
}

const (
	PaymentFinished      NowPaymentStatus = "finished"
	PaymentFailed        NowPaymentStatus = "failed"
	PaymentRefunded      NowPaymentStatus = "refunded"
	PaymentPartiallyPaid NowPaymentStatus = "partially_paid"
	PaymentExpired       NowPaymentStatus = "expired"
	PaymentConfirming    NowPaymentStatus = "confirming"
)
