// middleware.go
package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/sharanrohit7/gotoken/config"
)

// LoggingMiddleware logs all incoming requests and their responses
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Process response
		endTime := time.Now()

		// Extract relevant details from the request and response
		requestDetails := getRequestDetails(c.Request)
		responseDetails := getResponseDetails(c.Writer)

		// Combine request and response details
		message := map[string]interface{}{
			"timestamp":        time.Now(),
			"elapsed_time":     endTime.Sub(startTime).Seconds(),
			"request_details":  requestDetails,
			"response_details": responseDetails,
		}

		// Insert into the "log" collection
		_, err := db.DB.Collection("log").InsertOne(context.Background(), message)
		if err != nil {
			log.Printf("Error logging message: %v", err)
		}
	}
}

// getRequestDetails extracts relevant details from the request
func getRequestDetails(request *http.Request) map[string]interface{} {
	// Extract relevant details from the request (headers, method, etc.)
	// Example: Get headers
	headers := make(map[string]interface{})
	for key, values := range request.Header {
		headers[key] = values
	}

	requestDetails := map[string]interface{}{
		"method":  request.Method,
		"uri":     request.RequestURI,
		"headers": headers,
		// Add other relevant details as needed
	}

	return requestDetails
}

// getResponseDetails extracts relevant details from the response
func getResponseDetails(writer gin.ResponseWriter) map[string]interface{} {
	// Extract relevant details from the response (status code, headers, etc.)
	// Example: Get status code
	statusCode := writer.Status()

	responseDetails := map[string]interface{}{
		"status_code": statusCode,
		// Add other relevant details as needed
	}

	return responseDetails
}
