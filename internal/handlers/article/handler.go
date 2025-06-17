package article

import (
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

	temp, err := template.ParseFiles("web/templates/showPost.html", "web/templates/header.html", "web/templates/footer.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := temp.ExecuteTemplate(w, "showPost", articleById); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (handler *Handler) GetAllArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := handler.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	temp, err := template.ParseFiles("web/templates/index.html", "web/templates/header.html", "web/templates/footer.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := temp.ExecuteTemplate(w, "index", articles); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (handler *Handler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("web/templates/create.html", "web/templates/header.html", "web/templates/footer.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err := temp.ExecuteTemplate(w, "create", nil); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (handler *Handler) SaveArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	anons := r.FormValue("anons")
	fullText := r.FormValue("full_text")

	if err := handler.service.Create(title, anons, fullText); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
