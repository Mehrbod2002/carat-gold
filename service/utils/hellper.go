package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
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
	return rander.Intn(999999-100000+1) + 100000
}

func Sendotp(mobileNumber string, otp string) (bool, *string) {
	// otpSecret := os.Getenv("OTP_SECRET")
	// template := os.Getenv("OTP_TEMPLATE_ID")
	// api := kavenegar.New(otpSecret)
	// receptor := mobileNumber
	// params := &kavenegar.VerifyLookupParam{}
	// if _, err := api.Verify.Lookup(receptor, template, otp, params); err != nil {
	// 	errMsg := err.Error()
	// 	if strings.Contains(errMsg, ":") {
	// 		message := strings.TrimSpace(strings.Split(errMsg, ":")[1])
	// 		return false, &message
	// 	}
	// 	return false, nil
	// } else {
	return true, nil
	// }
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
