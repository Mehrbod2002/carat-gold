package models

import (
	"carat-gold/utils"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IsValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email)
}

func IsValidPhoneNumber(phoneNumber string) bool {
	phoneRegex := `^\+\d{1,4}\d{6,14}$`
	re := regexp.MustCompile(phoneRegex)

	return re.MatchString(phoneNumber)
}

func IsValidPassowrd(password string, c *gin.Context) bool {
	if len(password) > 36 || len(password) < 6 {
		utils.Method(c, "invalid password length")
		return false
	}
	if strings.Contains(password, ".") {
		utils.Method(c, "password not allowed to includes '.'")
		return false
	}
	return true
}

func (req *RequestSetDefineUser) Validate(c *gin.Context, Edit bool) bool {
	if len(*req.Name) > 20 || len(*req.Name) < 2 {
		utils.Method(c, "invalid name length")
		return false
	}
	if !IsValidPhoneNumber(*req.Phone) {
		utils.Method(c, "invalid phone")
		return false
	}
	if *req.IsSupport {
		if len(req.Permissions.Actions) == 0 {
			utils.Method(c, "as support , needs at least one permission")
			return false
		}
		if !ActionChecker(req.Permissions.Actions) {
			utils.Method(c, "invalid password permissions")
			return false
		}
		if !IsValidEmail(*req.Email) {
			utils.Method(c, "invalid email")
			return false
		}
		if !Edit {
			if valid := IsValidPassowrd(*req.Password, c); !valid {
				return false
			}
		}
	}
	if len(*req.Reason) > 200 {
		utils.Method(c, "invalid reason length")
		return false
	}
	if len(*req.Address) > 300 {
		utils.Method(c, "invalid address length")
		return false
	}
	return true
}

func ErrInSocket(ws *websocket.Conn, user *User, message string) error {
	err := ws.WriteJSON(Socket{
		ResponseTo: *user,
		Trigger:    "error",
		Validate:   false,
		Message:    message,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func ValidateSession(c *gin.Context) (*User, bool) {
	session := sessions.Default(c)
	token := session.Get("token")
	tokenString := c.GetHeader("Authorization")
	tokenAdmins := session.Get("token_admins")
	tokenSupports := session.Get("token_supports")
	cookie_token, err := c.Request.Cookie("token")
	if token == nil && tokenString == "" && err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"data":    "unauthorized",
			"message": "Unauthorized",
		})
		return nil, false
	}
	if tokenAdmins != "" {
		token = tokenAdmins
	}
	if token == "" && tokenSupports != "" {
		token = tokenSupports
	}
	if token == nil {
		token = tokenString
	}
	if token == "" {
		token = cookie_token.Value
	}
	jwtSecret := os.Getenv("SESSION_SECRET")
	if jwtSecret == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"data":    "unauthorized",
			"message": "JWT secret not configured",
		})
		return nil, false
	}

	parsedToken, err := jwt.Parse(token.(string), func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !parsedToken.Valid {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"data":    "unauthorized",
			"message": "Invalid token",
		})
		return nil, false
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		if claims["email"] == os.Getenv("ADMIN_USERNAME") {
			createdStr := claims["created_at"].(string)
			email := claims["email"].(string)
			createdAt, err := time.Parse(time.RFC3339, createdStr)
			if err != nil {
				log.Println(err)
				return nil, false
			}
			userID, ok := claims["_id"].(string)
			if !ok {
				return nil, false
			}
			if userID, err := primitive.ObjectIDFromHex(userID); err == nil {
				user := &User{
					ID:        userID,
					Email:     email,
					CreatedAt: createdAt,
				}
				return user, true
			}
			return nil, false
		}
		if claims["_id"] != nil {
			userID, ok := claims["_id"].(string)
			if !ok {
				return nil, false
			}
			if userID, err := primitive.ObjectIDFromHex(userID); err == nil {
				createdStr := claims["created_at"].(string)
				email := claims["email"].(string)
				createdAt, err := time.Parse(time.RFC3339, createdStr)
				if err != nil {
					log.Println(err)
					return nil, false
				}
				if err != nil {
					log.Println(err)
					return nil, false
				}
				user := &User{
					ID:        userID,
					Email:     email,
					CreatedAt: createdAt,
				}
				return user, false
			}
		}
		return nil, false
	}
	return nil, false
}

