package controlers

import (
	"carat-gold/models"
	"carat-gold/utils"
	"context"
	"encoding/base64"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

func CallCenter(c *gin.Context) {
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

	backData, err := io.ReadAll(backFile)
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	backBase64 := base64.StdEncoding.EncodeToString(backData)

	update := bson.M{
		"documents": models.Documents{
			Back: struct {
				Shot string "json:\"shot\" bson:\"shot\""
			}{
				Shot: backBase64,
			},
			Front: struct {
				Shot string "json:\"shot\" bson:\"shot\""
			}{
				Shot: frontBase64,
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
