package controllers

import (
	"carat-gold/models"
	"carat-gold/utils"
	"context"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetUserPermissions(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionSetPermission) {
		return
	}
	var permissions struct {
		UserID     string            `bson:"user_id"`
		Permission models.Permission `bson:"permissions" json:"permissions"`
	}
	userID, err := primitive.ObjectIDFromHex(permissions.UserID)
	if err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}
	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		log.Println(DBerr)
		return
	}
	_, err = db.Collection("users").UpdateOne(context.Background(), bson.M{
		"_id": userID,
	}, bson.M{"$addToSet": permissions}, options.Update().SetUpsert(true))
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
}

func EditUser(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}
	var user models.User
	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		log.Println(DBerr)
		return
	}
	_, err := db.Collection("users").UpdateOne(context.Background(), bson.M{
		"_id": user.ID,
	}, bson.M{"$set": user})
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
}

func DeleteUser(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionDelete) {
		return
	}
	var request struct {
		ID string `json:"user_id"`
	}
	if err := c.ShouldBindQuery(&request); err != nil {
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
	var user models.User
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	if err := db.Collection("users").FindOne(context.Background(), bson.M{
		"_id": userID,
	}).Decode(&user); err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	if _, err := db.Collection("users").DeleteOne(context.Background(), user); err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "done"})
}

func AdminLogout(c *gin.Context) {
	models.ValidateSession(c)
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Logging failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}

func GetAllUsers(c *gin.Context) {
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

func FreezeUser(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionFreeUser) {
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
	var user models.User
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	userID, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}
	if err := db.Collection("users").FindOne(context.Background(), bson.M{
		"_id": userID,
	}).Decode(&user); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	if user.Freeze {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "freezed before"})
		return
	}
	if _, err := db.Collection("users").UpdateOne(context.Background(), bson.M{"$set": bson.M{"_id": user.ID}}, bson.M{
		"freeze": true,
	}); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "done"})
}

func UnFreezeUse(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionUnfreezeUser) {
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
	userID, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}
	var user models.User
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	if err := db.Collection("users").FindOne(context.Background(), bson.M{
		"_id": userID,
	}).Decode(&user); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	if !user.Freeze {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "unfreezed before"})
		return
	}
	if _, err := db.Collection("users").UpdateOne(context.Background(), bson.M{"$set": bson.M{"_id": user.ID}}, bson.M{
		"freeze": false,
	}); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "done"})
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

func SetCurrencies(c *gin.Context) {

}

func SetProduct(c *gin.Context) {}

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

func SetSymbols(c *gin.Context) {

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
	var symbols []models.Symbols
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

func SetPaymentMethods(c *gin.Context) {

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

func SetDeliveryMethods(c *gin.Context) {

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

func ValidateUser(c *gin.Context) {

}

func SetGeneralData(c *gin.Context) {}

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

func ViewHistoryOrders(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionGeneralDataView) {
		return
	}
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	var historyOrders []models.HistoryOrders
	cursor, err := db.Collection("history_orders").Find(context.Background(), bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &historyOrders); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"message":        "done",
		"history_orders": historyOrders,
	})
}

func ViewCurrentOrders(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionGeneralDataView) {
		return
	}
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	var realTimeOrders []models.RealTimeOrders
	cursor, err := db.Collection("real_time_orders").Find(context.Background(), bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	defer cursor.Close(context.Background())
	if err := cursor.All(context.Background(), &realTimeOrders); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":          true,
		"message":          "done",
		"real_time_orders": realTimeOrders,
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

func SetMetaData(c *gin.Context) {
	// Percentage 10
	// Set Wallet
	// Set Bank metas
	// Set F&Q
	// Set Call center datas
}

func VerifyUser(c *gin.Context) {

}
