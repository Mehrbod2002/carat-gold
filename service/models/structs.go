package models

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CrispWebhookPayload struct {
	Payload struct {
		Data struct {
			Content     string        `json:"content"`
			Fingerprint float64       `json:"fingerprint"`
			From        string        `json:"from"`
			Mentions    []interface{} `json:"mentions"`
			Origin      string        `json:"origin"`
			SessionID   string        `json:"session_id"`
			Stamped     bool          `json:"stamped"`
			Timestamp   float64       `json:"timestamp"`
			Type        string        `json:"type"`
			User        struct {
				Nickname string `json:"nickname"`
				UserID   string `json:"user_id"`
			} `json:"user"`
			WebsiteID string `json:"website_id"`
		} `json:"data"`
		Event     string  `json:"event"`
		Timestamp float64 `json:"timestamp"`
		WebsiteID string  `json:"website_id"`
	} `json:"payload"`
}

type GeneralData struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	AedUsd float64            `bson:"aedusd" json:"aedusd"`
}

type MetaTraderAccounts struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Server   string             `bson:"server" json:"server"`
	Login    string             `bson:"login" json:"login"`
	Passowrd string             `bson:"password" json:"password"`
}

type RequestMetaTraderAccounts struct {
	Server   string `bson:"server" json:"server"`
	Login    string `bson:"login" json:"login"`
	Passowrd string `bson:"password" json:"password"`
}

type RequestSetCallCenter struct {
	WhatsApp         *string `bson:"whatsapp" json:"whatsapp"`
	LiveChat         *string `bson:"live" json:"live"`
	Telegram         *string `bson:"telegram" json:"telegram"`
	WebsiteReference *string `bson:"website_reference" json:"website_reference"`
	Email            *string `bson:"email" json:"email"`
	PhoneComapny     *string `bson:"phone" json:"phone"`
}

type CallCenterDatas struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	WhatsApp         string             `bson:"whatsapp" json:"whatsapp"`
	LiveChat         *string            `bson:"live" json:"live"`
	Telegram         string             `bson:"telegram" json:"telegram"`
	WebsiteReference string             `bson:"website_reference" json:"website_reference"`
	Email            string             `bson:"email" json:"email"`
	PhoneComapny     string             `bson:"phone" json:"phone"`
}

type DebitCard struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Vat       float64            `bson:"vat" json:"vat"`
	Address   string             `bson:"address" json:"address"`
	Access    string             `bson:"access" json:"access"`
	Token     string             `bson:"token" json:"token"`
	WhoDefine string             `bson:"who_define" json:"who_define"`
	Disable   bool               `bson:"disable" json:"disable"`
}

type PayPal struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Vat       float64            `bson:"vat" json:"vat"`
	WhoDefine string             `bson:"who_define" json:"who_define"`
	Address   string             `bson:"address" json:"address"`
	Access    string             `bson:"access" json:"access"`
	Token     string             `bson:"token" json:"token"`
	Disable   bool               `bson:"disable" json:"disable"`
}

type DefaultPayment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Vat       float64            `bson:"vat" json:"vat"`
	WhoDefine string             `bson:"who_define" json:"who_define"`
	Address   string             `bson:"address" json:"address"`
	Access    string             `bson:"access" json:"access"`
	Token     string             `bson:"token" json:"token"`
	Disable   bool               `bson:"disable" json:"disable"`
}

type RequestSetPayment struct {
	ID        *string `bson:"_id,omitempty" json:"_id"`
	Side      string  `bson:"side" json:"side"`
	Vat       float64 `bson:"vat" json:"vat"`
	WhoDefine string  `bson:"who_define" json:"who_define"`
	Address   string  `bson:"address" json:"address"`
	Access    string  `bson:"access" json:"access"`
	Token     string  `bson:"token" json:"token"`
}

type GooglePlay struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Vat       float64            `bson:"vat" json:"vat"`
	WhoDefine string             `bson:"who_define" json:"who_define"`
	Address   string             `bson:"address" json:"address"`
	Access    string             `bson:"access" json:"access"`
	Token     string             `bson:"token" json:"token"`
	Disable   bool               `bson:"disable" json:"disable"`
}

