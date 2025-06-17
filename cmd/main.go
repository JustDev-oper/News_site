package main

import (
	"News_site/internal/config"
	"News_site/internal/core"
	"News_site/internal/db"
	"News_site/internal/handlers/article"
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

	r := mux.NewRouter()
	r.HandleFunc("/", articleHandler.GetAllArticles).Methods("GET")
	r.HandleFunc("/create", articleHandler.CreateArticle).Methods("GET")
	r.HandleFunc("/save_article", articleHandler.SaveArticle).Methods("POST")
	r.HandleFunc("/post/{id:[0-9]+}", articleHandler.GetArticleByID).Methods("GET")

	log.Printf("Server starting on port %s...", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatal(err)
	}
}
