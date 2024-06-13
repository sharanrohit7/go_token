// models/log.go
package models

import (
	"context"
	"log"
	"net/http"
	"time"

	db "github.com/sharanrohit7/gotoken/config"
	"github.com/sharanrohit7/gotoken/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// Log represents the log data structure
type Log struct {
	Message         string    `bson:"message"`
	Timestamp       time.Time `bson:"timestamp"`
	ElapsedTime     float64   `bson:"elapsed_time"`     // Elapsed time for request processing
	RequestDetails  Request   `bson:"request_details"`  // Request details
	ResponseDetails Response  `bson:"response_details"` // Response details
}

// Request represents the details of an HTTP request
type Request struct {
	Method    string                 `bson:"method"`
	URI       string                 `bson:"uri"`
	Headers   map[string]interface{} `bson:"headers"`
	IPAddress string                 `bson:"ip_address"`
	// Add other relevant request details
}

// Response represents the details of an HTTP response
type Response struct {
	StatusCode int `bson:"status_code"`
	// Add other relevant response details
}

// InsertLog inserts a log entry into the "log" collection
func InsertLog(c *http.Request, message string, elapsedTime float64, responseDetails Response) error {
	collection := db.DB.Collection("log")

	requestDetails := Request{
		Method:    c.Method,
		URI:       c.RequestURI,
		Headers:   ExtractHeaders(c),
		IPAddress: utils.GetClientIPAddress(c),
	}

	logEntry := Log{
		Message:         message,
		Timestamp:       time.Now(),
		ElapsedTime:     elapsedTime,
		RequestDetails:  requestDetails,
		ResponseDetails: responseDetails,
	}

	_, err := collection.InsertOne(context.Background(), logEntry)
	if err != nil {
		log.Printf("Error inserting log entry: %v", err)
		return err
	}

	return nil
}

func ExtractHeaders(c *http.Request) map[string]interface{} {
	headers := make(map[string]interface{})

	for key, values := range c.Header {
		// If there is only one value, use it directly
		if len(values) == 1 {
			headers[key] = values[0]
		} else {
			// If there are multiple values, store them as a slice
			headers[key] = values
		}
	}

	return headers
}

// GetLogs retrieves all log entries from the "log" collection
func GetLogs() ([]Log, error) {
	collection := db.DB.Collection("log")

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Printf("Error retrieving logs: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var logs []Log
	err = cursor.All(context.Background(), &logs)
	if err != nil {
		log.Printf("Error decoding logs: %v", err)
		return nil, err
	}

	return logs, nil
}
