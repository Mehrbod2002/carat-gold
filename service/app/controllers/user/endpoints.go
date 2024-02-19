package controlers

import (
	"carat-gold/models"
	"carat-gold/utils"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

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

	var documents models.Documents
	if err := c.ShouldBindJSON(&documents); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	err := documents.Validate(c)
	if !err {
		log.Println(err)
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
		log.Println(exist)
		utils.InternalError(c)
		return
	}

	if user.StatusString == models.PendingStatus {
		c.JSON(401, gin.H{
			"success": false,
			"message": utils.Cap("user in proccessing"),
			"data":    "already_registered",
		})
		return
	}

	if user.UserVerified {
		c.JSON(401, gin.H{
			"success": false,
			"message": utils.Cap("user already verified"),
			"data":    "already_registered",
		})
		return
	}

	var update = bson.M{
		"documents": documents,
	}

	if (documents.Side == "front" && user.Documents.Side == "back") ||
		(documents.Side == "back" && user.Documents.Side == "front") {
		update["user_status"] = models.PendingStatus
	}

	if _, err := db.Collection("users").
		UpdateOne(context.Background(), bson.M{
			"_id": user.ID,
		}, bson.M{
			"$set": update,
		}); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": utils.Cap("document uploaded"),
		"data":    "document_uploaded",
	})
}
