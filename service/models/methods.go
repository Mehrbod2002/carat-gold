package models

import (
	"bytes"
	"carat-gold/utils"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/skip2/go-qrcode"
	"github.com/twilio/twilio-go"
	twilioclient "github.com/twilio/twilio-go/client"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GenerateTransactionID(name, last string, products int) string {
	if len(name) == 0 {
		name = "NN"
	}
	if len(last) == 0 {
		last = "NN"
	}

	formatted := fmt.Sprintf("%04d", products)
	rand := fmt.Sprintf("%d", utils.GenerateRandomCode())

	result := string(name[0]) + string(last[0]) + "-GB8-" + formatted + "-" + rand[0:3]
	return result
}

func GenerateOrderID(c *gin.Context) (*string, error) {
	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		return nil, DBerr
	}

	var lastOrderID LastOrderID
	err := db.Collection("last_order_id").FindOne(context.Background(), bson.M{}).Decode(&lastOrderID)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	var rand string
	if err == mongo.ErrNoDocuments {
		rand = "603155"
	} else {
		numberDig, err := strconv.Atoi(lastOrderID.OrderID)
		if err != nil {
			return nil, err
		}

		rand = fmt.Sprintf("%d", numberDig+1)
	}

	return &rand, nil
}

func CreateCryptoInvoice(c *gin.Context, price float64, orderID string) (*Invoice, bool, string) {
	url := "https://api.nowpayments.io/v1/invoice"

	payloadData := struct {
		PriceAmount      float64 `json:"price_amount"`
		PriceCurrency    string  `json:"price_currency"`
		PayCurrency      string  `json:"pay_currency"`
		IPNCallbackURL   string  `json:"ipn_callback_url"`
		OrderID          string  `json:"order_id"`
		OrderDescription string  `json:"order_description"`
		SuccessURL       string  `json:"success_url"`
		CancelURL        string  `json:"cancel_url"`
	}{
		PriceAmount:      price,
		PriceCurrency:    "usd",
		PayCurrency:      "usdtbsc",
		IPNCallbackURL:   os.Getenv("BASE_HOST") + os.Getenv("CALLBACK"),
		OrderID:          orderID,
		OrderDescription: "Fasih Products",
		SuccessURL:       os.Getenv("SUCCESS_URL"),
		CancelURL:        os.Getenv("CANCEL_URL"),
	}

	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		return nil, false, "internal error"
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, false, "internal error"
	}

	req.Header.Set("x-api-key", os.Getenv("CRYPTO_SECRET"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, false, "internal error"
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			return nil, false, "internal error"
		}

		var errMessage map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &errMessage)
		if err != nil {
			return nil, false, "internal error"
		}
		errMsg, ok := errMessage["message"].(string)
		if !ok {
			return nil, false, "internal error"
		}
		return nil, false, errMsg
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, false, "internal error"
	}

	var paymentResponse Invoice
	err = json.Unmarshal(buf.Bytes(), &paymentResponse)

	if err != nil {
		return nil, false, "internal error"
	}

	return &paymentResponse, true, ""
}

// func CreateCrypto(c *gin.Context, price float64, orderID string, secret string) (*PaymentResponse, error) {
// 	url := "https://api.nowpayments.io/v1/payment"

// 	payloadData := struct {
// 		PriceAmount      float64 `json:"price_amount"`
// 		PriceCurrency    string  `json:"price_currency"`
// 		PayCurrency      string  `json:"pay_currency"`
// 		IPNCallbackURL   string  `json:"ipn_callback_url"`
// 		OrderID          string  `json:"order_id"`
// 		OrderDescription string  `json:"order_description"`
// 	}{
// 		PriceAmount:      price,
// 		PriceCurrency:    "usd",
// 		PayCurrency:      "usdt",
// 		IPNCallbackURL:   os.Getenv("BASE_HOST") + "/" + os.Getenv("CALLBACK"),
// 		OrderID:          orderID,
// 		OrderDescription: "Fasih Products",
// 	}

