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

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetProducts(c *gin.Context) {
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
		"success":  true,
		"message":  "done",
		"products": products,
	})
}

func GeneralData(c *gin.Context) {
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
		"success":      true,
		"message":      "done",
		"call_center":  callCenter,
		"general_data": generalData,
	})
}

func SetCurrency(c *gin.Context) {
	authUser, _ := models.ValidateSession(c)

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
		"success":  true,
		"message":  "done",
		"products": products,
	})
}

func EditUser(c *gin.Context) {
	authUser, _ := models.ValidateSession(c)
	var request models.RequestEdit

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	valid := request.Validate(c)
	if !valid {
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
			"phone":   request.Phone,
			"address": request.Address,
			"name":    request.Name,
			"email":   request.Email,
		},
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
	authUser, _ := models.ValidateSession(c)

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
	authUser, _ := models.ValidateSession(c)

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
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": utils.Cap("user in processing"),
			"data":    "already_registered",
		})
		return
	}

	if user.UserVerified {
		c.JSON(http.StatusUnauthorized, gin.H{
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
	valid := utils.UploadPhoto(c, photoIDFront.Hex(), frontBase64)
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
	valid = utils.UploadPhoto(c, photoIDBack.Hex(), backBase64)
	if !valid {
		return
	}

	update := bson.M{
		"documents": models.Documents{
			Back: struct {
				Shot string "json:\"shot\" bson:\"shot\""
			}{
				Shot: photoIDBack.Hex(),
			},
			Front: struct {
				Shot string "json:\"shot\" bson:\"shot\""
			}{
				Shot: photoIDFront.Hex(),
			},
		},
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
	authUser, _ := models.ValidateSession(c)

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
	authUser, _ := models.ValidateSession(c)

	var request models.Transaction
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

	if request.PaymentMethod != models.CryptoPayment {
		c.JSON(http.StatusNotImplemented, gin.H{
			"success": false,
			"message": utils.Cap("not implement just yet."),
		})
		return
	}

	var crypto models.Crypto
	err := db.Collection("crypto").FindOne(context.Background(), bson.M{}).Decode(&crypto)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	if !crypto.Disable {
		c.JSON(http.StatusNotImplemented, gin.H{
			"success": false,
			"message": utils.Cap("payment method disabled"),
		})
		return
	}

	var products []models.Products
	cursor, err := db.Collection("products").Find(context.Background(), bson.M{
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
			TotalPrice:        request.TotalPrice,
			Vat:               crypto.Vat,
			PaymentCompletion: false,
		})
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	pay, err := models.CreateCrypto(c, request.TotalPrice+crypto.Vat, orderID, crypto.Token)
	if err != nil {
		utils.InternalError(c)
		return
	}

	qrCode := fmt.Sprintf("ethereum:%s?value=%f&token=usdt", pay.PayAddress, pay.PayAmount)
	url := fmt.Sprintf("https://nowpayments.io/payment/?iid=%s&paymentId=%s", pay.PurchaseID, pay.PaymentID)
	qr, err := models.CreateQr(qrCode)
	if err != nil {
		utils.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
		"data":    orderID,
		"url":     url,
		"qr":      qr,
	})
}

func MakeDepositTransaction(c *gin.Context) {
	authUser, _ := models.ValidateSession(c)

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

	if request.TotalPrice < 100 {
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
			PaymentCompletion: false,
		})
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	pay, err := models.CreateCrypto(c, request.TotalPrice+cryptoDetail.Vat, orderID, cryptoDetail.Token)
	if err != nil {
		utils.InternalError(c)
		return
	}

	qrCode := fmt.Sprintf("ethereum:%s?value=%f&token=usdt", pay.PayAddress, pay.PayAmount)
	url := fmt.Sprintf("https://nowpayments.io/payment/?iid=%s&paymentId=%s", pay.PurchaseID, pay.PaymentID)
	qr, err := models.CreateQr(qrCode)
	if err != nil {
		utils.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
		"data":    orderID,
		"url":     url,
		"qr":      qr,
	})
}

func Pay(c *gin.Context) {
	authUser, _ := models.ValidateSession(c)

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

	var transction models.Transaction
	if err = db.Collection("transactions").FindOne(context.Background(), bson.M{
		"$and": bson.M{
			"order_id": request.OrderID,
			"user_id":  authUser.ID,
		},
	}).Decode(&transction); err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	if transction.IsDebit {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": utils.Cap("This is debit operation"),
			"data":    "invalid_operation",
		})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"$push": bson.M{
				"wallet.purchased": models.Purchased{
					PaymentStatus:  transction.PaymentStatus,
					PaymentMethd:   transction.PaymentMethod,
					CreatePayment:  time.Now(),
					CreatedAt:      transction.CreatedAt,
					StatusDelivery: transction.StatusDelivery,
					Product:        transction.ProductIDs,
					OrderID:        transction.OrderID,
				},
			},
		},
	}

	if transction.StatusDelivery == models.Hold {
		update["$inc"] = bson.M{
			"wallet.balance": transction.TotalPrice - transction.Vat,
		}
	}

	if _, err := db.Collection("transactions").UpdateOne(context.Background(), bson.M{
		"$and": bson.M{
			"order_id": request.OrderID,
			"user_id":  authUser.ID,
		},
	}, bson.M{
		"$set": bson.M{
			"payment_status":     models.ApprovedStatus,
			"payment_completion": true,
		},
	}); err != nil {
		utils.BadBinding(c)
		return
	}

	if _, err := db.Collection("users").UpdateOne(context.Background(), bson.M{
		"_id": authUser.ID,
	}, update); err != nil {
		utils.BadBinding(c)
		return
	}

	utils.AutoOrder(c, 0)
	models.Notification(c, authUser.ID, "Payments set with #"+transction.OrderID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
	})
}
