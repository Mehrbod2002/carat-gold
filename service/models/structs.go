package models

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Avatar            string             `bson:"avatar" json:"avatar"`
	Name              string             `bson:"name" json:"name"`
	Email             string             `bson:"email" json:"email"`
	PhoneNumber       string             `bson:"phone" json:"phone"`
	Password          string             `bson:"password" json:"password"`
	PhoneVerified     bool               `bson:"phone_verified" json:"phone_verified"`
	OtpCode           *int               `bson:"otp_code" json:"otp_code"`
	Address           string             `bson:"address" json:"address"`
	RegisterCompleted bool               `bson:"register_completed" json:"register_completed"`
	ResetToken        string             `bson:"reset_token" json:"reset_token"`
	ResetTokenValid   time.Time          `bson:"reset_token_valid" json:"reset_token_valid"`
	ChangePhone       bool               `bson:"change_phone" json:"change_phone"`
	ExchangeMobile    string             `bson:"exchange_phone" json:"exchange_phone"`
	Freeze            bool               `bson:"freeze" json:"freeze"`
	File              string             `bson:"document" json:"document"`
	ChatList          []string           `bson:"chat_list" json:"chat_list"`
	Permissions       []Permission       `bson:"permissions" json:"permissions"`
	OtpValid          time.Time          `bson:"otp_valid" json:"otp_valid"`
	ReTryOtp          int                `bson:"retry_otp" json:"retry_otp"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	IsSupportOrAdmin  bool               `bson:"support_or_admin" json:"support_or_admin"`
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
	CreatedAt   time.Time `json:"created_at"`
	jwt.StandardClaims
}

type Message struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"-"`
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
	Currency   string `bson:"currency" json:"currency"`
	USDConvert int64  `bson:"usd_convert" json:"usd_convert"`
}

type Products struct {
	Currency   string `bson:"currency" json:"currency"`
	USDConvert int64  `bson:"usd_convert" json:"usd_convert"`
}

type Symbols struct {
	Symbols []string `bson:"symbols" json:"symbols"`
}

type PaymentMethods struct{}
type DeliveryMethods struct{}
type GeneralData struct{}
type HistoryOrders struct{}
type RealTimeOrders struct{}
type MetaData struct{}
type ChatHistories struct{}
