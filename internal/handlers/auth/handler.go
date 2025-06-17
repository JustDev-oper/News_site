package auth

import (
	"News_site/internal/auth/jwt"
	"News_site/internal/auth/middleware"
	"News_site/internal/services/user"
	"html/template"
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

// renderTemplate вспомогательная функция для рендеринга шаблонов
func (h *Handler) renderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}) {
	// Получаем данные пользователя из контекста
	userFromContext := middleware.GetUserFromContext(r.Context())

	templateData := struct {
		User *middleware.UserData
		Data interface{}
	}{
		User: userFromContext,
		Data: data,
	}

	temp, err := template.ParseFiles(
		"web/templates/"+templateName+".html",
		"web/templates/header.html",
		"web/templates/footer.html",
	)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := temp.ExecuteTemplate(w, templateName, templateData); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handler.renderTemplate(w, r, "login", nil)
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
		handler.renderTemplate(w, r, "register", nil)
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
