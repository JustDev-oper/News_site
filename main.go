package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

type Article struct {
	Id       uint16
	Title    string
	Anons    string
	FullText string
}

func index(w http.ResponseWriter, r *http.Request) {
	var articles []Article

	temp, err := template.ParseFiles("News_site/templates/index.html", "News_site/templates/header.html", "News_site/templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:Ar11042008@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM `articles`")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		articles = append(articles, post)
	}

	temp.ExecuteTemplate(w, "index", articles)
}

func create(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("News_site/templates/create.html", "News_site/templates/header.html", "News_site/templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	temp.ExecuteTemplate(w, "create", nil)
}

func saveArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	fullText := r.FormValue("full_text")

	if title == "" || anons == "" || fullText == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {

		db, err := sql.Open("mysql", "root:Ar11042008@tcp(127.0.0.1:3306)/golang")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		// Установка данных
		insert, err := db.Query("INSERT INTO `articles` (`title`, `anons`, `full_text`) VALUES (?, ?, ?)", title, anons, fullText)
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func showPost(w http.ResponseWriter, r *http.Request) {
	var post Article

	id := mux.Vars(r)["id"]

	db, err := sql.Open("mysql", "root:Ar11042008@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM `articles` WHERE `id`=?", id)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	for res.Next() {
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		post = post
	}

	temp, err := template.ParseFiles("News_site/templates/showPost.html", "News_site/templates/header.html", "News_site/templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	temp.ExecuteTemplate(w, "showPost", post)

}

func handleFunc() {
	router := mux.NewRouter()

	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/create", create).Methods("GET")
	router.HandleFunc("/save_article", saveArticle).Methods("POST")
	router.HandleFunc("/post/{id:[0-9]+}", showPost).Methods("GET")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func main() {

	handleFunc()
}
