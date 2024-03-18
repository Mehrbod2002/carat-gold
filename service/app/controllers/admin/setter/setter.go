package admin

import (
	"carat-gold/models"
	"carat-gold/utils"
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func SetDeletePayment(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}
	var request struct {
		ID   string `json:"id"`
		Side string `json:"side"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	paymentID, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	var payment models.DefaultPayment
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}

	if err := db.Collection(request.Side).FindOne(context.Background(), bson.M{
		"_id": paymentID,
	}).Decode(&payment); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	if _, err := db.Collection(request.ID).DeleteOne(context.Background(), bson.M{
		"_id": payment.ID,
	}); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "done"})
}

func SetEditPayment(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}

	var request models.RequestSetPayment
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.BadBinding(c)
		return
	}

	valid := request.Validate(c)
	if !valid {
		return
	}

	paymentID, valid := utils.ValidateID(*request.ID, c)
	if !valid {
		return
	}

	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		log.Println(DBerr)
		return
	}

	editPayment := models.DefaultPayment{
		Token:   request.Token,
		Address: request.Address,
		Access:  request.Access,
		Vat:     request.Vat,
	}

	_, err := db.Collection(request.Side).UpdateOne(context.Background(), bson.M{
		"_id": paymentID,
	}, bson.M{"$set": editPayment})

	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
}

func SetPayment(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}
	user, _ := models.ValidateSession(c)
	var request *models.RequestSetPayment

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

	_, err := db.Collection(request.Side).InsertOne(context.Background(),
		models.DefaultPayment{
			Vat:       request.Vat,
			Address:   request.Address,
			Access:    request.Access,
			Token:     request.Token,
			WhoDefine: user.Email,
		})
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
}

func SetEditFANDQ(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionContent) {
		return
	}

	var request models.RequestSetFANDQ
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	fandqID, valid := utils.ValidateID(*request.ID, c)
	if !valid {
		return
	}

	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		log.Println(DBerr)
		return
	}

	var fandq models.FANDQ
	err := db.Collection("f&q").FindOne(context.Background(),
		bson.M{"_id": fandqID},
	).Decode(&fandq)

	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	_, err = db.Collection("f&q").UpdateOne(context.Background(), bson.M{
		"_id": fandqID,
	}, bson.M{"$set": request})

	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
}

func SetFANDQ(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}
	user, _ := models.ValidateSession(c)
	var request models.RequestSetFANDQ
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
	_, err := db.Collection("f&q").InsertOne(context.Background(),
		models.FANDQ{
			Question:  request.Question,
			Answer:    request.Answer,
			WhoDefine: user.Email,
			CreatedAt: time.Now(),
		})
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
}

func SetDeleteFANDQ(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}
	var request struct {
		ID string `json:"fandq_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	fandqID, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	var fandq models.FANDQ
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}

	if err := db.Collection("f&q").FindOne(context.Background(), bson.M{
		"_id": fandqID,
	}).Decode(&fandq); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	if _, err := db.Collection("f&q").DeleteOne(context.Background(), bson.M{
		"_id": fandq.ID,
	}); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "done"})
}

func SetAedExchange(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}

	var request models.RequestSetGeneralData
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		utils.InternalError(c)
		return
	}

	result, err := db.Collection("general_data").UpdateOne(context.Background(), bson.M{}, bson.M{
		"$set": request,
	}, options.Update().SetUpsert(true))
	if err != nil {
		utils.InternalError(c)
		return
	}

	if result.UpsertedCount > 0 {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "General data inserted"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "General data updated"})
	}
}

func SetCallCenterDatas(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}

	var request models.RequestSetCallCenter
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		utils.InternalError(c)
		return
	}

	result, err := db.Collection("call_center").UpdateOne(context.Background(), bson.M{}, bson.M{
		"$set": request,
	}, options.Update().SetUpsert(true))
	if err != nil {
		utils.InternalError(c)
		return
	}

	if result.UpsertedCount > 0 {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Document inserted"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Document updated"})
	}
}

