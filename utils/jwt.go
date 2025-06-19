package utils

import (
	"fmt"
	"time"

	"backend/config"
	"backend/models"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID primitive.ObjectID, email string) (string, error) {
	expirationTime, err := time.ParseDuration(config.AppConfig.JWTAccessExpiry)
	if err != nil {
		expirationTime = 15 * time.Minute // fallback
	}

	claims := &JWTClaims{
		UserID: userID.Hex(),
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTAccessSecret))
}

func GenerateRefreshToken() string {
	return GenerateUUID()
}

func ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AppConfig.JWTAccessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func GenerateTokenPair(userID primitive.ObjectID, email string) (*models.TokenPair, error) {
	accessToken, err := GenerateAccessToken(userID, email)
	if err != nil {
		return nil, err
	}

	refreshToken := GenerateRefreshToken()

	expirationTime, err := time.ParseDuration(config.AppConfig.JWTAccessExpiry)
	if err != nil {
		expirationTime = 15 * time.Minute
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(expirationTime.Seconds()),
	}, nil
}
