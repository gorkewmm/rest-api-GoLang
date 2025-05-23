package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("supersecret")

func GenerateToken(email string, userId int64, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,  // Kullanıcının e-posta adresi
		"userId": userId, // Kullanıcının ID'si
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 1).Unix(), // Token'ın 1 saat sonra geçersiz olacağı zamanı belirtiyor

	})
	return token.SignedString(secretKey)
}

// verify that token
func VerifyToken(token string) (int64, string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return 0, "", errors.New("Could not parse token")
	}

	tokenIsValid := parsedToken.Valid
	if !tokenIsValid {
		return 0, "", errors.New("Invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("Invalid token")
	}

	// email := claims["email"].(string)
	userId := int64(claims["userId"].(float64))
	role := claims["role"].(string)
	return userId, role, nil
}