type ApplePlay struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Vat       float64            `bson:"vat" json:"vat"`
	WhoDefine string             `bson:"who_define" json:"who_define"`
	Address   string             `bson:"address" json:"address"`
	Access    string             `bson:"access" json:"access"`
	Token     string             `bson:"token" json:"token"`
	Disable   bool               `bson:"disable" json:"disable"`
}

type Crypto struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Coin      string             `bson:"coin" json:"coin"`
	Vat       float64            `bson:"vat" json:"vat"`
	Wallet    string             `bson:"wallet" json:"wallet"`
	WhoDefine string             `bson:"who_define" json:"who_define"`
	Address   string             `bson:"address" json:"address"`
	Access    string             `bson:"access" json:"access"`
	Token     string             `bson:"token" json:"token"`
	Disable   bool               `bson:"disable" json:"disable"`
}

type FANDQ struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Question  string             `bson:"question" json:"question"`
	Answer    string             `bson:"answer" json:"answer"`
	WhoDefine string             `bson:"who_define" json:"who_define"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type Invoice struct {
	ID               string    `json:"id"`
	OrderID          string    `json:"order_id"`
	OrderDescription string    `json:"order_description"`
	PriceAmount      string    `json:"price_amount"`
	PriceCurrency    string    `json:"price_currency"`
	PayCurrency      string    `json:"pay_currency"`
	IPNCallbackURL   string    `json:"ipn_callback_url"`
	InvoiceURL       string    `json:"invoice_url"`
	SuccessURL       string    `json:"success_url"`
	CancelURL        string    `json:"cancel_url"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type PaymentResponse struct {
	PaymentID              string      `json:"payment_id"`
	PaymentStatus          string      `json:"payment_status"`
	PayAddress             string      `json:"pay_address"`
	PriceAmount            float64     `json:"price_amount"`
	PriceCurrency          string      `json:"price_currency"`
	PayAmount              float64     `json:"pay_amount"`
	PayCurrency            string      `json:"pay_currency"`
	OrderID                string      `json:"order_id"`
	OrderDescription       string      `json:"order_description"`
	IPNCallbackURL         string      `json:"ipn_callback_url"`
	CreatedAt              string      `json:"created_at"`
	UpdatedAt              string      `json:"updated_at"`
	PurchaseID             string      `json:"purchase_id"`
	AmountReceived         *float64    `json:"amount_received"`
	PayinExtraID           interface{} `json:"payin_extra_id"`
	SmartContract          string      `json:"smart_contract"`
	Network                string      `json:"network"`
	NetworkPrecision       int         `json:"network_precision"`
	TimeLimit              interface{} `json:"time_limit"`
	BurningPercent         interface{} `json:"burning_percent"`
	ExpirationEstimateDate string      `json:"expiration_estimate_date"`
}

type Purchased struct {
	Product        []primitive.ObjectID `bson:"product" json:"product"`
	StatusDelivery DeliveryStatus       `bson:"status_delivery" json:"status_delivery"`
	PaymentStatus  UserStatus           `bson:"payment_status" json:"payment_status"`
	PaymentMethd   PaymentMethod        `bson:"payment_method" json:"payment_method"`
	CreatedAt      time.Time            `bson:"created_at" json:"created_at"`
	CreatePayment  time.Time            `bson:"created_payment" json:"created_payment"`
	OrderID        string               `bson:"order_id" json:"order_id"`
}

type Wallet struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID     primitive.ObjectID `bson:"user_id,omitempty" json:"user_id"`
	BalanceUSD float64            `bson:"balance" json:"balance"`
	Purchased  []Purchased        `bson:"purchased" json:"purchased"`
}