func ReceiveSession(c *gin.Context) *User {
	session := sessions.Default(c)
	token := session.Get("token")
	cookie_token, err := c.Request.Cookie("token")
	tokenString := c.GetHeader("Authorization")
	if token == nil && tokenString == "" && err != nil {
		log.Println(err)
		return nil
	}
	if token == nil {
		token = tokenString
	}
	if token == "" {
		token = cookie_token.Value
	}
	jwtSecret := os.Getenv("SESSION_SECRET")
	if jwtSecret == "" {
		return nil
	}

	parsedToken, err := jwt.Parse(token.(string), func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !parsedToken.Valid {
		log.Println(err)
		return nil
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && claims["_id"] != nil {
		userID, ok := claims["_id"].(string)
		if !ok {
			return nil
		}
		if userID, err := primitive.ObjectIDFromHex(userID); err == nil {
			createdStr := claims["created_at"].(string)
			email := claims["email"].(string)
			createdAt, errs := time.Parse(time.RFC3339, createdStr)
			if errs != nil {
				log.Println(errs)
				return nil
			}
			user := &User{
				ID:          userID,
				PhoneNumber: email,
				CreatedAt:   createdAt,
			}

			return user
		}
		return nil
	}
	return nil
}

func AllowedAction(c *gin.Context, action Action) bool {
	user, isAdmin := ValidateSession(c)
	if isAdmin {
		return true
	}
	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		log.Println(DBerr)
		return false
	}
	var currentUser User
	tableName := "users"
	err := db.Collection(tableName).
		FindOne(context.Background(), bson.M{
			"_id": user.ID,
		}).Decode(&currentUser)
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return false
	}
	for _, act := range currentUser.Permissions.Actions {
		if act == action {
			return true
		}
	}
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"success": false,
		"message": "you don't have this permission",
		"data":    "invalid_permission",
	})
	return false
}

func (user *User) GenerateToken() (string, error) {
	claims := &Claims{
		ID:          user.ID.Hex(),
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		PhoneNumber: user.PhoneNumber,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("SESSION_SECRET")))
	if err != nil {
		log.Println(err)
		return "", err
	}

	return signedToken, nil
}

