package controlers

import (
	"carat-gold/models"
	"carat-gold/utils"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetSingelTransaction(c *gin.Context) {
	authUser, valid := models.ValidateSession(c)
	if !valid {
		utils.Unauthorized(c)
		return
	}

	var request struct {
		ID string `json:"id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}

	var transaction models.Transaction
	err = db.Collection("transactions").FindOne(context.Background(), bson.M{"$and": []bson.M{
		{"user_id": authUser.ID},
		{"order_id": request.ID},
	}}).Decode(&transaction)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": utils.Cap("transaction not found"),
			"data":    "not_found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
		"data":    transaction,
	})
}

func GetTransactions(c *gin.Context) {
	authUser, valid := models.ValidateSession(c)
	if !valid {
		utils.Unauthorized(c)
		return
	}

	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}

	var transactions []models.Transaction
	cursor, err := db.Collection("transactions").Find(context.Background(), bson.M{
		"user_id": authUser.ID,
	})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &transactions); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
		"data":    transactions,
	})
}

func GetFANDQ(c *gin.Context) {
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	var fandq []models.FANDQ
	cursor, err := db.Collection("f&q").Find(context.Background(), bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &fandq); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
		"data":    fandq,
	})
}

func GetSymbol(c *gin.Context) {
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	var symbols []models.Symbol
	cursor, err := db.Collection("symbols").Find(context.Background(), bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &symbols); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
		"data":    symbols,
	})
}

func Crisp(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	var callCenter models.CallCenterDatas
	err = db.Collection("call_center").FindOne(context.Background(), bson.M{}).Decode(&callCenter)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	// sessionID := payload["data"].(map[string]interface{})["session_id"].(string)
	// if payload["event"] == "message:notify:unread:received" {
	// 	fmt.Printf("Received Crisp webhook payload: %+v\n", payload)
	// }

	// key := "887b6563dcda97bc5aa17be53af827ebaa3c36326527b249f456394c3bbe4c42"
	// id := "7739a54b-6450-41f2-9e2b-87771e6096ca"

	// endpoint := "https://api.crisp.chat/v1/website/" + *callCenter.LiveChat + "/conversation/" + sessionID
	// /v1/website/{website_id}/conversation/{session_id}/meta
	// client := crisp.WebsiteService{}
	// if callCenter.LiveChat != nil &&
	// 	len(*callCenter.LiveChat) != 0 &&
	// 	len(sessionID) != 0 {
	// 	// fmt.Println(*callCenter.LiveChat, sessionID, 123)
	// 	// a, b, f := client.GetMessagesInConversationLast(*callCenter.LiveChat, sessionID)
	// 	if f != nil {
	// 		// fmt.Println(a, b, f)
	// 	}
	// }
	// fmt.Println(endpoint)
	// fmt.Println("event: ", payload.Payload.Event, payload.Payload.Data.User.UserID, payload.Payload.Data.Type)
}

func GetProducts(c *gin.Context) {
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	var products []models.Products
	cursor, err := db.Collection("products").Find(context.Background(), bson.M{})
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &products); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	var outPut []models.Products
	for _, product := range products {
		if !product.Hide {
			outPut = append(outPut, product)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
		"data":    outPut,
	})
}

func GeneralData(c *gin.Context) {
	_, valid := models.ValidateSession(c)
	if !valid {
		utils.Unauthorized(c)
		return
	}

	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}

	var callCenter models.CallCenterDatas
	var generalData models.GeneralData
	err = db.Collection("call_center").FindOne(context.Background(), bson.M{}).Decode(&callCenter)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	err = db.Collection("general_data").FindOne(context.Background(), bson.M{}).Decode(&generalData)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
		"data":    map[string]any{"call_center": callCenter, "general_data": generalData},
	})
}

func SetCurrency(c *gin.Context) {
	authUser, valid := models.ValidateSession(c)
	if !valid {
		utils.Unauthorized(c)
		return
	}

	var request struct {
		Name string `json:"currency_name"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	var currencies []models.Currency

	cursor, err := db.Collection("currencies").Find(context.Background(), bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &currencies); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	var matched models.Currency
	for _, currency := range currencies {
		if currency.Currency == request.Name {
			matched = currency
			break
		}
	}

	if matched.Currency == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": utils.Cap("invalid currency"),
			"data":    "invalid_currency",
		})
	}

	if _, err := db.Collection("users").UpdateOne(context.Background(), bson.M{
		"_id": authUser.ID,
	}, bson.M{
		"$set": bson.M{
			"currency": matched.Currency,
		},
	}); err != nil {
		utils.BadBinding(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
	})
}