// 	payloadBytes, err := json.Marshal(payloadData)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("x-api-key", secret)
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != 201 {
// 		var serialized map[string]interface{}
// 		response, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to create url")
// 		}

// 		err = json.Unmarshal(response, &serialized)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to create url: %v", err)
// 		}

// 		if message, ok := serialized["message"].(string); ok {
// 			return nil, fmt.Errorf(message)
// 		}

// 		return nil, fmt.Errorf("failed to create url")
// 	}

// 	buf := new(bytes.Buffer)
// 	_, err = buf.ReadFrom(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var paymentResponse PaymentResponse
// 	err = json.Unmarshal(buf.Bytes(), &paymentResponse)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &paymentResponse, nil
// }

func CreateQr(payment string) (*string, error) {
	qrCode, err := qrcode.Encode(string(payment), qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	qrBase64 := base64.StdEncoding.EncodeToString(qrCode)
	return &qrBase64, nil
}

func Notification(userID primitive.ObjectID, title string, notification string) error {
	app := utils.GetApp()

	db, err := utils.GetDBServer()
	if err != nil {
		return err
	}

	var user User
	exist := db.Collection("users").FindOne(context.Background(), bson.M{
		"_id": userID,
	}).Decode(&user)
	if exist != nil {
		log.Println(exist)
		return exist
	}

	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		return err
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  notification,
		},
		Token: user.FcmToken,
	}

	_, err = client.Send(ctx, message)
	if err != nil {
		return err
	}

	return nil
}

func IsValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email)
}

