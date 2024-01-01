package models

const (
	ActionRead               Action = "personal_access_data"
	ActionDocument           Action = "review_document"
	ActionActDocument        Action = "act_on_document"
	ActionReview             Action = "review_transaction"
	ActionActTransaction     Action = "act_on_transaction"
	ActionWrite              Action = "edit_user"
	ActionDelete             Action = "delete_user"
	ActionFreeUser           Action = "freeze_user"
	ActionUnfreezeUser       Action = "unfreeze_user"
	ActionTicketAccess       Action = "review_tickets"
	ActionActOnTickets       Action = "act_on_tickets"
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
	ActionGetReferal         Action = "referal"
)

const (
	TextType TypeMessage = "text"
	FileType TypeMessage = "file"
)