type Address struct {
	Label   string `bson:"label" json:"label"`
	Country string `bson:"country" json:"country"`
	City    string `bson:"city" json:"city"`
	Region  string `bson:"region" json:"region"`
	Address string `bson:"address" json:"address"`
}

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Avatar           string             `bson:"avatar" json:"avatar"`
	FirstName        string             `bson:"first_name" json:"first_name"`
	LastName         string             `bson:"last_name" json:"last_name"`
	Email            string             `bson:"email" json:"email"`
	PhoneNumber      string             `bson:"phone" json:"phone"`
	Password         string             `bson:"password" json:"password"`
	PhoneVerified    bool               `bson:"phone_verified" json:"phone_verified"`
	OtpCode          *int               `bson:"otp_code" json:"otp_code"`
	Address          []Address          `bson:"address" json:"address"`
	ResetToken       string             `bson:"reset_token" json:"reset_token"`
	ResetTokenValid  time.Time          `bson:"reset_token_valid" json:"reset_token_valid"`
	ChangePhone      bool               `bson:"change_phone" json:"change_phone"`
	ExchangeMobile   string             `bson:"exchange_phone" json:"exchange_phone"`
	Freeze           bool               `bson:"freeze" json:"freeze"`
	Currency         string             `bson:"currency" json:"currency"`
	ChatList         []string           `bson:"chat_list" json:"chat_list"`
	Permissions      Permission         `bson:"permissions" json:"permissions"`
	OtpValid         time.Time          `bson:"otp_valid" json:"otp_valid"`
	ReTryOtp         int                `bson:"retry_otp" json:"retry_otp"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UserVerified     bool               `bson:"user_verified" json:"user_verified"`
	StatusString     UserStatus         `bson:"user_status" json:"user_status"`
	Reason           string             `bson:"reason" json:"reason"`
	IsSupportOrAdmin bool               `bson:"support_or_admin" json:"support_or_admin"`
	Wallet           Wallet             `bson:"wallet" json:"wallet"`
	FcmToken         string             `bson:"fcm_token" json:"fcm_token"`
	RefreshToken     string             `bson:"refresh_token" json:"-"`
}

type FeedBacks struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID   primitive.ObjectID `bson:"user_id" json:"user_id,omitempty"`
	FeedBack string             `bson:"feedback" json:"feedback"`
}

type RequestEdit struct {
	FirstName *string    `json:"first_name"`
	LastName  *string    `json:"last_name"`
	Address   *[]Address `json:"address"`
	Email     *string    `json:"email"`
	Phone     *string    `json:"phone"`
}

type Permission struct {
	Actions []Action `bson:"actions"`
}

type Socket struct {
	ResponseTo    User      `json:"user"`
	Trigger       string    `json:"trigger"`
	Validate      bool      `json:"validate"`
	Message       string    `json:"message"`
	Messages      []Message `json:"messages"`
	SingleMessage Message   `json:"single_message"`
}

type Client struct {
	Conn *websocket.Conn
	User *User
}

type NotificationAdmin struct {
	Type    string   `json:"notification"`
	Subject string   `json:"subject"`
	Message string   `json:"message"`
	Names   []string `json:"names"`
}

type Claims struct {
	ID          string    `json:"_id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	jwt.StandardClaims
}

type Message struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Sender           primitive.ObjectID `json:"sender" bson:"sender"`
	Receiver         primitive.ObjectID `json:"receiver" bson:"receiver"`
	SenderUsername   string             `json:"sender_username" bson:"sender_username"`
	ReceiverUsername string             `json:"receiver_username" bson:"receiver_username"`
	Content          string             `json:"content" bson:"content"`
	Seen             bool               `json:"seen" bson:"seen"`
	MessageType      TypeMessage        `json:"message_type" bson:"message_type"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
}

type UserMessage struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Content     string      `json:"message"`
	TypeMessage TypeMessage `json:"message_type"`
}

type Currency struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Currency  string             `bson:"currency,omitempty" json:"currency"`
	WhoDefine string             `bson:"who_define" json:"who_define"`
	Active    bool               `bson:"active" json:"active"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type Products struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name        string             `bson:"name,omitempty" json:"name"`
	Description string             `bson:"description" json:"description"`
	Images      []Image            `bson:"images" json:"images"`
	Width       float64            `bson:"width" json:"width"`
	Length      float64            `bson:"length" json:"length"`
	SubTitle    string             `bson:"sub_title" json:"sub_title"`
	Faq         string             `bson:"faq" json:"faq"`
	Answer      string             `bson:"answer" json:"answer"`
	WeightOZ    float64            `bson:"weight_oz" json:"weight_oz"`
	WeightGramm float64            `bson:"weight_gramm" json:"weight_gramm"`
	Purity      float64            `bson:"purity" json:"purity"`
	Percentage  float64            `bson:"percentage" json:"percentage"`
	Hide        bool               `bson:"hide" json:"hide"`
	Limited     bool               `bson:"limited" json:"limited"`
	WhoDefine   string             `bson:"who_define" json:"who_define"`
	Amount      int                `bson:"amount" json:"amount"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

type Symbol struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	SymbolName string             `bson:"name" json:"name"`
	SymbolType SymbolType         `bson:"type" json:"type"`
	SymbolSide SymbolSide         `bson:"side" json:"side"`
	Images     []Image            `bson:"images" json:"images"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

