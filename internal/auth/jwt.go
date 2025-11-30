package auth

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int64, secret string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "blog-api", 
			Subject:   "user_authentication",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
	log.Printf("🔍 ValidateToken called")
	log.Printf("🔍 Token length: %d", len(tokenString))
	log.Printf("🔍 Secret length: %d", len(secret))
	
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		log.Printf("🔍 Parsing token with method: %v", token.Method)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("❌ Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		log.Printf("❌ Parse error: %v", err)
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token is expired")
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, errors.New("token is not valid yet")
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, errors.New("invalid token signature")
		}
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	if !token.Valid {
		log.Printf("❌ Token is not valid")
		return nil, errors.New("invalid token")
	}
	
	log.Printf("✅ Token validated successfully for user_id: %d", claims.UserID)
	return claims, nil
}
