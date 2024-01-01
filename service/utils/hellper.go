package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
