package models

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Avatar           string             `bson:"avatar" json:"avatar"`
	Name             string             `bson:"name" json:"name"`
	Email            string             `bson:"email" json:"email"`
	PhoneNumber      string             `bson:"phone" json:"phone"`
	Password         string             `bson:"password" json:"password"`
	PhoneVerified    bool               `bson:"phone_verified" json:"phone_verified"`
	OtpCode          *int               `bson:"otp_code" json:"otp_code"`
	Address          string             `bson:"address" json:"address"`
	ResetToken       string             `bson:"reset_token" json:"reset_token"`
	ResetTokenValid  time.Time          `bson:"reset_token_valid" json:"reset_token_valid"`
	ChangePhone      bool               `bson:"change_phone" json:"change_phone"`
	ExchangeMobile   string             `bson:"exchange_phone" json:"exchange_phone"`
	Freeze           bool               `bson:"freeze" json:"freeze"`
	Documents        Documents          `bson:"documents" json:"documents"`
	ChatList         []string           `bson:"chat_list" json:"chat_list"`
	Permissions      Permission         `bson:"permissions" json:"permissions"`
	OtpValid         time.Time          `bson:"otp_valid" json:"otp_valid"`
	ReTryOtp         int                `bson:"retry_otp" json:"retry_otp"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UserVerified     bool               `bson:"user_verified" json:"user_verified"`
	StatusString     UserStatus         `bson:"user_status" json:"user_status"`
	Reason           string             `bson:"reason" json:"reason"`
	IsSupportOrAdmin bool               `bson:"support_or_admin" json:"support_or_admin"`
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
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Currency string             `bson:"currency,omitempty" json:"currency"`
}

type Products struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name      string             `bson:"name,omitempty" json:"name"`
	WhoDefine string             `bson:"who_define" json:"who_define"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type Symbol struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	SymbolName string             `bson:"name" json:"name"`
	SymbolType SymbolType         `bson:"type" json:"type"`
	SymbolSide SymbolSide         `bson:"side" json:"side"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

type PaymentMethods struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	PaymentName string             `bson:"name,omitempty" json:"name"`
	Provider    map[string]string  `bson:"provider,omitempty" json:"provider"`
	Description string             `bson:"description,omitempty" json:"description"`
	WhoDefine   string             `bson:"who_define" json:"who_define"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

type DeliveryMethods struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	DeliveryName string             `bson:"name,omitempty" json:"name"`
	Description  string             `bson:"description,omitempty" json:"description"`
	WhoDefine    string             `bson:"who_define" json:"who_define"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

type GeneralData struct{}
type HistoryOrders struct{}
type RealTimeOrders struct{}
type MetaData struct{}
type ChatHistories struct{}

type Transctions struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	OrderID       string             `json:"order_id" bson:"order_id,omitempty"`
	UserID        primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	Status        UserStatus         `bson:"status" json:"status"`
	PaymentMethod string             `bson:"payment_method" json:"payment_method"`
	Symbol        string             `bson:"symbol" json:"symbol"`
	ExternalData  map[string]string  `bson:"external_data" json:"external_data"`
}

type Purchaed struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	ProductID primitive.ObjectID `bson:"product_id,omitempty" json:"product_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type RequestSetDefineUser struct {
	UserID      *string     `json:"user_id"`
	Name        *string     `json:"name"`
	Email       *string     `json:"email"`
	Password    *string     `json:"password"`
	Phone       *string     `json:"phone"`
	PhoneVerify *bool       `json:"phone_verified"`
	IsSupport   *bool       `json:"support"`
	Freeze      *bool       `json:"freeze"`
	Status      *UserStatus `json:"status"`
	Reason      *string     `json:"reason"`
	Permissions *Permission `json:"permissions"`
	Address     *string     `json:"address"`
}

type LoginDataStep1 struct {
	Phone string `json:"phone"`
}

type LoginDataStep2 struct {
	Phone string `json:"phone"`
	Otp   *int   `json:"otp"`
}

type SendOTP struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone" binding:"required"`
}

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone" binding:"required"`
	OtpCode     *int   `json:"otp_code" binding:"required"`
}

type Documents struct {
	FrontShot string `json:"front_shot" bson:"front_shot"`
	BackShot  string `json:"back_shot" bson:"back_shot"`
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
	SymbolName *string     `json:"symbol_name"`
	SymbolType *SymbolType `json:"symbol_type"`
	SymbolSide *SymbolSide `json:"symbol_side"`
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
	Slippage   *float64 `json:"slippage"`
}

type RequestSetCancelTrade struct {
	SymbolName string `json:"name"`
	Ticket     int    `json:"ticket"`
}
