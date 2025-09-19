package auth

import (
	"fmt"
	"time"

	"github.com/akhilbisht798/gocrony/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenrateJWT(userId string, email string, provider string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"email":   email,
		"provider": provider,
		"exp":     time.Now().Add(time.Hour * 720).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwt_secret := config.GetEnv("JWT_SECRET", "superSecretKey")

	return token.SignedString([]byte(jwt_secret))
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	jwt_secret := config.GetEnv("JWT_SECRET", "superSecretKey")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(jwt_secret), nil
	})
}
