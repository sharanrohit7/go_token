package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sharanrohit7/gotoken/handlers"
	"github.com/sharanrohit7/gotoken/middleware"
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
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))

	// Logger Middleware
	router.Use(middleware.LoggingMiddleware())

	// Request Logging
	router.Use(func(c *gin.Context) {
		log.Printf("Received %s request for %s from %s", c.Request.Method, c.Request.URL, c.Request.RemoteAddr)
		c.Next()
		log.Printf("Sent %d response for %s to %s", c.Writer.Status(), c.Request.URL, c.Request.RemoteAddr)
	})

	// Welcome route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to your Golang microservice!"})
	})

	// Routes
	router.POST("/verify", handlers.VerifyHandler)

	router.POST("/signToken", func(c *gin.Context) {
		id := c.GetHeader("id")
		roleID := c.GetHeader("role_id")

		if id == "" || roleID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id or role_id in the request headers"})
			return
		}

		// Check cache first
		if token, found := GetCachedToken(id, roleID); found {
			c.JSON(http.StatusOK, gin.H{"token": token})
			return
		}

		// Generate new token
		claims := map[string]interface{}{
			"id":      id,
			"role_id": roleID,
			"exp":     time.Now().Add(time.Hour * 3).Unix(),
		}
		token, err := utils.SignJWT(claims)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign JWT token"})
			return
		}

		// Cache token
		CacheToken(id, roleID, token)
		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	router.POST("/extractData", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token not provided"})
			return
		}

		// Check cache
		if cachedClaims, found := checkTokenCache(token); found {
			log.Println("Cache hit: returning cached token claims")
			c.JSON(http.StatusOK, gin.H{"id": cachedClaims.ID, "role_id": cachedClaims.RoleID})
			return
		}

		// Extract claims
		id, roleID, err := utils.ExtractClaims(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token or failed to extract claims"})
			return
		}

		// Cache the claims
		cacheTokenClaims(token, id, roleID, 15*time.Minute)
		c.JSON(http.StatusOK, gin.H{"id": id, "role_id": roleID})
	})

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hello"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Start cache cleanup in a goroutine
	startCacheCleanup(10 * time.Minute)

	// Graceful shutdown setup
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Channel to listen for system interrupts
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		log.Printf("Server started on :%s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	// Wait for an interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// Gracefully shut down the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited gracefully")
}

var tokenCache sync.Map

func GetCachedToken(id, roleID string) (string, bool) {
	key := id + roleID
	if token, ok := tokenCache.Load(key); ok {
		return token.(string), true
	}
	return "", false
}

func CacheToken(id, roleID, token string) {
	key := id + roleID
	tokenCache.Store(key, token)
}

type TokenClaims struct {
	ID         string
	RoleID     string
	Expiration time.Time
}

func checkTokenCache(token string) (TokenClaims, bool) {
	if entry, ok := tokenCache.Load(token); ok {
		if claims, valid := entry.(TokenClaims); valid && time.Now().Before(claims.Expiration) {
			return claims, true
		}
		// Remove expired tokens
		tokenCache.Delete(token)
	}
	return TokenClaims{}, false
}

func cacheTokenClaims(token, id, roleID string, ttl time.Duration) {
	tokenCache.Store(token, TokenClaims{
		ID:         id,
		RoleID:     roleID,
		Expiration: time.Now().Add(ttl),
	})
}

func startCacheCleanup(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			tokenCache.Range(func(key, value interface{}) bool {
				claims, ok := value.(TokenClaims)
				if !ok {
					return true // Skip invalid entries
				}
				if time.Now().After(claims.Expiration) {
					tokenCache.Delete(key)
				}
				return true
			})
		}
	}()
}
