package main

import (
	"News_site/internal/auth/middleware"
	"News_site/internal/config"
	"News_site/internal/core"
	"News_site/internal/db"
	"News_site/internal/handlers/article"
	"News_site/internal/handlers/auth"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := db.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	container := core.NewContainer(dbConn.GetDB())
	defer container.Close()

	articleHandler := article.NewHandler(container.ArticleService)
	authHandler := auth.NewHandler(container.UserService, cfg.JWTSecret)

	r := mux.NewRouter()

	// Применяем middleware аутентификации ко всем маршрутам
	r.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	r.HandleFunc("/", articleHandler.GetAllArticles).Methods("GET")

	// User
	r.HandleFunc("/login", authHandler.Login).Methods("GET", "POST")
	r.HandleFunc("/register", authHandler.Register).Methods("GET", "POST")
	r.HandleFunc("/logout", authHandler.Logout).Methods("GET")
	r.HandleFunc("/my-posts", articleHandler.GetUserArticles).Methods("GET")

	// Articles
	r.HandleFunc("/post/{id:[0-9]+}", articleHandler.GetArticleByID).Methods("GET")
	r.HandleFunc("/create", articleHandler.CreateArticle).Methods("GET")
	r.HandleFunc("/save_article", articleHandler.SaveArticle).Methods("POST")
	r.HandleFunc("/edit/{id:[0-9]+}", articleHandler.EditArticle).Methods("GET", "POST")
	r.HandleFunc("/delete/{id:[0-9]+}", articleHandler.DeleteArticle).Methods("POST")
	r.HandleFunc("/post/{id:[0-9]+}/like", articleHandler.LikeArticleHandler).Methods("POST")

	log.Printf("Server starting on port %s...", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatal(err)
	}
}
