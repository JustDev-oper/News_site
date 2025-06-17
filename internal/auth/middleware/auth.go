package middleware

import (
	"context"
	"net/http"

	"News_site/internal/auth/jwt"
)

// UserData структура для хранения данных пользователя
type UserData struct {
	ID    uint
	Email string
}

// AuthMiddleware проверяет JWT токен в cookie
func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем токен из cookie
			cookie, err := r.Cookie("token")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			// Валидируем токен
			claims, err := jwt.ValidateToken(cookie.Value, secretKey)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			// Создаем новый контекст с данными пользователя
			ctx := r.Context()
			ctx = context.WithValue(ctx, "user", &UserData{
				ID:    claims.UserID,
				Email: claims.Email,
			})

			// Вызываем следующий обработчик с обновленным контекстом
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext получает данные пользователя из контекста
func GetUserFromContext(ctx context.Context) *UserData {
	if user, ok := ctx.Value("user").(*UserData); ok {
		return user
	}
	return nil
}