func ViewProducts(c *gin.Context) {
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	var products []models.Products
	cursor, err := db.Collection("products").Find(context.Background(), bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &products); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
		"data":    products,
	})
}

func SetFeedback(c *gin.Context) {
	authUser, valid := models.ValidateSession(c)
	if !valid {
		utils.Unauthorized(c)
		return
	}

	var request struct {
		FeedBack string `json:"feedback"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadBinding(c)
		return
	}

	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		log.Println(DBerr)
		return
	}

	if len(request.FeedBack) > 2000 {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"message": utils.Cap("allowed to send less than 2000 characters"),
		})
		return
	}

	if len(request.FeedBack) < 5 {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"message": utils.Cap("allowed to send more than 5 characters"),
		})
		return
	}

	if count, err := db.Collection("feedbacks").CountDocuments(context.Background(),
		bson.M{
			"user_id": authUser.ID,
		}); err != nil {
		utils.InternalError(c)
		return
	} else if count >= 10 {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"message": utils.Cap("you exceed the max feedbacks"),
		})
		return
	}

	if _, err := db.Collection("feedbacks").InsertOne(context.Background(), models.FeedBacks{
		FeedBack: request.FeedBack,
		UserID:   authUser.ID,
	}); err != nil {
		utils.BadBinding(c)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": utils.Cap("done"),
	})
}

func EditUser(c *gin.Context) {
	authUser, valid := models.ValidateSession(c)
	if !valid {
		utils.Unauthorized(c)
		return
	}
	var request models.RequestEdit

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadBinding(c)
		return
	}

	valid = request.Validate(c)
	if !valid {
		return
	}

	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		log.Println(DBerr)
		return
	}

	updateFields := bson.M{}
	if request.FirstName != nil {
		if len(*request.LastName) != 0 {
			updateFields["first_name"] = utils.TrimAndLowerCase(utils.DerefStringPtr(request.FirstName))
		}
	}
	if request.LastName != nil {
		if len(*request.LastName) != 0 {
			updateFields["last_name"] = utils.TrimAndLowerCase(utils.DerefStringPtr(request.LastName))
		}
	}
	if request.Email != nil {
		updateFields["email"] = utils.TrimAndLowerCase(utils.DerefStringPtr(request.Email))
	}
	if request.Phone != nil {
		updateFields["phone"] = request.Phone
	}
	if request.Address != nil {
		updateFields["address"] = request.Address
	}

	if _, err := db.Collection("users").UpdateOne(context.Background(), bson.M{
		"_id": authUser.ID,
	}, bson.M{
		"$set": updateFields,
	}); err != nil {
		utils.BadBinding(c)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": utils.Cap("done"),
	})
}

func GetUser(c *gin.Context) {
	authUser, valid := models.ValidateSession(c)
	if !valid {
		utils.Unauthorized(c)
		return
	}

	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		log.Println(DBerr)
		return
	}

	var user models.User
	exist := db.Collection("users").FindOne(context.Background(), bson.M{"$and": []bson.M{
		{"_id": authUser.ID},
	}}).Decode(&user)
	if exist != nil {
		if exist == mongo.ErrNoDocuments {
			utils.Unauthorized(c)
			return
		}
		log.Println(exist)
		utils.InternalError(c)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": utils.Cap("done"),
		"data":    user,
	})
}

func SendDocuments(c *gin.Context) {
	authUser, valid := models.ValidateSession(c)
	if !valid {
		utils.Unauthorized(c)
		return
	}

	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		log.Println(DBerr)
		utils.InternalError(c)
		return
	}

	var user models.User
	err := db.Collection("users").FindOne(context.Background(), bson.M{"_id": authUser.ID}).Decode(&user)
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	if user.StatusString == models.PendingStatus {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": utils.Cap("user is awaiting for admin"),
			"data":    "already_registered",
		})
		return
	}

	if user.UserVerified {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": utils.Cap("user already verified"),
			"data":    "already_registered",
		})
		return
	}

	err = c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.Method(c, "front shot size exceeds 10 MB")
		return
	}

	frontFile, _, err := c.Request.FormFile("front")
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer frontFile.Close()

	backFile, _, err := c.Request.FormFile("back")
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer backFile.Close()

	frontData, err := io.ReadAll(frontFile)
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	frontBase64 := base64.StdEncoding.EncodeToString(frontData)
	photoIDFront := primitive.NewObjectID()
	valid = utils.UploadPhoto(c, photoIDFront.Hex(), frontBase64, true)
	if !valid {
		return
	}

	backData, err := io.ReadAll(backFile)
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	backBase64 := base64.StdEncoding.EncodeToString(backData)
	photoIDBack := primitive.NewObjectID()
	valid = utils.UploadPhoto(c, photoIDBack.Hex(), backBase64, true)
	if !valid {
		return
	}

	update := bson.M{
		"user_id":              authUser.ID,
		"documents.back.shot":  photoIDBack.Hex(),
		"documents.front.shot": photoIDFront.Hex(),
	}

	if _, err := db.Collection("documents").
		UpdateOne(context.Background(), bson.M{
			"user_id": user.ID,
		}, bson.M{
			"$set": update,
		}, options.Update().SetUpsert(true)); err != nil {
		utils.BadBinding(c)
		return
	}

	if _, err := db.Collection("users").UpdateOne(context.Background(), bson.M{
		"_id": user.ID,
	}, bson.M{
		"$set": bson.M{
			"user_status": models.PendingStatus,
		},
	}); err != nil {
		utils.BadBinding(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": utils.Cap("document uploaded"),
		"data":    "document_uploaded",
	})
}

func UpdateFcm(c *gin.Context) {
	authUser, valid := models.ValidateSession(c)
	if !valid {
		utils.Unauthorized(c)
		return
	}

	var request struct {
		FcmToken string `json:"fcm_token"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		log.Println(DBerr)
		return
	}

	if _, err := db.Collection("users").UpdateOne(context.Background(), bson.M{
		"_id": authUser.ID,
	}, bson.M{
		"$set": bson.M{
			"fcm_token": request.FcmToken,
		},
	}); err != nil {
		utils.BadBinding(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
	})
}

