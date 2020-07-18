package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "start server on port 8080")
	})

	fmt.Println(http.ListenAndServe(":8080", nil))
}
