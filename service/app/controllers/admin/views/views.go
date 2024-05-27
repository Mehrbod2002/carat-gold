package admin

import (
	"carat-gold/models"
	"carat-gold/utils"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUsersFeedBacks(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionReadOnly) {
		return
	}
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	var feedbacks []models.FeedBacks
	cursor, err := db.Collection("feedbacks").Find(context.Background(), bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &feedbacks); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "done",
		"data":    feedbacks,
	})
}

func ViewAllUsers(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionReadOnly) {
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
	if !models.AllowedAction(c, models.ActionReadOnly) {
		return
	}

	var request struct {
		ID string `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
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

	var purchases []models.Purchased
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
	if !models.AllowedAction(c, models.ActionReadOnly) {
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
	if !models.AllowedAction(c, models.ActionReadOnly) {
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
	if !models.AllowedAction(c, models.ActionReadOnly) {
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
	if !models.AllowedAction(c, models.ActionReadOnly) {
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

	cryptoAll := make([]models.Crypto, 0)
	cryptoAll = append(cryptoAll, crypto)

	debitAll := make([]models.DebitCard, 0)
	debitAll = append(debitAll, debitCard)

	payPalAll := make([]models.PayPal, 0)
	payPalAll = append(payPalAll, payPal)

	applePlayAll := make([]models.ApplePlay, 0)
	applePlayAll = append(applePlayAll, applePlay)

	googlePlayAll := make([]models.GooglePlay, 0)
	googlePlayAll = append(googlePlayAll, googlePlay)
	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "done",
		"crypto":      cryptoAll,
		"debit_card":  debitAll,
		"paypal":      payPalAll,
		"apple_play":  applePlayAll,
		"google_play": googlePlayAll,
	})
}

func ViewDeliveryMethods(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionReadOnly) {
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
	if !models.AllowedAction(c, models.ActionReadOnly) {
		return
	}
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	var generalData []models.GeneralData
	cursor, err := db.Collection("general_data").Find(context.Background(), bson.M{})
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
		"success":      true,
		"message":      "done",
		"general_data": generalData,
	})
}

func ViewMetric(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionReadOnly) {
		return
	}

	var request struct {
		RangeTime int `json:"time"`
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

	if request.RangeTime <= 0 {
		utils.BadBinding(c)
		return
	}

	var generalData models.GeneralData
	if err := db.Collection("general_data").FindOne(context.Background(), bson.M{}).Decode(&generalData); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	var users []models.User
	cursor, err := db.Collection("users").Find(context.Background(), bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &users); err != nil {
		utils.InternalError(c)
		return
	}

	var rangedUsers []models.User
	var totalBars int = 0
	for _, u := range users {
		if u.CreatedAt.Hour() <= request.RangeTime {
			rangedUsers = append(rangedUsers, u)
		}

		for _, p := range u.Wallet.Purchased {
			totalBars += len(p.Product)
		}
	}

	var transactions []models.Transaction
	cursor, err = db.Collection("transactions").Find(context.Background(), bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &transactions); err != nil {
		utils.InternalError(c)
		return
	}

	var aed float64 = 0
	var usd float64 = 0
	for _, t := range transactions {
		if t.CreatedAt.Hour() <= request.RangeTime && t.PaymentCompletion {
			if !t.IsDebit {
				aed += t.TotalPrice * generalData.AedUsd

				if t.PaymentMethod == models.CryptoPayment {
					usd += t.TotalPrice
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "done",
		"users":     rangedUsers,
		"gold_bars": totalBars,
		"aed":       aed,
		"usd":       usd,
	})
}

func ViewFANDQ(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionReadOnly) {
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

func ViewUser(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionReadOnly) {
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

	userID, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
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
	if err := db.Collection("user").FindOne(context.Background(),
		bson.M{"_id": userID}).Decode(&user); err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	var documents models.Documents
	if err := db.Collection("documents").FindOne(context.Background(),
		bson.M{"user_id": userID}).Decode(&documents); err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	var transactions []models.Transaction
	cursor, err := db.Collection("transactions").Find(context.Background(), bson.M{
		"user_id": user.ID,
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
		"success":      true,
		"message":      "done",
		"user":         user,
		"transactions": transactions,
		"documents":    documents,
	})
}

func ViewMetaTrader(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionReadOnly) {
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
	if !models.AllowedAction(c, models.ActionReadOnly) {
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
	if !models.AllowedAction(c, models.ActionReadOnly) {
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
	if !models.AllowedAction(c, models.ActionReadOnly) {
		return
	}

	response, valid := utils.GetRequest("get_history")

	if !valid {
		utils.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response["data"], "message": "done"})
}

func ViewMetaTraderFromWindows(c *gin.Context) {
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
