package article

import (
	"News_site/internal/auth/middleware"
	"News_site/internal/services/article"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	service article.Service
}

func NewHandler(service article.Service) *Handler {
	return &Handler{service: service}
}

// renderTemplate вспомогательная функция для рендеринга шаблонов
func (h *Handler) renderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}) {

	user := middleware.GetUserFromContext(r.Context())

	templateData := struct {
		User *middleware.UserData
		Data interface{}
	}{
		User: user,
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

func (handler *Handler) GetArticleByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 16)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	articleById, err := handler.service.GetByID(uint16(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	handler.renderTemplate(w, r, "showPost", articleById)
}

func (handler *Handler) GetAllArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := handler.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handler.renderTemplate(w, r, "index", articles)
}

func (handler *Handler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	handler.renderTemplate(w, r, "createPost", nil)
}

func (handler *Handler) SaveArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	title := r.FormValue("title")
	anons := r.FormValue("anons")
	fullText := r.FormValue("full_text")

	if err := handler.service.Create(title, anons, fullText, user.ID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (handler *Handler) GetUserArticles(w http.ResponseWriter, r *http.Request) {

	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	articles, err := handler.service.GetByUserID(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handler.renderTemplate(w, r, "userPosts", articles)
}

func (handler *Handler) EditArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 16)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	if r.Method == "GET" {

		article, err := handler.service.GetByID(uint16(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if article.UserID != user.ID {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		handler.renderTemplate(w, r, "editPost", article)
		return
	}

	if r.Method == "POST" {
		title := r.FormValue("title")
		anons := r.FormValue("anons")
		fullText := r.FormValue("full_text")

		if err := handler.service.Update(uint16(id), title, anons, fullText, user.ID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/my-posts", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (handler *Handler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 16)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Проверяем, что статья принадлежит пользователю
	article, err := handler.service.GetByID(uint16(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if article.UserID != user.ID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	if err := handler.service.Delete(uint16(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/my-posts", http.StatusSeeOther)
}
