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
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BadBinding(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": "invalid request parameters",
		"data":    "invalid_parameters",
	})
}

func InternalErrorMsg(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": Cap(message),
		"data":    "internal_error",
	})
}

func InternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": "internal server connection",
		"data":    "internal_error",
	})
}

func Method(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
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

func UploadPhoto(c *gin.Context, id string, photo string) bool {
	pngPath := filepath.Join("CDN", id+".png")

	var photoData []byte
	if strings.Contains(photo, "base64,") {
		photoData = []byte(strings.Split(photo, ",")[1])
	} else {
		InternalError(c)
		return false
	}

	decodedData, err := base64.StdEncoding.DecodeString(string(photoData))
	if err != nil {
		InternalError(c)
		return false
	}

	err = os.WriteFile(pngPath, decodedData, 0644)
	if err != nil {
		InternalError(c)
		return false
	}

	// ppmPath := filepath.Join("temp", id+".ppm")
	// cmdConvert := exec.Command("convert", pngPath, ppmPath)
	// err = cmdConvert.Run()
	// if err != nil {
	// 	InternalError(c)
	// 	return false
	// }

	// cmd := exec.Command("potrace", ppmPath, "-s", "-o", svgPath)
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	fmt.Println("Error running potrace:", err)
	// 	fmt.Println("Output:", string(output))
	// 	InternalError(c)
	// 	return false
	// }

	// err = os.Remove(pngPath)
	// if err != nil {
	// 	log.Println(err, 567)
	// 	InternalError(c)
	// 	return false
	// }

	// err = os.Remove(ppmPath)
	// if err != nil {
	// 	log.Println(err, 567)
	// 	InternalError(c)
	// 	return false
	// }

	return true
}

func AutoOrder(c *gin.Context, price float64) bool {
	volume := price - (price / 10)
	payload := map[string]interface{}{
		"comment":   "User Payment Stream",
		"symbol":    "XAUUSD",
		"type":      1,
		"volume":    volume,
		"deviation": 0,
		"sl":        0,
		"tp":        0,
		"stoplimit": 0,
	}
	result, valid := PostRequest(payload, "send_order")

	if !valid || result["status"] == false {
		return false
	}

	return true
}
