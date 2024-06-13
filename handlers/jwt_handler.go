package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sharanrohit7/gotoken/utils"
)

func VerifyHandler(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token not provided"})
		return
	}

	claims, err := utils.VerifyJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token verified",
		"claims":  claims,
	})
}
