package article

import (
	"News_site/internal/auth/middleware"
	"News_site/internal/services/article"
	"News_site/internal/utils"
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

	utils.RenderTemplate(w, r, "showPost", articleById)
}

func (handler *Handler) GetAllArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := handler.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RenderTemplate(w, r, "index", articles)
}

func (handler *Handler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, r, "createPost", nil)
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

	utils.RenderTemplate(w, r, "userPosts", articles)
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

		utils.RenderTemplate(w, r, "editPost", article)
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

func (handler *Handler) LikeArticleHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 16)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	liked, err := handler.service.IsArticleLikedByUser(user.ID, uint16(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if liked {
		err = handler.service.UnlikeArticle(user.ID, uint16(id))
	} else {
		err = handler.service.LikeArticle(user.ID, uint16(id))
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