func SetMetaData(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}

	var request models.RequestMetaTraderAccounts
	if err := c.ShouldBindJSON(&request); err != nil {
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

	_, err := db.Collection("metatrader_accounts").UpdateOne(context.Background(), bson.M{}, bson.M{
		"$set": request,
	}, options.Update().SetUpsert(true))
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	result, valid := utils.PostRequest(make(map[string]interface{}), "reinitialize")
	if !valid || result["status"] == false {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": result["data"], "data": result["data"]})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "done"})
}

func SetEditCurrency(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}

	var request struct {
		ID string `json:"currency_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	currencyID, valid := utils.ValidateID(request.ID, c)
	if !valid {
		return
	}

	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		log.Println(DBerr)
		return
	}

	var currency models.Currency
	err := db.Collection("currencies").FindOne(context.Background(),
		bson.M{"_id": currencyID},
	).Decode(&currency)

	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	_, err = db.Collection("currencies").UpdateOne(context.Background(), bson.M{
		"_id": currencyID,
	}, bson.M{"$set": bson.M{"active": !currency.Active}})

	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
}

func SetEditProduct(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}

	var images []models.Image
	var request models.RequestSetProduct
	if err := c.ShouldBindJSON(&request); err != nil {
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

	for _, photo := range request.Images {
		photoID := primitive.NewObjectID()
		valid := utils.UploadPhoto(c, photoID.Hex(), photo)
		if !valid {
			return
		}

		images = append(images, models.Image{PhotoID: photoID})
	}

	editedUser := models.Products{
		Percentage:  request.Percentage,
		Name:        request.Name,
		Description: request.Description,
		WeightOZ:    request.WeightOZ,
		WeightGramm: request.WeightGramm,
		Purity:      request.Purity,
		Length:      request.Length,
		Width:       request.Width,
		Amount:      request.Amount,
		Images:      images,
	}

	productId, valid := utils.ValidateID(*request.ProductID, c)
	if !valid {
		return
	}

	_, err := db.Collection("products").UpdateOne(context.Background(), bson.M{
		"_id": productId,
	}, bson.M{"$set": editedUser})

	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
}

func SetProduct(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}
	user, _ := models.ValidateSession(c)
	var request models.RequestSetProduct
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

	var images []models.Image
	for _, photo := range request.Images {
		photoID := primitive.NewObjectID()
		valid := utils.UploadPhoto(c, photoID.Hex(), photo)
		if !valid {
			return
		}

		images = append(images, models.Image{PhotoID: photoID})
	}

	_, err := db.Collection("products").InsertOne(context.Background(),
		models.Products{
			Name:        request.Name,
			Description: request.Description,
			WeightOZ:    request.WeightOZ,
			WeightGramm: request.WeightGramm,
			Purity:      request.Purity,
			Length:      request.Length,
			Width:       request.Width,
			Percentage:  request.Percentage,
			Images:      images,
			Amount:      request.Amount,
			WhoDefine:   user.Email,
			CreatedAt:   time.Now(),
		})
	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
}

func SetSymbols(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}

	var request models.RequestSetSymbol
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

	var images []models.Image
	photoID := primitive.NewObjectID()
	valid := utils.UploadPhoto(c, photoID.Hex(), request.Image)
	if !valid {
		return
	}

	images = append(images, models.Image{PhotoID: photoID})
	var symbols []models.Symbol
	cursor, err := db.Collection("symbols").Find(context.Background(),
		bson.M{"name": request.Name})
	if err != nil {
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

	if len(symbols) == 0 {
		newSymbol := models.Symbol{
			SymbolName: request.Name,
			Images:     images,
			CreatedAt:  time.Now(),
		}
		_, err := db.Collection("symbols").InsertOne(context.Background(), newSymbol)
		if err != nil {
			log.Println(err)
			utils.InternalError(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "symbol already exsits"})
}

func SetUserPermissions(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
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
	if !models.AllowedAction(c, models.ActionWrite) {
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
		Email:            utils.DerefStringPtr(request.Email),
		PhoneNumber:      utils.DerefStringPtr(request.Phone),
		Name:             utils.DerefStringPtr(request.Name),
		UserVerified:     UserVerified,
		Freeze:           utils.DerefBoolPtr(request.Freeze),
		Permissions:      *request.Permissions,
		IsSupportOrAdmin: utils.DerefBoolPtr(request.IsSupport),
		PhoneVerified:    utils.DerefBoolPtr(request.PhoneVerify),
		StatusString:     *request.Status,
		Reason:           utils.DerefStringPtr(request.Reason),
		Address:          utils.DerefStringPtr(request.Address),
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

func SetDeleteProduct(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}

	var request struct {
		ID string `json:"product_id"`
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

	var product models.Products
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	if err := db.Collection("products").FindOne(context.Background(), bson.M{
		"_id": userID,
	}).Decode(&product); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	for _, image := range product.Images {
		path := filepath.Join("CDN", image.PhotoID.Hex()+".png")
		os.Remove(path)
	}

	if _, err := db.Collection("products").DeleteOne(context.Background(), bson.M{
		"_id": product.ID,
	}); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "done"})
}

func SetDeleteUser(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
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
	if !models.AllowedAction(c, models.ActionWrite) {
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
	if !models.AllowedAction(c, models.ActionWrite) {
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
	if !models.AllowedAction(c, models.ActionWrite) {
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
	if request.Phone != nil {
		filter["phone"] = *request.Phone
	}
	if request.Email != nil {
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
			Currency:         "USD",
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

func SetEditDeliveryMethods(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}
	var request models.RequestSetDeliveryMethod
	if err := c.ShouldBindJSON(&request); err != nil {
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

	deliveryID, valid := utils.ValidateID(*request.DeliveryID, c)
	if !valid {
		return
	}

	user, _ := models.ValidateSession(c)
	editedDelivery := models.DeliveryMethods{
		Title:         request.Title,
		EstimatedTime: request.EstimatedTime,
		TimeProvided:  request.TimeProvided,
		Description:   request.Description,
		Fee:           request.Fee,
		WhoDefine:     user.Email,
		CreatedAt:     time.Now(),
		Disable:       request.Disable,
	}

	_, err := db.Collection("delivery_methods").UpdateOne(context.Background(), bson.M{
		"_id": deliveryID,
	}, bson.M{"$set": editedDelivery})

	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
}

func SetDeliveryMethods(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}
	user, _ := models.ValidateSession(c)
	var request *models.RequestSetDeliveryMethod
	if err := c.ShouldBindJSON(&request); err != nil {
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

	_, err := db.Collection("delivery_methods").InsertOne(context.Background(),
		models.DeliveryMethods{
			Title:         request.Title,
			EstimatedTime: request.EstimatedTime,
			TimeProvided:  request.TimeProvided,
			Description:   request.Description,
			Fee:           request.Fee,
			WhoDefine:     user.Email,
			CreatedAt:     time.Now(),
			Disable:       request.Disable,
		})

	if err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
}

func SetDeleteSymbol(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}
	var request struct {
		ID string `json:"symbol_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	symolID, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	var symbol models.Symbol
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}

	if err := db.Collection("symbols").FindOne(context.Background(), bson.M{
		"_id": symolID,
	}).Decode(&symbol); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	for _, image := range symbol.Images {
		path := filepath.Join("CDN", image.PhotoID.Hex()+".png")
		os.Remove(path)
	}

	if _, err := db.Collection("symbols").DeleteOne(context.Background(), bson.M{
		"_id": symbol.ID,
	}); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "done"})
}