func CreateTranscations(c *gin.Context) {
	authUser, valid := models.ValidateSession(c)
	if !valid {
		utils.Unauthorized(c)
		return
	}

	var request models.Transaction
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		log.Println(DBerr)
		utils.InternalError(c)
		return
	}

	if request.PaymentMethod != models.CryptoPayment {
		c.JSON(http.StatusNotImplemented, gin.H{
			"success": false,
			"message": utils.Cap("not implement just yet."),
		})
		return
	}

	var deliveryMethods []models.UserDeliveryMethods
	cursor, err := db.Collection("delivery_methods").Find(context.Background(), bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &deliveryMethods); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	var crypto models.Crypto
	err = db.Collection("crypto").FindOne(context.Background(), bson.M{}).Decode(&crypto)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	if len(request.ProductIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": utils.Cap("invalid products"),
		})
		return
	}

	if crypto.Disable {
		c.JSON(http.StatusNotImplemented, gin.H{
			"success": false,
			"message": utils.Cap("payment method disabled"),
		})
		return
	}

	var products []models.Products
	cursor, err = db.Collection("products").Find(context.Background(), bson.M{
		"_id": request,
	})
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	if err := cursor.All(context.Background(), &products); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	// accepted := false
	// for _, i := range deliveryMethods {
	// 	if request.DeliveryMethod == models.DeliveryMethod(i.Title) {
	// 		accepted = true
	// 	}
	// }

	// if !accepted {
	// 	c.JSON(http.StatusNotFound, gin.H{
	// 		"success": false,
	// 		"message": utils.Cap("invalid delivery method"),
	// 	})
	// 	return
	// }

	for _, product := range products {
		if product.Amount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": utils.Cap("there is no more product left"),
			})
		}

		if !product.Hide {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": utils.Cap("invalid product"),
			})
		}
	}

	var orderID = primitive.NewObjectID().Hex()
	_, err = db.Collection("transactions").InsertOne(context.Background(),
		models.Transaction{
			OrderID:           orderID,
			UserID:            authUser.ID,
			CreatedAt:         time.Now(),
			StatusDelivery:    request.StatusDelivery,
			PaymentMethod:     request.PaymentMethod,
			ProductIDs:        request.ProductIDs,
			DeliveryMethod:    request.DeliveryMethod,
			TotalPrice:        request.TotalPrice,
			Vat:               crypto.Vat,
			PaymentCompletion: false,
			PaymentStatus:     models.PendingStatus,
		})
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	pay, valid, errStr := models.CreateCryptoInvoice(c, request.TotalPrice, orderID)
	if !valid {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": errStr,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
		"data":    map[string]string{"order_id": orderID, "url": pay.InvoiceURL},
	})
}

