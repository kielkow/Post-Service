package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kielkow/Post-Service/database"
	"github.com/kielkow/Post-Service/foo"
	"github.com/kielkow/Post-Service/receipt"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()

	receipt.SetupRoutes(apiBasePath)
	foo.SetupRoutes(apiBasePath)

	log.Fatal(http.ListenAndServe(":3333", nil))
}
