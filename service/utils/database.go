package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var Client *mongo.Client

func InitializeDB() error {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbSRV := os.Getenv("DB_SRV")
	dbReplicaSet := os.Getenv("DB_REPLICA")
	dbPassword := os.Getenv("DB_PASSWORD")

	uri := fmt.Sprintf("mongodb%s://%s:%s@%s?%sssl=false", dbSRV, dbUser, dbPassword, dbHost, dbReplicaSet)
	otps := options.Client().ApplyURI(uri)
	var err error
	Client, err = mongo.Connect(context.TODO(), otps)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetDBServer() (*mongo.Database, error) {
	err := Client.Ping(context.Background(), nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return Client.Database(os.Getenv("DB_NAME")), nil
}

func GetDB(c *gin.Context) (*mongo.Database, error) {
	err := Client.Ping(context.Background(), nil)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "internal server connection",
			"data":    "internal_error",
		})
		return nil, err
	}
	return Client.Database(os.Getenv("DB_NAME")), nil
}

func GetDBWSS() (*mongo.Database, error) {
	err := Client.Ping(context.Background(), nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return Client.Database(os.Getenv("DB_NAME")), nil
}

func CloseDB() {
	if Client != nil {
		Client.Disconnect(context.TODO())
	}
}
