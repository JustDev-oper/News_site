package utils

import (
	"News_site/internal/auth/middleware"
	"html/template"
	"net/http"
)

// RenderTemplate вспомогательная функция для рендеринга шаблонов
func RenderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}) {
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