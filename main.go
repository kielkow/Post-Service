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
	r.HandleFunc("/{id}/view", ViewHandler)
	r.HandleFunc("/addpost/{title}/{body}", AddPost)

	fmt.Println(http.ListenAndServe(":8080", r))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/layout.html", "templates/list.html"))

	if err := t.ExecuteTemplate(w, "layout.html", ListPosts()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	t := template.Must(template.ParseFiles("templates/layout.html", "templates/view.html"))

	if err := t.ExecuteTemplate(w, "layout.html", GetPostById(id)); err != nil {
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

func GetPostById(id string) *Post {
	row := db.QueryRow("select * from posts where id = ?", id)
	post := Post{}
	row.Scan(&post.Id, &post.Title, &post.Body)

	return &post
}

func AddPost(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]
	body := mux.Vars(r)["body"]

	stmt, err := db.Prepare("insert into posts(title, body) values(?, ?)")
	checkErr(err)

	_, err = stmt.Exec(title, body)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