type DeliveryMethods struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Description   string             `bson:"description,omitempty" json:"description"`
	WhoDefine     string             `bson:"who_define" json:"who_define"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	Title         string             `bson:"title,omitempty" json:"title"`
	Fee           float64            `bson:"fee" json:"fee"`
	SubTitle      string             `bson:"sub_title" json:"sub_title"`
	EstimatedTime string             `bson:"estimated_time,omitempty" json:"estimated_time"`
	TimeProvided  bool               `bson:"time_provided,omitempty" json:"time_provided"`
	Disable       bool               `bson:"disable" json:"disable"`
}

type Transaction struct {
	ID                primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	OrderID           string               `bson:"order_id" json:"order_id,omitempty"`
	ProductIDs        []primitive.ObjectID `bson:"product_ids" json:"product_ids"`
	UserID            primitive.ObjectID   `bson:"user_id" json:"user_id,omitempty"`
	CreatedAt         time.Time            `bson:"created_at" json:"created_at"`
	PaymentMethod     PaymentMethod        `bson:"payment_method" json:"payment_method"`
	DeliveryMethod    DeliveryMethod       `bson:"delivery_method" json:"delivery_method"`
	PaymentCompletion PaymentCompletion    `bson:"payment_completion" json:"payment_completion"`
	StatusDelivery    DeliveryStatus       `bson:"status_delivery" json:"status_delivery"`
	Symbol            string               `bson:"symbol" json:"symbol"`
	TotalPrice        float64              `bson:"total_price" json:"total_price"`
	Vat               float64              `bson:"vat" json:"vat"`
	PaymentStatus     UserStatus           `bson:"payment_status" json:"payment_status"`
	IsDebit           bool                 `bson:"is_debit" json:"is_debit"`
	ExternalData      map[string]string    `bson:"external_data" json:"external_data"`
	MetatraderID      string               `bson:"metatrader_id" json:"metatrader_id"`
}

type RequestSetDefineUser struct {
	UserID      *string     `json:"user_id"`
	FirstName   *string     `json:"name"`
	LastName    *string     `json:"last_name"`
	Email       *string     `json:"email"`
	Password    *string     `json:"password"`
	Phone       *string     `json:"phone"`
	PhoneVerify *bool       `json:"phone_verified"`
	IsSupport   *bool       `json:"support"`
	Freeze      *bool       `json:"freeze"`
	Status      *UserStatus `json:"status"`
	Reason      *string     `json:"reason"`
	Permissions *Permission `json:"permissions"`
	BalanceUSD  *float64    `json:"balance"`
}

type LoginDataStep1 struct {
	Phone string `json:"phone"`
}

type LoginDataStep2 struct {
	Phone string `json:"phone"`
	Otp   *int   `json:"otp"`
}

type SendOTP struct {
	PhoneNumber string `json:"phone" binding:"required"`
}

type RegisterRequest struct {
	PhoneNumber string `json:"phone" binding:"required"`
	OtpCode     *int   `json:"otp_code" binding:"required"`
}

type Shot struct {
	Shot string `json:"shot" bson:"shot"`
}

type DocumentShots struct {
	Back  Shot `json:"back" bson:"back"`
	Front Shot `json:"front" bson:"front"`
}

type Documents struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id,omitempty"`
	Documents DocumentShots      `bson:"documents" json:"documents"`
}

