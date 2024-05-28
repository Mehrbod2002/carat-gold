package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BadBinding(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": Cap("invalid request parameters"),
		"data":    "invalid_parameters",
	})
}

func Unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"message": Cap("unauthorized"),
		"data":    "unauthorized",
	})
}

func InternalErrorMsg(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": Cap(message),
		"data":    "internal_error",
	})
}

func InternalError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": Cap("internal server connection"),
		"data":    "internal_error",
	})
}

func AdminError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
		"success": false,
		"message": Cap("Currenctly we don't get new payments"),
		"data":    "internal_error",
	})
}

func Method(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": Cap(message),
		"data":    "invalid_parameters",
	})
}

func Cap(s string) string {
	if len(s) == 0 {
		return s
	}

	firstLetter := string(s[0])
	firstLetter = strings.ToUpper(firstLetter)

	return firstLetter + s[1:]
}

func ValidateID(Id string, c *gin.Context) (primitive.ObjectID, bool) {
	userID, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		BadBinding(c)
		return primitive.ObjectID{}, false
	}
	return userID, true
}

func GenerateRandomCode() int {
	rander := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rander.Intn(90000) + 10000
}

func ValidateAdmin(token string) bool {
	jwtSecret := os.Getenv("SESSION_SECRET")
	if jwtSecret == "" {
		return false
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !parsedToken.Valid {
		return false
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if ok && claims["_id"] != nil {
		userID, ok := claims["_id"].(string)
		if !ok {
			return false
		}
		if _, err := primitive.ObjectIDFromHex(userID); err == nil {
			return true
		}
		return false
	}

	return false
}

func PostRequest(data map[string]interface{}, endPoint string) (map[string]interface{}, bool) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, false
	}

	url := "http://54.163.221.86/" + endPoint
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, false
	}
	defer req.Body.Close()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Secret-Header", "MysticalDragon$7392&WhisperingWinds&SunsetHaven$AuroraBorealis")
	resp, err := client.Do(req)

	if err != nil {
		return nil, false
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, false
	}

	return result, true
}

func GetRequest(endPoint string) (map[string]interface{}, bool) {
	url := "http://54.163.221.86/" + endPoint

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(nil))
	if err != nil {
		return nil, false
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Secret-Header", "MysticalDragon$7392&WhisperingWinds&SunsetHaven$AuroraBorealis")

	resp, err := client.Do(req)
	if err != nil {
		return nil, false
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, false
	}
	return result, true
}

func DerefStringPtr(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

func DerefBoolPtr(ptr *bool) bool {
	if ptr != nil {
		return *ptr
	}
	return false
}

func DerefIntPtr(ptr *int) int {
	if ptr != nil {
		return *ptr
	}
	return 0
}

func UploadPhoto(c *gin.Context, id string, photo string, doc bool) bool {
	svgPath := filepath.Join("CDN", id+".svg")
	if doc {
		svgPath = filepath.Join("CDN", id+".png")
	}

	var photoData string
	var decoded []byte
	if strings.Contains(photo, "base64,") {
		photoData = strings.Split(photo, ",")[1]
	} else {
		decodedData, err := base64.StdEncoding.DecodeString(photo)
		if err != nil {
			InternalError(c)
			return false
		}
		decoded = decodedData
	}

	if len(decoded) == 0 {
		decodedData, err := base64.StdEncoding.DecodeString(photoData)
		if err != nil {
			InternalError(c)
			return false
		}
		decoded = decodedData
	}

	err := os.WriteFile(svgPath, decoded, 0644)
	if err != nil {
		InternalError(c)
		return false
	}

	return true
}

func AutoOrder(price float64) (*uint64, bool) {
	result, valid := GetRequest("get_last_price")

	if !valid {
		return nil, false
	}
	if !result["status"].(bool) {
		return nil, false
	}

	// volumn := result["data"].(float64) / price
	// fmt.Println(volumn)
	payload := map[string]interface{}{
		"comment":   "User Payment Stream",
		"symbol":    "XAUUSD",
		"type":      1,
		"volume":    0.1,
		"deviation": 0,
		"sl":        0,
		"tp":        0,
		"stoplimit": 0,
	}
	result, valid = PostRequest(payload, "send_order")

	status := result["status"].(bool)
	if !valid || !status {
		return nil, false
	}

	dataString, ok := result["data"].(string)
	if !ok {
		return nil, false
	}

	orderID, err := strconv.ParseUint(dataString, 10, 64)
	if err != nil {
		return nil, false
	}
	return &orderID, true
}

func TrimAndLowerCase(data string) string {
	return strings.ToLower(strings.TrimSpace(data))
}
