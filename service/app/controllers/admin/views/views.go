package admin

import (
	"carat-gold/models"
	"carat-gold/utils"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ViewAllUsers(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionGeneralDataView) {
		return
	}
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	var users []models.User
	cursor, err := db.Collection("users").Find(context.Background(), bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &users); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
		"users":   users,
	})
}

func ViewPurchase(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionGeneralDataView) {
		return
	}

	var request struct {
		ID string `json:"user_id"`
	}

	purchaseID, valid := utils.ValidateID(request.ID, c)
	if !valid {
		return
	}

	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}

	var purchases []models.Purchaed
	cursor, err := db.Collection("purchases").Find(context.Background(), bson.M{
		"user_id": purchaseID,
	})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &purchases); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "done",
		"purchases": purchases,
	})
}

func ViewCurrencies(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionGeneralDataView) {
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
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "done",
		"currencies": currencies,
	})
}

func ViewProducts(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionGeneralDataView) {
		return
	}
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

func ViewSymbols(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionGeneralDataView) {
		return
	}
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
		"symbols": symbols,
	})
}

func ViewPaymentMethods(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionGeneralDataView) {
		return
	}
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}

	var googlePlay models.GooglePlay
	var applePlay models.ApplePlay
	var payPal models.PayPal
	var debitCard models.DebitCard
	var crypto models.Crypto
	err = db.Collection("google_play").FindOne(context.Background(), bson.M{}).Decode(&googlePlay)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	err = db.Collection("apple_play").FindOne(context.Background(), bson.M{}).Decode(&applePlay)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	err = db.Collection("paypal").FindOne(context.Background(), bson.M{}).Decode(&payPal)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	err = db.Collection("debit_card").FindOne(context.Background(), bson.M{}).Decode(&debitCard)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	err = db.Collection("crypto").FindOne(context.Background(), bson.M{}).Decode(&crypto)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "done",
		"crypto":      crypto,
		"debit_card":  debitCard,
		"paypal":      payPal,
		"apple_play":  applePlay,
		"google_play": googlePlay,
	})
}

func ViewDeliveryMethods(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionGeneralDataView) {
		return
	}
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	var deliveryMethods []models.DeliveryMethods
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
	c.JSON(http.StatusOK, gin.H{
		"success":          true,
		"message":          "done",
		"delivery_methods": deliveryMethods,
	})
}

func ViewGeneralData(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionGeneralDataView) {
		return
	}
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	var generalData []models.GeneralData
	cursor, err := db.Collection("general_datas").Find(context.Background(), bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &generalData); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"message":       "done",
		"general_datas": generalData,
	})
}

func ViewChatsHistories(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionGeneralDataView) {
		return
	}
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	var chatHistories []models.ChatHistories
	cursor, err := db.Collection("chat_histories").Find(context.Background(), bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &chatHistories); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"message":        "done",
		"chat_histories": chatHistories,
	})
}

func ViewFANDQ(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionGeneralDataView) {
		return
	}
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
		"fandq":   fandq,
	})
}

func ViewMetaTrader(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionMetaTrader) {
		return
	}

	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		log.Println(DBerr)
		return
	}

	var metaTraderAccounts models.MetaTraderAccounts
	err := db.Collection("metatrader_accounts").FindOne(context.Background(), bson.M{}, options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}})).Decode(&metaTraderAccounts)
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "done",
		"accounts": metaTraderAccounts,
	})
}

func ViewCallCenter(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionGeneralDataView) {
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

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "done",
		"call_center": callCenter,
	})
}

func ViewCurrentOrders(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionMetaTrader) {
		return
	}

	response, valid := utils.GetRequest("positions")

	if !valid {
		utils.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response["data"], "message": "done"})
}

func ViewHistoryOrders(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionMetaTrader) {
		return
	}

	response, valid := utils.GetRequest("get_history")

	if !valid {
		utils.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response["data"], "message": "done"})
}
