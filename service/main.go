package main

import (
	"carat-gold/models"
	"carat-gold/routes"
	"carat-gold/utils"
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	executablePath, err := os.UserHomeDir()
	if err != nil {
		log.Println("failed to get executable path")
		return
	}
	executableDir := filepath.Dir(executablePath)
	envFilePath := filepath.Join(executableDir, "admin", "carat-gold", "service", ".env")
	err = godotenv.Load(envFilePath)
	if err != nil {
		log.Println("failed to load env file")
		return
	}

	logFile, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()
	log.SetOutput(os.Stdout)

	DBerr := utils.InitializeDB()
	if DBerr != nil {
		log.Println("failed to connect to mongodb")
		return
	}
	defer utils.CloseDB()

	err = utils.InitializeApp()
	if err != nil {
		log.Println(err)
		return
	}

	db, DBerr := utils.GetDBWSS()
	if DBerr != nil {
		log.Println(DBerr)
		log.Println("failed to connect to mongodb")
		return
	}

	var user models.User
	var adminUsername = os.Getenv("ADMIN_USERNAME")
	var adminPassword = os.Getenv("ADMIN_PASSWORD")
	var dataChannel = make(chan interface{}, 50)
	if exist := db.Collection("users").FindOne(context.Background(), bson.M{
		"email": adminUsername,
	}).Decode(&user); exist != nil && exist == mongo.ErrNoDocuments {
		hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		if errHash != nil {
			log.Println("failed to create admin user ", errHash)
			return
		}
		if _, err := db.Collection("users").InsertOne(context.Background(), bson.M{
			"email":    os.Getenv("ADMIN_USERNAME"),
			"password": hashedPassword,
		}); err != nil {
			log.Println("failed to create admin user ", err)
			return
		}
	}

	routes := routes.SetupRouter(dataChannel)
	runningErr := routes.Run(":3000")
	log.Println("start serving ...")
	if runningErr != nil {
		log.Println(runningErr)
		return
	}
}
