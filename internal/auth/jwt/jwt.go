package jwt

import (
	"errors"
	"time"

	"News_site/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("недействительный токен")
	ErrExpiredToken = errors.New("срок действия токена истек")
)

// Claims представляет структуру JWT токена

// GenerateToken создает новый JWT токен для пользователя
func GenerateToken(user *models.User, secretKey string, duration time.Duration) (string, error) {
	claims := models.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ValidateToken проверяет JWT токен и возвращает claims
func ValidateToken(tokenString string, secretKey string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ExtractUserID извлекает ID пользователя из токена
func ExtractUserID(tokenString string, secretKey string) (uint, error) {
	claims, err := ValidateToken(tokenString, secretKey)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
