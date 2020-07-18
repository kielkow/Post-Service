package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

type Post struct {
	Id    int
	Title string
	Body  string
}

var db, err = sql.Open("mysql", "root:root@/go_course?charset=utf8")

func main() {
	// stmt, err := db.Prepare("insert into posts(title, body) values(?, ?)")
	// checkErr(err)

	// _, err = stmt.Exec("My first post", "My first content")
	// checkErr(err)

	r := mux.NewRouter()
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

	r.HandleFunc("/", HomeHandler)

	fmt.Println(http.ListenAndServe(":8080", r))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/index.html"))

	if err := t.ExecuteTemplate(w, "index.html", ListPosts()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ListPosts() []Post {
	rows, err := db.Query("select * from posts")
	checkErr(err)

	items := []Post{}

	for rows.Next() {
		post := Post{}

		rows.Scan(&post.Id, &post.Title, &post.Body)
		items = append(items, post)
	}

	return items
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