func IsValidPhoneNumber(phoneNumber string) bool {
	phoneRegex := `^\+\d{1,4}\d{6,14}$`
	re := regexp.MustCompile(phoneRegex)

	return re.MatchString("+" + phoneNumber)
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

func (req *RequestEdit) Validate(c *gin.Context) bool {
	if req.FirstName != nil {
		if len(*req.FirstName) != 0 {
			if len(*req.FirstName) > 20 || len(*req.FirstName) < 2 {
				utils.Method(c, "invalid first name")
				return false
			}
		}
	}
	if req.LastName != nil && *req.LastName != "" {
		if len(*req.LastName) != 0 {
			if len(*req.LastName) > 20 || len(*req.LastName) < 2 {
				utils.Method(c, "invalid last name ")
				return false
			}
		}
	}

	if req.Email != nil && *req.Email != "" {
		if !IsValidEmail(*req.Email) {
			utils.Method(c, "invalid email address")
			return false
		}
	}
	if req.Phone != nil && *req.Phone != "" {
		if !IsValidPhoneNumber(*req.Phone) {
			utils.Method(c, "invalid phone")
			return false
		}
	}
	if req.Address != nil && len(*req.Address) != 0 {
		for _, address := range *req.Address {
			if len(address.City) > 20 || len(address.City) == 0 {
				utils.Method(c, "invalid city")
				return false
			}
			if len(address.Country) > 20 || len(address.Country) == 0 {
				utils.Method(c, "invalid country")
				return false
			}
			if len(address.Region) > 20 || len(address.Region) == 0 {
				utils.Method(c, "invalid region")
				return false
			}
			if len(address.Address) > 300 || len(address.Address) == 0 {
				utils.Method(c, "invalid address")
				return false
			}
			if len(address.Label) > 20 || len(address.Label) < 3 {
				utils.Method(c, "invalid label")
				return false
			}
		}
	}
	return true
}

func (req *RequestSetDefineUser) Validate(c *gin.Context, Edit bool) bool {
	if req.FirstName != nil {
		if len(*req.FirstName) > 20 || len(*req.FirstName) < 2 {
			utils.Method(c, "invalid first name length")
			return false
		}
	}
	if req.LastName != nil {
		if len(*req.LastName) > 20 || len(*req.LastName) < 2 {
			utils.Method(c, "invalid last name length")
			return false
		}
	}
	if req.BalanceUSD != nil {
		balance, _ := strconv.ParseFloat(*req.BalanceUSD, 64)
		if balance < 0 {
			utils.Method(c, "invalid balance")
			return false
		}
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
			utils.Method(c, "invalid email address")
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

	// if req.Address != nil && len(*req.Address) != 0 {
	// 	for _, address := range *req.Address {
	// 		if len(address.City) > 20 || len(address.City) == 0 {
	// 			utils.Method(c, "invalid city length")
	// 			return false
	// 		}
	// 		if len(address.Country) > 20 || len(address.Country) == 0 {
	// 			utils.Method(c, "invalid country length")
	// 			return false
	// 		}
	// 		if len(address.Region) > 20 || len(address.Region) == 0 {
	// 			utils.Method(c, "invalid region length")
	// 			return false
	// 		}
	// 		if len(address.Address) > 300 || len(address.Address) == 0 {
	// 			utils.Method(c, "invalid address length")
	// 			return false
	// 		}
	// 		if len(address.Label) > 50 || len(address.Label) == 0 {
	// 			utils.Method(c, "invalid address label")
	// 			return false
	// 		}
	// 	}
	// }
	return true
}

func UserExists(c *gin.Context, id primitive.ObjectID) bool {
	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		log.Println(DBerr)
		return false
	}
	var currentUser User
	err := db.Collection("users").
		FindOne(context.Background(), bson.M{
			"_id": id,
		}).Decode(&currentUser)
	return err == nil
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
		utils.Unauthorized(c)
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
		utils.Unauthorized(c)
		return nil, false
	}

	parsedToken, err := jwt.Parse(token.(string), func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !parsedToken.Valid {
		log.Println(err)
		utils.Unauthorized(c)
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

				exists := UserExists(c, userID)
				if !exists {
					session.Delete("token")
					err = session.Save()
					if err != nil {
						log.Println(err)
						return nil, false
					}
					return nil, false
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
				user := &User{
					ID:        userID,
					Email:     email,
					CreatedAt: createdAt,
				}

				exists := UserExists(c, userID)
				if !exists {
					session.Delete("token")
					err = session.Save()
					if err != nil {
						return nil, false
					}
					return nil, false
				}
				return user, true
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
	fmt.Println("ff:", user.PhoneNumber)
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
	if !IsValidPhoneNumber(loginData.Phone) {
		utils.Method(c, "invalid phone number")
		return false
	}
	return true
}

func (loginData *LoginDataStep2) Validate(c *gin.Context) bool {
	if !IsValidPhoneNumber(loginData.Phone) {
		utils.Method(c, "invalid phone number")
		return false
	}
	return true
}

func (sendOTPData *SendOTP) Validate(c *gin.Context) bool {
	if !IsValidPhoneNumber(sendOTPData.PhoneNumber) {
		utils.Method(c, "invalid phone number")
		return false
	}
	return true
}

func (registerRequest *RegisterRequest) Validate(c *gin.Context) bool {
	if !IsValidPhoneNumber(registerRequest.PhoneNumber) {
		utils.Method(c, "invalid phone number")
		return false
	}
	return true
}

func (body *Documents) Validate(c *gin.Context) bool {
	decodedFileFront, err := base64.StdEncoding.DecodeString(body.Documents.Front.Shot)
	if err != nil {
		utils.Method(c, "invalid front file format")
		return false
	}
	fileSizeMBFront := float64(len(decodedFileFront)) / (1024 * 1024)
	if fileSizeMBFront > 10 {
		utils.Method(c, "front shot size exceeds 10 MB")
		return false
	}

	decodedFile, err := base64.StdEncoding.DecodeString(body.Documents.Back.Shot)
	if err != nil {
		utils.Method(c, "invalid front file format")
		return false
	}
	fileSizeMB := float64(len(decodedFile)) / (1024 * 1024)
	if fileSizeMB > 8 {
		utils.Method(c, "front shot size exceeds 8 MB")
		return false
	}

	return true
}

func (requestSymbol *RequestSetSymbol) Validate(c *gin.Context) bool {
	decodedFile, err := base64.StdEncoding.DecodeString(strings.Split(requestSymbol.Image, ",")[1])
	if err != nil {
		utils.Method(c, "invalid file format")
		return false
	}
	fileSizeMB := float64(len(decodedFile)) / (1024 * 1024)
	if fileSizeMB > 30 {
		utils.Method(c, "front shot size exceeds 30 MB")
		return false
	}
	if len(requestSymbol.Name) < 3 {
		utils.Method(c, "symbol name is short")
		return false
	}
	if len(requestSymbol.SymbolMetaName) == 0 {
		utils.Method(c, "symbol name is invalid")
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
	for _, i := range *product.Images {
		decodedFile, err := base64.StdEncoding.DecodeString(strings.Split(i, ",")[1])
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

	if len(*product.Name) == 0 {
		utils.Method(c, "invalid name")
		return false
	}
	if len(*product.Description) == 0 {
		utils.Method(c, "invalid description")
		return false
	}
	if len(*product.SubTitle) == 0 || len(*product.SubTitle) > 100 {
		utils.Method(c, "invalid sub name")
		return false
	}
	if len(*product.Faq) == 0 || len(*product.Faq) > 100 {
		utils.Method(c, "invalid faq name")
		return false
	}
	if len(*product.Answer) == 0 || len(*product.Answer) == 200 {
		utils.Method(c, "invalid answer")
		return false
	}
	if *product.Width <= 0 {
		utils.Method(c, "invalid width")
		return false
	}
	if *product.Length <= 0 {
		utils.Method(c, "invalid length")
		return false
	}
	if *product.WeightOZ <= 0 {
		utils.Method(c, "invalid weight oz")
		return false
	}
	if *product.WeightGramm <= 0 {
		utils.Method(c, "invalid weight gramm")
		return false
	}
	if product.Purity != nil {
		if *product.Purity != 999 && *product.Purity != 995 {
			utils.Method(c, "invalid purity,not 995 either 999")
			return false
		}
	}

	if product.PurityStr != nil {
		if len(*product.PurityStr) == 0 {
			utils.Method(c, "invalid purity string")
			return false
		}
	}
	if product.Percentage != nil {
		if *product.Percentage < 0 {
			utils.Method(c, "invalid percentage")
			return false
		}
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
	if delivery.Fee < 0 {
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
			utils.Method(c, "invalid email address")
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

func (meta *RequestMetaTraderAccounts) Validate(c *gin.Context) bool {
	if len(meta.Server) == 0 {
		utils.Method(c, "invalid server")
		return false
	}
	if len(meta.Login) == 0 {
		utils.Method(c, "invalid login")
		return false
	}
	if len(meta.Passowrd) == 0 {
		utils.Method(c, "invalid passowrd")
		return false
	}
	return true
}

func (p *PaymentCallBack) UnmarshalJSON(data []byte) error {
	var tmp struct {
		UpdatedAt          int64            `json:"updated_at"`
		PaymentID          int64            `json:"payment_id"`
		ParentPaymentID    int64            `json:"parent_payment_id"`
		InvoiceID          interface{}      `json:"invoice_id"`
		PaymentStatus      NowPaymentStatus `json:"payment_status"`
		PayAddress         string           `json:"pay_address"`
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

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	p.UpdatedAt = time.Unix(tmp.UpdatedAt/1000, 0)
	return nil
}

func HandleIPN(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return
	}

	// signature := c.GetHeader("x-nowpayments-sig")
	// verified := VerifyIPN(c.Request, signature)
	// if !verified {
	// 	utils.InternalError(c)
	// 	return
	// }

	var payment PaymentCallBack
	if err := payment.UnmarshalJSON(body); err != nil {
		utils.InternalError(c)
		return
	}

	if payment.PaymentStatus == PaymentFinished {
		valid := Pay(payment.OrderID)
		if !valid {
			user := GetUserByPayment(payment.OrderID)
			if user != nil {
				Notification(user.ID, "Invoice", "Payments cancelled #"+payment.OrderID.(string))
			}
		}
	} else if payment.PaymentStatus == PaymentConfirming {
		transaction := GetTransaction(payment.OrderID)
		if transaction != nil {
			if !transaction.IsDebit {
				mID, valid := utils.AutoOrder(transaction.TotalPrice, false)
				if valid {
					StoreMetatraderID(transaction.OrderID, fmt.Sprintf("%d", mID))
				}
			}
		}
	} else if payment.PaymentStatus == PaymentFailed ||
		payment.PaymentStatus == PaymentRefunded ||
		payment.PaymentStatus == PaymentPartiallyPaid ||
		payment.PaymentStatus == PaymentExpired {
		valid, isDebit := Cancel(payment.OrderID)
		if !valid {
			user := GetUserByPayment(payment.OrderID)
			if user != nil {
				Notification(user.ID, "Invoice", "Your invoice is cancelled . Please contact supports for invoice #"+payment.OrderID.(string))
			}
		} else {
			if !isDebit {
				transaction := GetTransaction(payment.OrderID)
				if transaction != nil {
					payload := map[string]interface{}{
						"ticket_id": transaction.MetatraderID,
					}
					utils.PostRequest(payload, "cancel_order")
				}
			}
		}
	}
}

func StoreMetatraderID(orderIDInterface interface{}, metatraderID string) {
	orderID, valid := orderIDInterface.(string)
	if !valid {
		return
	}

	db, err := utils.GetDBServer()
	if err != nil {
		log.Println(err)
		return
	}

	if _, err = db.Collection("transactions").UpdateOne(context.Background(), bson.M{
		"order_id": orderID,
	}, bson.M{
		"metatrader_id": metatraderID,
	}); err != nil {
		return
	}
}

func GetTransaction(orderIDInterface interface{}) *Transaction {
	orderID, valid := orderIDInterface.(string)
	if !valid {
		return nil
	}

	db, err := utils.GetDBServer()
	if err != nil {
		log.Println(err)
		return nil
	}

	var transaction Transaction
	if err = db.Collection("transactions").FindOne(context.Background(), bson.M{
		"order_id": orderID,
	}).Decode(&transaction); err != nil && err != mongo.ErrNoDocuments {
		return nil
	}

	return &transaction
}

func VerifyIPN(req *http.Request, receivedHMAC string) bool {
	requestJSON, err := io.ReadAll(req.Body)
	if err != nil {
		return false
	}

	db, err := utils.GetDB(&gin.Context{})
	if err != nil {
		return false
	}

	var cryptoDetail Crypto
	err = db.Collection("crypto").FindOne(context.Background(), bson.M{}).Decode(&cryptoDetail)
	if err != nil {
		return false
	}

	var requestData map[string]interface{}
	err = json.Unmarshal(requestJSON, &requestData)
	if err != nil {
		return false
	}

	keys := make([]string, 0, len(requestData))
	for k := range requestData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sortedRequestJSON string
	for _, k := range keys {
		sortedRequestJSON += fmt.Sprintf(`"%s":%v,`, k, requestData[k])
	}
	sortedRequestJSON = "{" + sortedRequestJSON[:len(sortedRequestJSON)-1] + "}"

	hmacHash := hmac.New(sha512.New, []byte(cryptoDetail.Access))
	hmacHash.Write([]byte(sortedRequestJSON))
	signatureCalculated := hex.EncodeToString(hmacHash.Sum(nil))

	return signatureCalculated == receivedHMAC
}

func SortedParamsToString(params map[string]interface{}) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var sortedString string
	for _, key := range keys {
		sortedString += fmt.Sprintf(`"%s":"%v",`, key, params[key])
	}
	return "{" + sortedString[:len(sortedString)-1] + "}"
}

func Sendotp(mobileNumber string) (bool, string) {
	mobileNumber = fmt.Sprintf("+%s", mobileNumber)
	accountSID := os.Getenv("SID")
	authToken := os.Getenv("SMS_TOKEN")
	verificationSid := os.Getenv("VERIFY")
	channel := "sms"

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	_, err := client.VerifyV2.CreateVerification(verificationSid, &openapi.CreateVerificationParams{
		To:      &mobileNumber,
		Channel: &channel,
	})

	if err != nil {
		twilioError := err.(*twilioclient.TwilioRestError)
		return false, twilioError.Error()
	}

	return true, ""
}

func GetUserByPayment(orderIDInterface interface{}) *User {
	orderID, valid := orderIDInterface.(string)
	if !valid {
		return nil
	}

	db, err := utils.GetDBServer()
	if err != nil {
		log.Println(err)
		return nil
	}

	var transaction Transaction
	if err = db.Collection("transactions").FindOne(context.Background(), bson.M{
		"order_id": orderID,
	}).Decode(&transaction); err != nil && err != mongo.ErrNoDocuments {
		return nil
	}

	var user User
	if err = db.Collection("transactions").FindOne(context.Background(), bson.M{
		"_id": transaction.UserID,
	}).Decode(&transaction); err != nil && err != mongo.ErrNoDocuments {
		return nil
	}

	return &user
}

func Pay(orderIDInterface interface{}) bool {
	orderID, valid := orderIDInterface.(string)
	if !valid {
		return false
	}

	db, err := utils.GetDBServer()
	if err != nil {
		log.Println(err)
		return false
	}

	var transaction Transaction
	if err = db.Collection("transactions").FindOne(context.Background(), bson.M{
		"order_id": orderID,
	}).Decode(&transaction); err != nil && err != mongo.ErrNoDocuments {
		return false
	}

	update := bson.M{
		"$set": bson.M{
			"$push": bson.M{
				"wallet.purchased": Purchased{
					PaymentStatus:  ApprovedStatus,
					PaymentMethd:   transaction.PaymentMethod,
					CreatePayment:  time.Now(),
					CreatedAt:      transaction.CreatedAt,
					StatusDelivery: transaction.StatusDelivery,
					Product:        transaction.ProductIDs,
					OrderID:        transaction.OrderID,
				},
			},
		},
	}

	if transaction.IsDebit {
		if _, err := db.Collection("users").UpdateOne(context.Background(), bson.M{
			"_id": transaction.UserID,
		}, bson.M{"$set": bson.M{
			"$inc": bson.M{
				"wallet.balance": transaction.TotalPrice - transaction.Vat,
			},
		}}); err != nil {
			return false
		}
	}

	if _, err := db.Collection("transactions").UpdateOne(context.Background(), bson.M{
		"order_id": orderID,
	}, bson.M{
		"$set": bson.M{
			"payment_status":     ApprovedStatus,
			"payment_completion": true,
		},
	}); err != nil {
		return false
	}

	if _, err := db.Collection("users").UpdateOne(context.Background(), bson.M{
		"_id": transaction.UserID,
	}, update); err != nil {
		return false
	}

	Notification(transaction.UserID, "Invoice", "Payments set with #"+transaction.OrderID)
	return true
}

func Cancel(orderIDInterface interface{}) (bool, bool) {
	orderID, valid := orderIDInterface.(string)
	if !valid {
		return false, false
	}

	db, err := utils.GetDBServer()
	if err != nil {
		log.Println(err)
		return false, false
	}

	var transaction Transaction
	if err = db.Collection("transactions").FindOne(context.Background(), bson.M{
		"order_id": orderID,
	}).Decode(&transaction); err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		return false, false
	}

	if _, err := db.Collection("transactions").UpdateOne(context.Background(), bson.M{
		"order_id": orderID,
	}, bson.M{
		"$set": bson.M{
			"payment_status":     RejectedStatus,
			"payment_completion": false,
		},
	}); err != nil {
		return false, transaction.IsDebit
	}

	Notification(transaction.ID, "Invoice", "Payments cancelled with #"+transaction.OrderID)
	return true, transaction.IsDebit
}

func ValidateUser(name, lastName, id, front, back string) bool {
	requestBody := RequestKYC{
		Reference:   "1234561",
		CallbackURL: "https://9999gold.ae",
		Email:       "fasih@icloud.com",
		Country:     "EU",
		Language:    "EN",
		EKYC: EKYC{
			DocumentTwo: DocumentTwo{
				Proof:           front,
				AdditionalProof: back,
				Name: Name{
					FirstName:  name,
					MiddleName: "",
					LastName:   lastName,
				},
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return false
	}

	req, err := http.NewRequest("POST", os.Getenv("KYCUrl"), bytes.NewBuffer(jsonBody))
	if err != nil {
		return false
	}

	credentials := os.Getenv("KYCClient") + ":" + os.Getenv("KYCSecret")
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))
	authorizationHeader := "Basic " + encodedCredentials

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+authorizationHeader)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return false
	}

	var response ResponseKYC
	err = json.Unmarshal([]byte(responseBody), &response)
	if err != nil {
		return false
	}

	verificationResult := response.VerificationResult.EKYC
	return verificationResult == 1
}
