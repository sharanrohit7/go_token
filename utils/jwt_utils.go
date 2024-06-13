package utils

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte(os.Getenv("JWT_SECRET"))

func VerifyJWT(tokenString string) (jwt.MapClaims, error) {
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Invalid token")
}

func SignJWT(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ExtractClaims(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, idExists := claims["id"].(string)
		roleID, roleIDExists := claims["role_id"].(string)

		if !idExists || !roleIDExists {
			log.Printf("Claims: %v", claims)
			log.Println("Failed to extract ID or role_id from claims")
			return "", "", errors.New("Failed to extract ID or role_id from claims")
		}

		return id, roleID, nil
	}

	log.Printf("Claims: %v", token.Claims)
	log.Println("Invalid token claims")
	return "", "", errors.New("Invalid token claims")
}

func GetClientIPAddress(req *http.Request) string {
	// Check if the request went through a proxy or load balancer
	ipAddress := req.Header.Get("X-Real-IP")
	if ipAddress == "" {
		// If X-Real-IP header is not present, try X-Forwarded-For
		ipAddress = req.Header.Get("X-Forwarded-For")
	}

	// If neither header is present, get the IP from the request
	if ipAddress == "" {
		ipAddress = req.RemoteAddr
	}

	return ipAddress
}
