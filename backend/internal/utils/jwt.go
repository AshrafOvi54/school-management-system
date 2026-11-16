package utils

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT that expires after 24 hours.
func GenerateToken(userID uint, role string) (string, error) {
	secretStr := os.Getenv("JWT_SECRET")
	if secretStr == "" {
		log.Fatal("FATAL SECURITY ERROR: JWT_SECRET environment variable is missing in .env")
	}
	jwtSecret := []byte(secretStr)

	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
