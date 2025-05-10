package app

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/invinciblewest/gophermart/internal/logger"
	"github.com/invinciblewest/gophermart/internal/model"
	"go.uber.org/zap"
	"time"
)

type AuthUseCase struct {
	secretKey string
}

func NewAuthUseCase(secretKey string) *AuthUseCase {
	return &AuthUseCase{
		secretKey: secretKey,
	}
}

func (as *AuthUseCase) GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(as.secretKey))
}

func (as *AuthUseCase) ParseToken(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(as.secretKey), nil
	})

	if err != nil || !token.Valid {
		logger.Log.Error("failed to parse token", zap.Error(err))
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id not found")
	}

	return int(userIDFloat), nil
}

func (as *AuthUseCase) HashPassword(password string) string {
	hash := hmac.New(sha256.New, []byte(as.secretKey))
	hash.Write([]byte(password))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func (as *AuthUseCase) VerifyPassword(user *model.User, password string) bool {
	hashedPassword := as.HashPassword(password)
	return hmac.Equal([]byte(user.Password), []byte(hashedPassword))
}
