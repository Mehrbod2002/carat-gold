package controlers

import (
	"carat-gold/models"
	"carat-gold/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func LoginOneTimeLoginStep1(c *gin.Context) {
	var loginData models.LoginDataStep1
	if err := c.ShouldBindJSON(&loginData); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}
	err := loginData.Validate(c)
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
		{"phone": loginData.Phone},
	}}).Decode(&user)
	if exist != nil {
		if exist == mongo.ErrNoDocuments {
			c.JSON(400, gin.H{
				"success": false,
				"message": "invalid phone",
				"data":    "invalid_phone",
			})
			return
		}
		log.Println(exist)
		utils.InternalError(c)
		return
	}

	if user.Freeze {
		c.JSON(401, gin.H{
			"success": false,
			"message": "your account freezed by admin",
			"data":    "freezed_account",
		})
		return
	}
	if !user.PhoneVerified {
		c.JSON(406, gin.H{
			"success": false,
			"message": "phone isn't validated",
			"data":    "unverified_user",
		})
		return
	}
	// if user.ReTryOtp == 5 && time.Since(user.OtpValid) < time.Hour { // Test
	// 	c.JSON(406, gin.H{
	// 		"success": false,
	// 		"message": "otp freezed for 1 hour",
	// 		"data":    "otp_freezed_for_1_hour",
	// 	})
	// 	return
	// }
	otp_code := 12345
	// otp_code := utils.GenerateRandomCode()
	// sent, errMessage := models.Sendotp(user.PhoneNumber, fmt.Sprint(otp_code))
	// if !sent {
	// 	log.Println("otp : ", errMessage)
	// 	c.JSON(500, gin.H{
	// 		"success": false,
	// 		"message": errMessage,
	// 		"data":    "failed_otp",
	// 	})
	// 	return
	// }
	if user.ReTryOtp == 5 && time.Since(user.OtpValid) > time.Hour {
		user.ReTryOtp = 0
	}
	_, errs := db.Collection("users").UpdateOne(context.Background(),
		bson.M{"_id": user.ID}, bson.M{
			"$set": bson.M{
				"otp_code":  &otp_code,
				"otp_valid": time.Now().UTC(),
				"retry_otp": user.ReTryOtp + 1,
			},
		})
	if errs != nil {
		log.Println(err)
		utils.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "logged in",
		"data":    "otp_sent",
	})
}

func LoginOneTimeLoginStep2(c *gin.Context) {
	session := sessions.Default(c)
	var loginData models.LoginDataStep2
	if err := c.ShouldBindJSON(&loginData); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}
	err := loginData.Validate(c)
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
		{"phone": loginData.Phone},
	}}).Decode(&user)
	if exist != nil {
		if exist == mongo.ErrNoDocuments {
			c.JSON(400, gin.H{
				"success": false,
				"message": "invalid phone",
				"data":    "invalid_phone",
			})
			return
		}
		log.Println(exist)
		utils.InternalError(c)
		return
	}
	if user.Freeze {
		c.JSON(401, gin.H{
			"success": false,
			"message": "your account freezed by admin",
			"data":    "freezed_account",
		})
		return
	}
	if !user.PhoneVerified {
		c.JSON(406, gin.H{
			"success": false,
			"message": "phone isn't validated",
			"data":    "unverified_user",
		})
		return
	}

	if time.Since(user.OtpValid) > time.Minute*2 { // Test
		c.JSON(400, gin.H{
			"success": false,
			"message": "request for otp first",
			"data":    "otp_expired",
		})
		return
	}
	if *user.OtpCode != *loginData.Otp {
		c.JSON(406, gin.H{
			"success": false,
			"message": "invalid otp",
			"data":    "invalid_otp",
		})
	}
	_, errs := db.Collection("users").UpdateOne(context.Background(), bson.M{
		"_id": user.ID,
	}, bson.M{"$set": bson.M{
		"otp_code":  nil,
		"otp_valid": time.Now().UTC(),
	}})
	if errs != nil {
		log.Println(errs)
		utils.InternalError(c)
		return
	}
	token, errs := user.GenerateToken()
	if errs != nil {
		log.Println(errs)
		utils.InternalError(c)
		return
	}

	refreshToken, errs := user.GenerateToken()
	if errs != nil {
		log.Println(errs)
		utils.InternalError(c)
		return
	}
	_, errs = db.Collection("users").UpdateOne(context.Background(), bson.M{
		"_id": user.ID,
	}, bson.M{"$set": bson.M{
		"refresh_token": refreshToken,
	}})
	if errs != nil {
		log.Println(errs)
		utils.InternalError(c)
		return
	}

	session.Set("token", token)
	errs = session.Save()
	if errs != nil {
		log.Println(errs)
		utils.InternalError(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "logged in",
		"data": map[string]string{
			"token":         token,
			"refresh_token": refreshToken,
		},
	})
}