func MakeDepositTransaction(c *gin.Context) {
	authUser, valid := models.ValidateSession(c)
	if !valid {
		utils.Unauthorized(c)
		return
	}

	var request struct {
		PaymentMethod models.PaymentMethod `bson:"payment_method" json:"payment_method"`
		TotalPrice    float64              `bson:"total_price" json:"total_price"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	if request.PaymentMethod != models.CryptoPayment {
		c.JSON(http.StatusNotImplemented, gin.H{
			"success": false,
			"message": utils.Cap("not implement just yet."),
		})
		return
	}

	if request.TotalPrice < 25 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": utils.Cap("Deposit price should be greater than 99$"),
			"data":    "invalid_price",
		})
		return
	}

	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}

	var cryptoDetail models.Crypto
	err = db.Collection("crypto").FindOne(context.Background(), bson.M{}).Decode(&cryptoDetail)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	var orderID = primitive.NewObjectID().Hex()
	_, err = db.Collection("transactions").InsertOne(context.Background(),
		models.Transaction{
			OrderID:           orderID,
			UserID:            authUser.ID,
			CreatedAt:         time.Now(),
			PaymentMethod:     request.PaymentMethod,
			Vat:               cryptoDetail.Vat,
			TotalPrice:        request.TotalPrice,
			IsDebit:           true,
			PaymentStatus:     models.PendingStatus,
			PaymentCompletion: false,
		})
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	pay, valid, errStr := models.CreateCryptoInvoice(c, request.TotalPrice+cryptoDetail.Vat, orderID)
	if !valid {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": errStr,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
		"data":    map[string]string{"order_id": orderID, "url": pay.InvoiceURL},
	})
}

func Cancel(c *gin.Context) {
	authUser, valid := models.ValidateSession(c)
	if !valid {
		utils.Unauthorized(c)
		return
	}

	var request struct {
		OrderID string `json:"order_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}

	var transaction models.Transaction
	err = db.Collection("transactions").FindOne(context.Background(), bson.M{"$and": []bson.M{
		{"user_id": authUser.ID},
		{"order_id": request.OrderID},
	}}).Decode(&transaction)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": utils.Cap("transaction not found"),
			"data":    "not_found",
		})
		return
	}

	if _, err = db.Collection("transactions").UpdateOne(context.Background(), bson.M{"$and": []bson.M{
		{"user_id": authUser.ID},
		{"order_id": request.OrderID},
	}}, bson.M{
		"$set": bson.M{
			"payment_status": models.RejectedStatus,
		},
	}); err != nil {
		log.Println(err)
		utils.InternalError(c)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
	})
}

