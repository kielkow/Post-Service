package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kielkow/Post-Service/env"
	"github.com/kielkow/Post-Service/modules/author"
	"github.com/kielkow/Post-Service/modules/post"
	"github.com/kielkow/Post-Service/modules/session"
	"github.com/kielkow/Post-Service/shared/database"
)

const apiBasePath = "/api"

func main() {
	env.SetEnv()

	database.SetupDatabase()

	session.SetupRoutes(apiBasePath)
	author.SetupRoutes(apiBasePath)
	post.SetupRoutes(apiBasePath)

	log.Fatal(http.ListenAndServe(":3333", nil))
}
