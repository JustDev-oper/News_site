package auth

import (
	"News_site/internal/auth/jwt"
	"News_site/internal/services/user"
	"News_site/internal/utils"
	"net/http"
	"time"
)

type Handler struct {
	service   user.Service
	secretKey string
}

func NewHandler(userService user.Service, secretKey string) *Handler {
	return &Handler{
		service:   userService,
		secretKey: secretKey,
	}
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RenderTemplate(w, r, "login", nil)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	userByEmail, err := handler.service.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if !handler.service.ValidatePassword(userByEmail, password) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := jwt.GenerateToken(userByEmail, handler.secretKey, 24*time.Hour)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   86400, // 24 часа
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RenderTemplate(w, r, "register", nil)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	username := r.FormValue("username")

	if err := handler.service.CreateUser(email, password, username); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (handler *Handler) Logout(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1, // Устанавливаем отрицательное значение для немедленного удаления
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