func RevalidateToken(c *gin.Context) {
	session := sessions.Default(c)
	authUser, valid := models.ValidateSession(c)

	if !valid {
		return
	}

	var request struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}

	var user models.User
	if err = db.Collection("users").FindOne(context.Background(), bson.M{
		"_id":           authUser.ID,
		"refresh_token": request.RefreshToken,
	}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": utils.Cap("user not found"),
			})
			return
		}
		log.Println(err)
		utils.InternalError(c)
		return
	}

	session.Delete("token")
	err = session.Save()
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	token, err := user.GenerateToken()
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	refreshToken, err := user.GenerateToken()
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	_, err = db.Collection("users").UpdateOne(context.Background(), bson.M{
		"_id": user.ID,
	}, bson.M{"$set": bson.M{
		"refresh_token": refreshToken,
	}})
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	session.Set("token", token)
	err = session.Save()
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
		"data": map[string]string{
			"token":         token,
			"refresh_token": refreshToken,
		}})
}

func UserPaymentMethods(c *gin.Context) {
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}

	// var googlePlay models.GooglePlay
	// var applePlay models.ApplePlay
	// var payPal models.PayPal
	var debitCard []models.UserDebitCard
	var crypto []models.UserCrypto
	// err = db.Collection("google_play").FindOne(context.Background(), bson.M{}).Decode(&googlePlay)
	// if err != nil && err != mongo.ErrNoDocuments {
	// 	log.Println(err)
	// 	utils.InternalError(c)
	// 	return
	// }
	// err = db.Collection("apple_play").FindOne(context.Background(), bson.M{}).Decode(&applePlay)
	// if err != nil && err != mongo.ErrNoDocuments {
	// 	log.Println(err)
	// 	utils.InternalError(c)
	// 	return
	// }
	// err = db.Collection("paypal").FindOne(context.Background(), bson.M{}).Decode(&payPal)
	// if err != nil && err != mongo.ErrNoDocuments {
	// 	log.Println(err)
	// 	utils.InternalError(c)
	// 	return
	// }
	cursor, err := db.Collection("debit_card").Find(context.Background(), bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &debitCard); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	cursor, err = db.Collection("crypto").Find(context.Background(), bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &crypto); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	// payPalAll := make([]models.PayPal, 0)
	// payPalAll = append(payPalAll, payPal)

	// applePlayAll := make([]models.ApplePlay, 0)
	// applePlayAll = append(applePlayAll, applePlay)

	// googlePlayAll := make([]models.GooglePlay, 0)
	// googlePlayAll = append(googlePlayAll, googlePlay)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
		"data": map[string]any{
			"crypto":     crypto,
			"debit_card": debitCard,
			// "paypal":      payPalAll,
			// "apple_play":  applePlayAll,
			// "google_play": googlePlayAll,
		},
	})
}

func UserDeliveryMethods(c *gin.Context) {
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	var deliveryMethods []models.UserDeliveryMethods
	cursor, err := db.Collection("delivery_methods").Find(context.Background(), bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &deliveryMethods); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	var finalDeliveryMethods []models.UserDeliveryMethods
	for _, deliveryMethod := range deliveryMethods {
		if !deliveryMethod.Disable {
			finalDeliveryMethods = append(finalDeliveryMethods, deliveryMethod)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
		"data":    finalDeliveryMethods,
	})
}

