package utils

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BadBinding(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": "invalid request parameters",
		"data":    "invalid_parameters",
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