type DataMeta struct {
	Time      string  `json:"time" bson:"time"`
	Symbol    string  `json:"symbol" bson:"symbol"`
	Ask       float64 `json:"ask" bson:"ask"`
	Bid       float64 `json:"bid" bson:"bid"`
	High      float64 `json:"high" bson:"high"`
	Low       float64 `json:"low" bson:"low"`
	Open      float64 `json:"open" bson:"open"`
	Close     float64 `json:"close" bson:"close"`
	Type      string  `json:"type" bson:"type"`
	Timeframe string  `json:"timeframe" bson:"timeframe"`
}

type OpenTrade struct {
	Time      string  `json:"time" bson:"time"`
	Symbol    string  `json:"symbol" bson:"symbol"`
	Ask       float64 `json:"ask" bson:"ask"`
	Bid       float64 `json:"bid" bson:"bid"`
	High      float64 `json:"high" bson:"high"`
	Low       float64 `json:"low" bson:"low"`
	Open      float64 `json:"open" bson:"open"`
	Close     float64 `json:"close" bson:"close"`
	Type      string  `json:"type" bson:"type"`
	Timeframe string  `json:"timeframe" bson:"timeframe"`
}

type OperationMetaTrader struct {
	OP_BUY       string `json:"OP_BUY"`
	OP_SELL      string `json:"OP_SELL"`
	OP_BUYLIMIT  string `json:"OP_BUYLIMIT"`
	OP_SELLLIMIT string `json:"OP_SELLLIMIT"`
	OP_BUYSTOP   string `json:"OP_BUYSTOP"`
	OP_SELLSTOP  string `json:"OP_SELLSTOP"`
	OP_BALANCE   string `json:"OP_BALANCE"`
	OP_CREDIT    string `json:"OP_CREDIT"`
	OP_BUYTO     string `json:"OP_BUYTO"`
	OP_SELLTO    string `json:"OP_SELLTO"`
}

type RequestSetSymbol struct {
	Name  string     `json:"symbol_name"`
	Image string     `json:"image"`
	Type  SymbolType `json:"type"`
}

type RequestSetGeneralData struct {
	AedUsd float64 `json:"aedusd"`
}

type MetaTraderAdmin struct {
	Data   string    `json:"data"`
	ID     int       `json:"id"`
	Time   time.Time `json:"time"`
	Status string    `json:"status"`
}

type RequestSetTrade struct {
	SymbolName string   `json:"name"`
	Volumn     float64  `json:"volumn"`
	Operation  int      `json:"operation"`
	StopLoss   *float64 `json:"stoploss"`
	TakeProfit *float64 `json:"takeprofit"`
	Comment    *string  `json:"comment"`
	Deviation  *float64 `json:"deviation"`
	Stoplimit  *float64 `json:"stop_limit"`
}

type Image struct {
	PhotoID primitive.ObjectID `bson:"image" json:"image"`
}