func PayWithWallet(c *gin.Context) {
	authUser, valid := models.ValidateSession(c)
	if !valid {
		utils.Unauthorized(c)
		return
	}

	var request models.Transaction
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		log.Println(DBerr)
		utils.InternalError(c)
		return
	}

	if request.PaymentMethod != models.CryptoPayment && request.PaymentMethod != models.WalletPayment {
		c.JSON(http.StatusNotImplemented, gin.H{
			"success": false,
			"message": utils.Cap("not implement just yet."),
		})
		return
	}

	if len(request.ProductIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": utils.Cap("invalid products"),
		})
		return
	}

	var deliveryMethods []models.UserDeliveryMethods
	cursor, err := db.Collection("delivery_methods").Find(context.Background(), bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &deliveryMethods); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	var crypto models.Crypto
	err = db.Collection("crypto").FindOne(context.Background(), bson.M{}).Decode(&crypto)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	var user models.User
	err = db.Collection("users").FindOne(context.Background(), bson.M{
		"_id": authUser.ID,
	}).Decode(&user)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	if request.TotalPrice == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": utils.Cap("invalid price"),
		})
		return
	}

	if (user.Wallet.BalanceUSD < request.TotalPrice) &&
		request.TotalPrice != 0 &&
		user.Wallet.BalanceUSD != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": utils.Cap("insufficient balance"),
		})
		return
	}

	if request.IsAED {
		var generalData models.GeneralData
		err = db.Collection("general_data").FindOne(context.Background(), bson.M{}).Decode(&generalData)
		if err != nil && err != mongo.ErrNoDocuments {
			log.Println(err)
			utils.InternalError(c)
			return
		}

		request.TotalPrice = request.TotalPrice / generalData.AedUsd
	}

	// result, valid := utils.GetRequest("get_last_price")
	// if !valid {
	// 	utils.AdminError(c)
	// 	return
	// }
	// lastGoldPrice := result["data"].(float64)
	// lengths := float64(len(request.ProductIDs))
	// eachGoldBar := request.TotalPrice / lengths
	// if (-10 < eachGoldBar-lastGoldPrice || eachGoldBar-lastGoldPrice > 10) &&
	// 	request.TotalPrice != 0 {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"success": false,
	// 		"message": utils.Cap("each gold bar price is more than 10 USD difference"),
	// 	})
	// 	return
	// }

	if crypto.Disable {
		c.JSON(http.StatusNotImplemented, gin.H{
			"success": false,
			"message": utils.Cap("payment method disabled"),
		})
		return
	}

	var products []models.Products
	cursor, err = db.Collection("products").Find(context.Background(), bson.M{
		"_id": request,
	})
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	if err := cursor.All(context.Background(), &products); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	accepted := false
	for _, i := range deliveryMethods {
		if request.DeliveryMethod == models.DeliveryMethod(i.Title) {
			accepted = true
		}
	}

	if !accepted {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": utils.Cap("invalid delivery method"),
		})
		return
	}

	for _, product := range products {
		if product.Amount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": utils.Cap("there is no more product left"),
			})
		}

		if !product.Hide {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": utils.Cap("invalid product"),
			})
		}
	}

	var orderID = primitive.NewObjectID().Hex()
	mID, valid := utils.AutoOrder(request.TotalPrice)
	if valid {
		models.StoreMetatraderID(orderID, fmt.Sprintf("%d", mID))
	}
	if !valid {
		utils.AdminError(c)
		return
	}

	_, err = db.Collection("transactions").InsertOne(context.Background(),
		models.Transaction{
			OrderID:           orderID,
			UserID:            authUser.ID,
			CreatedAt:         time.Now(),
			StatusDelivery:    request.StatusDelivery,
			PaymentMethod:     request.PaymentMethod,
			ProductIDs:        request.ProductIDs,
			DeliveryMethod:    request.DeliveryMethod,
			TotalPrice:        request.TotalPrice,
			Vat:               crypto.Vat,
			PaymentCompletion: true,
			PaymentStatus:     models.ApprovedStatus,
		})
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	if _, err := db.Collection("users").UpdateOne(context.Background(), bson.M{
		"_id": authUser.ID,
	}, bson.M{
		"$inc": bson.M{
			"wallet.balance": -request.TotalPrice,
		},
	}); err != nil {
		utils.InternalError(c)
		return
	}

	if user.Wallet.Purchased == nil {
		user.Wallet.Purchased = []models.Purchased{}
	}

	newPurchase := models.Purchased{
		PaymentStatus:  models.ApprovedStatus,
		PaymentMethd:   request.PaymentMethod,
		CreatePayment:  time.Now(),
		CreatedAt:      time.Now(),
		StatusDelivery: "",
		Product:        request.ProductIDs,
		OrderID:        orderID,
	}

	user.Wallet.Purchased = append(user.Wallet.Purchased, newPurchase)

	update := bson.M{
		"$set": bson.M{
			"wallet": user.Wallet,
		},
	}

	if _, err := db.Collection("users").UpdateOne(context.Background(), bson.M{
		"_id": authUser.ID,
	}, update); err != nil {
		utils.InternalError(c)
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
		"data":    map[string]string{"order_id": orderID},
	})
}
