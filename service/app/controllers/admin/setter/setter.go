package admin

import (
	"carat-gold/models"
	"carat-gold/utils"
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
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

	var request models.RequestSetDefineUser
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadBinding(c)
		return
	}

	valid := request.Validate(c, true)
	if !valid {
		return
	}

	var password []byte
	var passwordSet bool
	if !strings.Contains(*request.Password, ".") {
		passwordSet = true
		hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(*request.Password), bcrypt.DefaultCost)
		if errHash != nil {
			log.Println(errHash)
			c.JSON(401, gin.H{"success": false, "message": "Invalid email or password"})
			return
		}
		password = hashedPassword
	}

	var UserVerified bool = false
	if *request.Status == models.ApprovedStatus {
		UserVerified = true
	}

	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		log.Println(DBerr)
		return
	}

	editedUser := models.User{
		Email:            *request.Email,
		PhoneNumber:      *request.Phone,
		Name:             *request.Name,
		UserVerified:     UserVerified,
		Freeze:           *request.Freeze,
		Permissions:      *request.Permissions,
		IsSupportOrAdmin: *request.IsSupport,
		PhoneVerified:    *request.PhoneVerify,
		StatusString:     *request.Status,
		Reason:           *request.Reason,
		Address:          *request.Address,
		CreatedAt:        time.Now(),
	}

	userID, valid := utils.ValidateID(*request.UserID, c)
	if !valid {
		return
	}

	if passwordSet {
		editedUser.Password = string(password)
	}
	_, err := db.Collection("users").UpdateOne(context.Background(), bson.M{
		"_id": userID,
	}, bson.M{"$set": editedUser})

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
	if _, err := db.Collection("users").DeleteOne(context.Background(), bson.M{
		"_id": user.ID,
	}); err != nil {
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

func SetDefineUser(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionSetUser) {
		return
	}
	var request models.RequestSetDefineUser

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadBinding(c)
		return
	}

	valid := request.Validate(c, false)
	if !valid {
		return
	}

	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		log.Println(DBerr)
		return
	}

	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(*request.Password), bcrypt.DefaultCost)
	if errHash != nil {
		log.Println(errHash)
		c.JSON(401, gin.H{"success": false, "message": "Invalid email or password"})
		return
	}

	filter := bson.M{}
	if *request.Phone != "" {
		filter["phone"] = *request.Phone
	}
	if *request.Email != "" {
		filter["email"] = *request.Email
	}

	var users []models.User
	cursor, err := db.Collection("users").Find(context.Background(), filter)
	if err != nil {
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

	var UserVerified bool = false
	if *request.Status == models.ApprovedStatus {
		UserVerified = true
	}
	if len(users) == 0 {
		newUser := models.User{
			Email:            *request.Email,
			PhoneNumber:      *request.Phone,
			Name:             *request.Name,
			UserVerified:     UserVerified,
			Password:         string(hashedPassword),
			Freeze:           *request.Freeze,
			Permissions:      *request.Permissions,
			IsSupportOrAdmin: *request.IsSupport,
			PhoneVerified:    *request.PhoneVerify,
			StatusString:     *request.Status,
			Reason:           *request.Reason,
			Address:          *request.Address,
			CreatedAt:        time.Now(),
		}
		_, err := db.Collection("users").InsertOne(context.Background(), newUser)
		if err != nil {
			log.Println(err)
			utils.InternalError(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "user already exsits"})
}
