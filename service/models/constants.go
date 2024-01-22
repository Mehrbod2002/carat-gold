package models

import "net"

var MetaTraderSocket *net.Conn

const (
	ActionRead               Action = "personal_access_data"
	ActionDocument           Action = "review_document"
	ActionActDocument        Action = "act_on_document"
	ActionReview             Action = "review_transaction"
	ActionActTransaction     Action = "act_on_transaction"
	ActionFreeUser           Action = "freeze_user"
	ActionSetUser            Action = "define_user"
	ActionChat               Action = "chat"
	ActionSupport            Action = "support"
	ActionSendNotification   Action = "send_notification"
	ActionChangeTradeReport  Action = "trade_report"
	ActionReportsTransaction Action = "reports_transactions"
	ActionSetTransactions    Action = "transactions_data"
	ActionGeneralDataView    Action = "general_data_view"
	ActionGeneralDataEdit    Action = "general_data_edit"
	ActionReviewMessage      Action = "review_message"
	ActionOpenChat           Action = "open_message"
	ActionSetPermission      Action = "set_permissions"
	ActionEditUser           Action = "edit_user"
	ActionDeleteUser         Action = "delete_user"
	ActionMetaTrader         Action = "metatrader"
)

const (
	TextType TypeMessage = "text"
	FileType TypeMessage = "file"
)

const (
	ForexType     SymbolType = "forex"
	CommodityType SymbolType = "commodity"
	IndexesType   SymbolType = "index"
	CryptoType    SymbolType = "crypto"
	StockType     SymbolType = "stock"
	FiatType      SymbolType = "fiat"
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

var AllActions = []Action{
	ActionRead,
	ActionDocument,
	ActionActDocument,
	ActionReview,
	ActionActTransaction,
	ActionFreeUser,
	ActionSetUser,
	ActionChat,
	ActionSupport,
	ActionSendNotification,
	ActionChangeTradeReport,
	ActionReportsTransaction,
	ActionSetTransactions,
	ActionGeneralDataView,
	ActionGeneralDataEdit,
	ActionReviewMessage,
	ActionOpenChat,
	ActionSetPermission,
	ActionEditUser,
	ActionDeleteUser,
	ActionMetaTrader,
}
