package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Post struct {
	Id    int
	Title string
	Body  string
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		post := Post{Id: 1, Title: "Unamed Post", Body: "No content"}

		if title := r.FormValue("title"); title != "" {
			post.Title = title
		}

		if body := r.FormValue("body"); body != "" {
			post.Body = body
		}

		t := template.Must(template.ParseFiles("templates/index.html"))

		if err := t.ExecuteTemplate(w, "index.html", post); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println(http.ListenAndServe(":8080", nil))
}
