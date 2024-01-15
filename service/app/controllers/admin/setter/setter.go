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

func SetSupport(c *gin.Context) {

}

func SetMetaData(c *gin.Context) {
	// Percentage 10
	// Set Wallet
	// Set Bank metas
	// Set F&Q
	// Set Call center datas
}

func SetValidateUser(c *gin.Context) {

}

func SetGeneralData(c *gin.Context) {}

func SetCurrencies(c *gin.Context) {

}

func SetProduct(c *gin.Context) {}

func SetPaymentMethods(c *gin.Context) {}

func SetDeliveryMethods(c *gin.Context) {

}

func SetSymbols(c *gin.Context) {

}

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

func SetUser(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionEditUser) {
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

func SetDeleteUser(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionDeleteUser) {
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

func SetFreezeUser(c *gin.Context) {
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

func SetUnFreezeUser(c *gin.Context) {
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