func SendOTP(c *gin.Context) {
	var sendOTPData models.SendOTP
	if err := c.ShouldBindJSON(&sendOTPData); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}
	err := sendOTPData.Validate(c)
	if !err {
		log.Println(err)
		return
	}
	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		log.Println(DBerr)
		return
	}

	var existingUser models.User
	exist := db.Collection("users").FindOne(context.Background(), bson.M{"$or": []bson.M{
		{"phone": sendOTPData.PhoneNumber},
	}}).Decode(&existingUser)
	if exist != nil {
		if exist == mongo.ErrNoDocuments {
			otp_code := 12345 // utils.GenerateRandomCode()
			// sent, errMessage := models.Sendotp(sendOTPData.PhoneNumber, fmt.Sprint(otp_code))
			// if !sent {
			// 	log.Println("otp : ", errMessage)
			// 	c.JSON(500, gin.H{
			// 		"success": false,
			// 		"message": errMessage,
			// 		"data":    "failed_otp",
			// 	})
			// 	return
			// }
			var user models.User
			user.Email = ""
			user.PhoneNumber = sendOTPData.PhoneNumber
			user.OtpCode = &otp_code
			user.ReTryOtp = 0
			user.OtpValid = time.Now().UTC()
			user.CreatedAt = time.Now().UTC()
			user.PhoneVerified = false
			_, err := db.Collection("users").InsertOne(context.Background(), user)
			if err != nil {
				log.Println(err)
				utils.InternalError(c)
				return
			}
			c.JSON(200, gin.H{
				"success": true,
				"message": "done",
				"data":    "otp_sent",
			})
			return
		}
		log.Println(exist)
		utils.InternalError(c)
		return
	}
	// if !existingUser.PhoneVerified {
	allowToSend := time.Since(existingUser.OtpValid) > time.Minute*2
	if allowToSend {
		// if existingUser.ReTryOtp == 5 && time.Since(existingUser.OtpValid) < time.Hour { // Test
		// 	c.JSON(406, gin.H{
		// 		"success": false,
		// 		"message": "otp freezed for 1 hour",
		// 		"data":    "otp_freezed_for_1_hour",
		// 	})
		// 	return
		// }
		otp_code := 12345 // utils.GenerateRandomCode()
		// sent, errMessage := models.Sendotp(sendOTPData.PhoneNumber, fmt.Sprint(otp_code))
		// if !sent {
		// 	log.Println("otp : ", errMessage)
		// 	c.JSON(500, gin.H{
		// 		"success": false,
		// 		"message": errMessage,
		// 		"data":    "failed_otp",
		// 	})
		// 	return
		// }
		if existingUser.ReTryOtp == 5 && time.Since(existingUser.OtpValid) > time.Hour {
			existingUser.ReTryOtp = 0
		}
		_, err := db.Collection("users").UpdateOne(context.Background(),
			bson.M{"_id": existingUser.ID}, bson.M{
				"$set": bson.M{
					"otp_code":  &otp_code,
					"otp_valid": time.Now().UTC(),
					"retry_otp": existingUser.ReTryOtp + 1,
				},
			})
		if err != nil {
			log.Println(err)
			utils.InternalError(c)
			return
		}
		c.JSON(200, gin.H{
			"success": true,
			"message": "done",
			"data":    "otp_sent",
		})
		return
	} else {
		c.JSON(406, gin.H{
			"success": false,
			"message": "2 minutes should pass to send sms",
			"data":    "not_allowed_to_send_sms",
		})
	}
}