func ActionChecker(actions []Action) bool {
	for _, reqAct := range actions {
		found := false
		for _, act := range AllActions {
			if reqAct == act {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func (loginData *LoginDataStep1) Validate(c *gin.Context) bool {
	if IsValidPhoneNumber(loginData.Phone) {
		utils.Method(c, "invalid phone number")
		return false
	}
	return true
}

func (loginData *LoginDataStep2) Validate(c *gin.Context) bool {
	if IsValidPhoneNumber(loginData.Phone) {
		utils.Method(c, "invalid phone number")
		return false
	}
	return true
}

func (sendOTPData *SendOTP) Validate(c *gin.Context) bool {
	if IsValidPhoneNumber(sendOTPData.PhoneNumber) {
		utils.Method(c, "invalid phone number")
		return false
	}
	return true
}

func (registerRequest *RegisterRequest) Validate(c *gin.Context) bool {
	if len(registerRequest.Name) < 3 {
		utils.Method(c, "name is short")
		return false
	}
	if !IsValidPhoneNumber(registerRequest.PhoneNumber) {
		utils.Method(c, "invalid phone number")
		return false
	}
	return true
}

func (body *Documents) Validate(c *gin.Context) bool {
	decodedFile, err := base64.StdEncoding.DecodeString(body.FrontShot)
	if err != nil {
		utils.Method(c, "invalid front file format")
		return false
	}
	fileSizeMB := float64(len(decodedFile)) / (1024 * 1024)
	if fileSizeMB > 30 {
		utils.Method(c, "front shot size exceeds 30 MB")
		return false
	}

	decodedFile, err = base64.StdEncoding.DecodeString(body.BackShot)
	if err != nil {
		utils.Method(c, "invalid back file format")
		return false
	}
	fileSizeMB = float64(len(decodedFile)) / (1024 * 1024)
	if fileSizeMB > 30 {
		utils.Method(c, "back shot size exceeds 30 MB")
		return false
	}

	return true
}

func (registerRequest *RequestSetSymbol) Validate(c *gin.Context) bool {
	if len(*registerRequest.SymbolName) < 3 {
		utils.Method(c, "symbol name is short")
		return false
	}
	return true
}

func (order *RequestSetTrade) Validate(c *gin.Context) bool {
	if len(order.SymbolName) == 0 {
		utils.Method(c, "symbol is missed")
		return false
	}
	if order.Volumn == 0 {
		utils.Method(c, "volumn is missed")
		return false
	}
	if order.Comment == nil {
		defaultComment := ""
		order.Comment = &defaultComment
	}
	if order.TakeProfit == nil {
		defaultTakeProfit := 0.0
		order.TakeProfit = &defaultTakeProfit
	}
	if order.StopLoss == nil {
		defaultStopLoss := 0.0
		order.StopLoss = &defaultStopLoss
	}
	if order.Deviation == nil {
		defaultDeviation := 0.0
		order.Deviation = &defaultDeviation
	}
	if order.Stoplimit == nil {
		defaultStoplimit := 0.0
		order.Stoplimit = &defaultStoplimit
	}
	return true
}

func (order *RequestSetCancelTrade) Validate(c *gin.Context) bool {
	return true
}

func (product *RequestSetProduct) Validate(c *gin.Context) bool {
	for _, i := range product.Images {
		decodedFile, err := base64.StdEncoding.DecodeString(i)
		if err != nil {
			utils.Method(c, "invalid file format")
			return false
		}
		fileSizeMB := float64(len(decodedFile)) / (1024 * 1024)
		if fileSizeMB > 30 {
			utils.Method(c, "front shot size exceeds 30 MB")
			return false
		}
	}

	if len(product.Name) == 0 {
		utils.Method(c, "invalid name")
		return false
	}
	if len(product.Description) == 0 {
		utils.Method(c, "invalid description")
		return false
	}
	if product.Width <= 0 {
		utils.Method(c, "invalid width")
		return false
	}
	if product.Length <= 0 {
		utils.Method(c, "invalid length")
		return false
	}
	if product.WeightOZ <= 0 {
		utils.Method(c, "invalid weight oz")
		return false
	}
	if product.WeightGramm <= 0 {
		utils.Method(c, "invalid weight gramm")
		return false
	}
	if product.Purity <= 0 {
		utils.Method(c, "invalid name")
		return false
	}
	if product.Percentage <= 0 {
		utils.Method(c, "invalid percentage")
		return false
	}
	return true
}

func CreateOrder(order *RequestSetTrade) (string, string) {
	slippage := "0"
	stopLoss := "0"
	takeProfit := "0"
	comment := "Default comment"
	requestID := fmt.Sprintf("%d", utils.GenerateRandomCode())[1:]
	volumn := fmt.Sprintf("%f", order.Volumn)
	operation := fmt.Sprintf("%d", order.Operation)
	if order.Deviation != nil {
		slippage = fmt.Sprintf("%f", *order.Deviation)
	}
	if order.TakeProfit != nil {
		takeProfit = fmt.Sprintf("%f", *order.TakeProfit)
	}
	if order.StopLoss != nil {
		stopLoss = fmt.Sprintf("%f", *order.StopLoss)
	}
	if order.Comment != nil {
		comment = *order.Comment
	}
	expirationTime := time.Now().Add(24 * time.Hour)

	expirationTimeString := expirationTime.Format("2006.01.02 15:04:00")

	orderStr := requestID + "|OPEN_TRADE|" + order.SymbolName + "|" +
		operation + "|" + volumn + "|" + slippage + "|" +
		stopLoss + "|" + takeProfit + "|" +
		comment + "|" + requestID + "|" + expirationTimeString

	return requestID, orderStr
}

func CancelOrder(order *RequestSetCancelTrade) (string, string) {
	requestID := fmt.Sprintf("%d", utils.GenerateRandomCode())[1:]
	ticket := fmt.Sprintf("%d", order.Ticket)

	orderStr := requestID + "|CLOSE_TRADE|" + "|" + ticket
	return requestID, orderStr
}

func GetCurrentOrder() (string, string) {
	requestID := fmt.Sprintf("%d", utils.GenerateRandomCode())[1:]

	orderStr := requestID + "|CURRENT_ORDERS|"
	return requestID, orderStr
}

func GetHistoryOrder() (string, string) {
	requestID := fmt.Sprintf("%d", utils.GenerateRandomCode())[1:]

	orderStr := requestID + "|HISTORY_ORDERS|"
	return requestID, orderStr
}

func (delivery *RequestSetDeliveryMethod) Validate(c *gin.Context) bool {
	if len(delivery.Title) == 0 {
		utils.Method(c, "invalid title")
		return false
	}
	if len(delivery.Description) == 0 {
		utils.Method(c, "invalid description")
		return false
	}
	if delivery.TimeProvided {
		if delivery.EstimatedTime == "" {
			utils.Method(c, "invalid time")
			return false
		}
	}
	if delivery.Fee == 0 {
		utils.Method(c, "invalid Fee")
		return false
	}
	return true
}

func (delivery *RequestSetFANDQ) Validate(c *gin.Context) bool {
	if len(delivery.Question) == 0 {
		utils.Method(c, "invalid question")
		return false
	}
	if len(delivery.Answer) == 0 {
		utils.Method(c, "invalid answer")
		return false
	}
	return true
}

func (req *RequestSetPayment) Validate(c *gin.Context) bool {
	if len(req.Access) == 0 {
		utils.Method(c, "access is missed")
		return false
	}
	if len(req.Address) == 0 {
		utils.Method(c, "address is missed")
		return false
	}
	if len(req.Token) == 0 {
		utils.Method(c, "token is missed")
		return false
	}
	if req.Vat == 0 {
		utils.Method(c, "vat is missed")
		return false
	}
	return true
}

func (req *RequestSetCallCenter) Validate(c *gin.Context) bool {
	if len(*req.Email) != 0 {
		if !IsValidEmail(*req.Email) {
			utils.Method(c, "invalid email")
			return false
		}
	}
	if len(*req.PhoneComapny) == 0 {
		if !IsValidPhoneNumber(*req.PhoneComapny) {
			utils.Method(c, "invalid phone")
			return false
		}
	}
	return true
}