type RequestSetProduct struct {
	ProductID   *string   `json:"product_id"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	Images      *[]string `json:"images"`
	Width       *float64  `json:"width"`
	Length      *float64  `json:"length"`
	WeightOZ    *float64  `json:"weight_oz"`
	WeightGramm *float64  `json:"weight_gramm"`
	Purity      *float64  `json:"purity"`
	PriceGramm  *float64  `json:"price_gramm"`
	Percentage  *float64  `json:"percentage"`
	Hide        *bool     `bson:"hide" json:"hide"`
	Limited     *bool     `bson:"limited" json:"limited"`
	Amount      *int      `bson:"amount" json:"amount"`
	SubTitle    *string   `bson:"sub_title" json:"sub_title"`
	Faq         *string   `bson:"faq" json:"faq"`
	Answer      *string   `bson:"answer" json:"answer"`
}

type RequestSetDeliveryMethod struct {
	DeliveryID    *string `json:"delivery_id"`
	Title         string  `bson:"title,omitempty" json:"title"`
	Fee           float64 `bson:"fee" json:"fee"`
	Description   string  `bson:"description" json:"description"`
	EstimatedTime string  `bson:"estimated_time" json:"estimated_time"`
	TimeProvided  bool    `bson:"time_provided" json:"time_provided"`
	SubTitle      string  `bson:"sub_title" json:"sub_title"`
	Disable       bool    `bson:"disable" json:"disable"`
}

type RequestSetCancelTrade struct {
	Ticket string `json:"ticket_id"`
}

type RequestSetFANDQ struct {
	ID       *string `json:"_id"`
	Question string  `bson:"question" json:"question"`
	Answer   string  `bson:"answer" json:"answer"`
}

type PaymentCallBack struct {
	PaymentID          int64            `json:"payment_id"`
	ParentPaymentID    int64            `json:"parent_payment_id"`
	InvoiceID          interface{}      `json:"invoice_id"`
	PaymentStatus      NowPaymentStatus `json:"payment_status"`
	PayAddress         string           `json:"pay_address"`
	UpdatedAt          time.Time        `json:"updated_at"`
	PayinExtraID       interface{}      `json:"payin_extra_id"`
	PriceAmount        float64          `json:"price_amount"`
	PriceCurrency      string           `json:"price_currency"`
	PayAmount          float64          `json:"pay_amount"`
	ActuallyPaid       float64          `json:"actually_paid"`
	ActuallyPaidAtFiat float64          `json:"actually_paid_at_fiat"`
	PayCurrency        string           `json:"pay_currency"`
	OrderID            interface{}      `json:"order_id"`
	OrderDescription   interface{}      `json:"order_description"`
	PurchaseID         string           `json:"purchase_id"`
	OutcomeAmount      float64          `json:"outcome_amount"`
	OutcomeCurrency    string           `json:"outcome_currency"`
	PaymentExtraIDs    interface{}      `json:"payment_extra_ids"`
	Fee                Fee              `json:"fee"`
}

type Fee struct {
	Currency      string  `json:"currency"`
	DepositFee    float64 `json:"depositFee"`
	WithdrawalFee float64 `json:"withdrawalFee"`
	ServiceFee    float64 `json:"serviceFee"`
}

type TwilioVerifyVerification struct {
	To      string `json:"to"`
	Channel string `json:"channel"`
}

type TwilioVerifyVerificationCheck struct {
	To   string `json:"to"`
	Code string `json:"code"`
}

type UserDeliveryMethods struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Description   string             `bson:"description,omitempty" json:"description"`
	WhoDefine     string             `bson:"who_define" json:"-"`
	CreatedAt     time.Time          `bson:"created_at" json:"-"`
	Title         string             `bson:"title,omitempty" json:"title"`
	Fee           float64            `bson:"fee" json:"fee"`
	SubTitle      string             `bson:"sub_title" json:"sub_title"`
	EstimatedTime string             `bson:"estimated_time,omitempty" json:"estimated_time"`
	TimeProvided  bool               `bson:"time_provided,omitempty" json:"time_provided"`
	Disable       bool               `bson:"disable" json:"-"`
}

type UserCrypto struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Coin      string             `bson:"coin" json:"coin"`
	Vat       float64            `bson:"vat" json:"vat"`
	Wallet    string             `bson:"wallet" json:"-"`
	WhoDefine string             `bson:"who_define" json:"-"`
	Address   string             `bson:"address" json:"-"`
	Access    string             `bson:"access" json:"-"`
	Token     string             `bson:"token" json:"-"`
	Disable   bool               `bson:"disable" json:"-"`
}

type UserDebitCard struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Vat       float64            `bson:"vat" json:"vat"`
	Address   string             `bson:"address" json:"-"`
	Access    string             `bson:"access" json:"-"`
	Token     string             `bson:"token" json:"-"`
	WhoDefine string             `bson:"who_define" json:"-"`
	Disable   bool               `bson:"disable" json:"-"`
}