func SetDeleteDeliveryMethodl(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}
	var request struct {
		ID string `json:"delivery_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	deliveryID, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	var delivery models.DeliveryMethods
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}

	if err := db.Collection("delivery_methods").FindOne(context.Background(), bson.M{
		"_id": deliveryID,
	}).Decode(&delivery); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	if _, err := db.Collection("delivery_methods").DeleteOne(context.Background(), bson.M{
		"_id": delivery.ID,
	}); err != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "done"})
}

func SetOrders(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}

	var request models.RequestSetTrade
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	valid := request.Validate(c)
	if !valid {
		return
	}

	payload := map[string]interface{}{
		"comment":   *request.Comment,
		"symbol":    request.SymbolName,
		"type":      request.Operation,
		"volume":    request.Volumn,
		"deviation": *request.Deviation,
		"sl":        *request.StopLoss,
		"tp":        *request.TakeProfit,
		"stoplimit": *request.Stoplimit,
	}
	result, valid := utils.PostRequest(payload, "send_order")

	if !valid || result["status"] == false {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": result["data"], "data": result["data"]})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "done"})
}

func SetCancelOrder(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}

	var request models.RequestSetCancelTrade
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	payload := map[string]interface{}{
		"ticket_id": request.Ticket,
	}
	result, valid := utils.PostRequest(payload, "cancel_order")

	if !valid || result["status"] == false {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": result["data"], "data": result["data"]})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "done"})
}

func SendNotification(c *gin.Context) {
	if !models.AllowedAction(c, models.ActionWrite) {
		return
	}

	var request struct {
		All     bool     `json:"all"`
		UserIDS []string `json:"users"`
		Text    string   `json:"text"`
		Title   string   `json:"title"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}

	if request.Title == "" || len(request.Title) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid title"})
		return
	}

	if request.Text == "" || len(request.Text) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid text"})
		return
	}
	var users []models.User
	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}
	cursor, err := db.Collection("users").Find(context.Background(), bson.M{})
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

	if request.All {
		for _, u := range users {
			if u.FcmToken != "" {
				models.Notification(c, u.ID, request.Title, request.Text)
			}
		}
	} else {
		for _, u := range users {
			for _, ru := range request.UserIDS {
				userID, err := primitive.ObjectIDFromHex(ru)
				if err != nil {
					continue
				}
				if u.ID == userID {
					if u.FcmToken != "" {
						models.Notification(c, u.ID, request.Title, request.Text)
					}
				}
			}
		}
	}
}