func Register(c *gin.Context) {
	session := sessions.Default(c)
	var registerData models.RegisterRequest
	if err := c.ShouldBindJSON(&registerData); err != nil {
		log.Println(err)
		utils.BadBinding(c)
		return
	}
	err := registerData.Validate(c)
	if !err {
		log.Println(err)
		return
	}
	db, DBerr := utils.GetDB(c)
	if DBerr != nil {
		log.Println(DBerr)
		return
	}
	var existingUser models.User
	var newRegister bool = false
	exist := db.Collection("users").
		FindOne(context.Background(), bson.M{
			"phone": registerData.PhoneNumber,
		}).Decode(&existingUser)
	if exist != nil {
		if exist == mongo.ErrNoDocuments {
			c.JSON(400, gin.H{
				"success": false,
				"message": "request for otp first",
				"data":    "invalid_otp",
			})
			return
		}
		log.Println(exist)
		utils.InternalError(c)
		return
	}
	if existingUser.PhoneNumber != registerData.PhoneNumber {
		c.JSON(400, gin.H{
			"success": false,
			"message": "request for otp first",
			"data":    "invalid_otp",
		})
		return
	}
	if existingUser.PhoneVerified {
		newRegister = true
	}
	if time.Since(existingUser.OtpValid) > time.Minute*5 {
		c.JSON(400, gin.H{
			"success": false,
			"message": utils.Cap("otp expired"),
			"data":    "otp_expired",
		})
		return
	}
	if *existingUser.OtpCode != *registerData.OtpCode {
		c.JSON(400, gin.H{
			"success": false,
			"message": "invalid otp code",
			"data":    "invalid_otp",
		})
		return
	}

	if newRegister {
		newUser := models.User{
			Email:         "",
			PhoneNumber:   registerData.PhoneNumber,
			CreatedAt:     time.Now(),
			Currency:      "USD",
			PhoneVerified: true,
			UserVerified:  true,
			StatusString:  models.ApprovedStatus,
			OtpCode:       nil,
			Reason:        "",
		}
		_, errs := db.Collection("users").UpdateOne(context.Background(), bson.M{
			"phone": registerData.PhoneNumber,
		}, bson.M{"$set": newUser})
		if errs != nil {
			log.Println(errs)
			utils.InternalError(c)
			return
		}
		errs = db.Collection("users").
			FindOne(context.Background(), bson.M{
				"phone": registerData.PhoneNumber,
			}).Decode(&existingUser)
		if errs != nil {
			log.Println(errs)
			utils.InternalError(c)
			return
		}
		newUser.ID = existingUser.ID

		fmt.Println("Create token", newUser.PhoneNumber)
		token, er := newUser.GenerateToken()
		if er != nil {
			log.Println(errs)
			utils.InternalError(c)
			return
		}
		session.Set("token", token)
		save_err := session.Save()
		if save_err != nil {
			log.Println(save_err)
			utils.InternalError(c)
			return
		}

		refreshToken, errs := newUser.GenerateToken()
		if errs != nil {
			log.Println(errs)
			utils.InternalError(c)
			return
		}
		_, errs = db.Collection("users").UpdateOne(context.Background(), bson.M{
			"_id": newUser.ID,
		}, bson.M{"$set": bson.M{
			"refresh_token": refreshToken,
		}})
		if errs != nil {
			log.Println(errs)
			utils.InternalError(c)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "registered",
			"data": map[string]string{
				"token":         token,
				"refresh_token": refreshToken,
			},
		})
		return
	} else {
		_, errs := db.Collection("users").UpdateOne(context.Background(), bson.M{
			"_id": existingUser.ID,
		}, bson.M{"$set": bson.M{
			"otp_code":  nil,
			"otp_valid": time.Now().UTC(),
		}})
		if errs != nil {
			log.Println(errs)
			utils.InternalError(c)
			return
		}
		token, errs := existingUser.GenerateToken()
		if errs != nil {
			log.Println(errs)
			utils.InternalError(c)
			return
		}

		refreshToken, errs := existingUser.GenerateToken()
		if errs != nil {
			log.Println(errs)
			utils.InternalError(c)
			return
		}
		_, errs = db.Collection("users").UpdateOne(context.Background(), bson.M{
			"_id": existingUser.ID,
		}, bson.M{"$set": bson.M{
			"refresh_token": refreshToken,
		}})
		if errs != nil {
			log.Println(errs)
			utils.InternalError(c)
			return
		}

		session.Set("token", token)
		errs = session.Save()
		if errs != nil {
			log.Println(errs)
			utils.InternalError(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "logged in",
			"data": map[string]string{
				"token":         token,
				"refresh_token": refreshToken,
			},
		})
		return
	}
}

func ValidateSession(c *gin.Context) {
	session := sessions.Default(c)
	token := session.Get("token")
	tokenString := c.GetHeader("Authorization")
	if token == nil && tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false})
		return
	}
	if token == nil {
		token = tokenString
	}

	jwtSecret := os.Getenv("SESSION_SECRET")
	if jwtSecret == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false})
		return
	}

	parsedToken, err := jwt.Parse(token.(string), func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !parsedToken.Valid {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"success": false})
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if ok && claims["_id"] != nil {
		userID, ok := claims["_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false})
			return
		}
		if _, err := primitive.ObjectIDFromHex(userID); err == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"success": true})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"success": false})
}
