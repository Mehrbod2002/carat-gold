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
	var paymentMethods []models.PaymentMethods
	cursor, err := db.Collection("payment_methods").Find(context.Background(), bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &paymentMethods); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"message":         "done",
		"payment_methods": paymentMethods,
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

func ViewMetaData(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionGeneralDataView) {
		return
	}
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	var metaData []models.MetaData
	cursor, err := db.Collection("meta_datas").Find(context.Background(), bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &metaData); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "done",
		"meta_datas": metaData,
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

func ViewCurrentOrders(c *gin.Context) {
	// if !models.AllowedAction(c, models.ActionMetaTrader) {
	// 	return
	// }
	metaTrader, connected := utils.GetSharedSocket(c)

	if !connected {
		utils.InternalErrorMsg(c, "metatrader connection channel is closed")
		return
	}

	requestID, data := models.GetCurrentOrder()

	n, err := metaTrader.Write([]byte(data))
	if err != nil || n == 0 {
		utils.InternalErrorMsg(c, "metatrader connection channel is closed")
		return
	}

	response, connected := utils.GetSharedReader(c, requestID)

	if !connected {
		utils.InternalErrorMsg(c, "metatrader connection channel is closed")
		return
	}

	if response["status"] == "true" {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": response["data"], "message": response["data"]})
		return
	}

	utils.InternalErrorMsg(c, "metatrader connection channel is closed")
}

func ViewHistoryOrders(c *gin.Context) {
	// if !models.AllowedAction(c, models.ActionMetaTrader) {
	// 	return
	// }
	metaTrader, connected := utils.GetSharedSocket(c)

	if !connected {
		utils.InternalErrorMsg(c, "metatrader connection channel is closed")
		return
	}

	requestID, data := models.GetHistoryOrder()

	n, err := metaTrader.Write([]byte(data))
	if err != nil || n == 0 {
		utils.InternalErrorMsg(c, "metatrader connection channel is closed")
		return
	}

	response, connected := utils.GetSharedReader(c, requestID)

	if !connected {
		utils.InternalErrorMsg(c, "metatrader connection channel is closed")
		return
	}

	if response["status"] == "true" {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": response["data"], "message": "done"})
		return
	}

	utils.InternalErrorMsg(c, "metatrader connection channel is closed")
}
