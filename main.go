package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kielkow/Post-Service/author"
	"github.com/kielkow/Post-Service/database"
	"github.com/kielkow/Post-Service/post"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()

	author.SetupRoutes(apiBasePath)
	post.SetupRoutes(apiBasePath)

	log.Fatal(http.ListenAndServe(":3333", nil))
}
