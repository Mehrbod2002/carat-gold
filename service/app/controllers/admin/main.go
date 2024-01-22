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
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

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

func AdminLogin(c *gin.Context) {
	session := sessions.Default(c)
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	if loginData.Email == "" || loginData.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "email or password is empty",
		})
		return
	}

	db, err := utils.GetDB(c)
	if err != nil {
		log.Println(err)
		return
	}

	var user models.User
	if err := db.Collection("admin").FindOne(context.Background(), bson.M{
		"email": loginData.Email,
	}).Decode(&user); err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			c.JSON(401, gin.H{"success": false, "message": "Invalid email or password"})
		} else {
			utils.InternalError(c)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		log.Println(err)
		c.JSON(401, gin.H{"success": false, "message": "Invalid email or password"})
		return
	}

	token, err := user.GenerateToken()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate token",
		})
		return
	}

	session.Set("token_admins", token)
	err = session.Save()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate token",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged in successful",
		"token":   token,
	})
}
