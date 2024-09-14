package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sharanrohit7/gotoken/handlers"

	// "github.com/sharanrohit7/gotoken/models"

	// "github.com/sharanrohit7/gotoken/config"
	"github.com/sharanrohit7/gotoken/utils"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	router := gin.Default()
	// Configure CORS middleware
	// config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"*"} // Add your domain here
	// router.Use(cors.New(config))

	// Configure logger middleware

	// log.Printf("MongoDB Connection Details:\n")
	// log.Printf(" - URI: %s\n", os.Getenv("MONGODB_URI"))
	// log.Printf(" - Database: %s\n", os.Getenv("MONGODB_DB"))
	// db.init()

	// router.Use(middleware.LoggingMiddleware())

	router.Use(func(c *gin.Context) {
		// Log request
		log.Printf("Received %s request for %s from %s", c.Request.Method, c.Request.URL, c.Request.RemoteAddr)

		// Continue processing
		c.Next()

		// Log response
		log.Printf("Sent %d response for %s to %s", c.Writer.Status(), c.Request.URL, c.Request.RemoteAddr)
	})
	// Welcome route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to your Golang microservice!"})
	})

	// Verify JWT route
	router.POST("/verify", handlers.VerifyHandler)

	// router.POST("/signToken", func(c *gin.Context) {
	// 	// Check if id and role_id are available in the request body
	// 	id, idExists := c.GetPostForm("id")
	// 	roleID, roleIDExists := c.GetPostForm("role_id")

	// 	if !idExists || !roleIDExists || id == "" || roleID == "" {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id or role_id in the request body"})
	// 		return
	// 	}

	// 	// Create a map with claims
	// 	claims := map[string]interface{}{
	// 		"id":      id,
	// 		"role_id": roleID,
	// 		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (e.g., 24 hours from now)
	// 	}

	// 	// Sign the JWT token
	// 	token, err := utils.SignJWT(claims)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign JWT token"})
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, gin.H{"token": token})
	// })
	router.POST("/signToken", func(c *gin.Context) {
		// Extract id and role_id from headers
		id := c.GetHeader("id")
		roleID := c.GetHeader("role_id")

		if id == "" || roleID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id or role_id in the request headers"})
			return
		}

		// Create a map with claims
		claims := map[string]interface{}{
			"id":      id,
			"role_id": roleID,
			"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (e.g., 24 hours from now)
		}

		// Sign the JWT token
		token, err := utils.SignJWT(claims)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign JWT token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	router.GET("/extractData", func(c *gin.Context) {
		// Get the token from the Authorization header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token not provided"})
			return
		}

		// Extract ID and role_id using the new function
		id, roleID, err := utils.ExtractClaims(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token or failed to extract claims"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id, "role_id": roleID})
	})
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hello"})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server started on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
