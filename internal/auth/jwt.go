package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(userID uint, role string, secretKey string) (string, error) {
	// Create a new JWT token with the user ID as a claim and sign it with the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"expiresAt": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expires in 24 hours
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString string, secretKey string) (*jwt.Token, uint, error) {
	// Implementation for validating JWT token goes here
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if userID, ok := claims["user_id"].(float64); ok {
			return token, uint(userID), nil
		}
	}
	return nil, 0, nil
}