package author

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/kielkow/Post-Service/apperror"
	"github.com/kielkow/Post-Service/cors"
)

const authorsBasePath = "authors"

// SetupRoutes function
func SetupRoutes(apiBasePath string) {
	handleAuthors := http.HandlerFunc(authorsHandler)
	handleAuthor := http.HandlerFunc(authorHandler)

	reportHandler := http.HandlerFunc(handleAuthorReport)

	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, authorsBasePath), cors.Middleware(handleAuthors))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, authorsBasePath), cors.Middleware(handleAuthor))

	http.Handle(fmt.Sprintf("%s/%s/reports", apiBasePath, authorsBasePath), cors.Middleware(reportHandler))
}

func authorsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		authorList, err := getAuthorList()

		if err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		authorsJSON, err := json.Marshal(authorList)

		if err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(authorsJSON)

	case http.MethodPost:
		var newAuthor CreateAuthor
		bodyBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			error := apperror.GenerateError(400, err.Error())

			w.WriteHeader(http.StatusBadRequest)
			w.Write(error)
			return
		}

		err = json.Unmarshal(bodyBytes, &newAuthor)

		if err != nil {
			error := apperror.GenerateError(400, err.Error())

			w.WriteHeader(http.StatusBadRequest)
			w.Write(error)
			return
		}

		_, err = insertAuthor(newAuthor)

		if err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return

	case http.MethodOptions:
		return
	}
}

func authorHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "authors/")
	id, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])

	if err != nil {
		error := apperror.GenerateError(500, err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(error)
		return
	}

	author, err := GetAuthor(id)

	if err != nil {
		error := apperror.GenerateError(500, err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(error)
		return
	}

	if author == nil {
		error := apperror.GenerateError(404, "Author not found")

		w.WriteHeader(http.StatusNotFound)
		w.Write(error)
		return
	}

	switch r.Method {
	case http.MethodGet:
		authorJSON, err := json.Marshal(author)

		if err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(authorJSON)

	case http.MethodPut:
		var updatedAuthor Author

		bodyBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			error := apperror.GenerateError(400, err.Error())

			w.WriteHeader(http.StatusBadRequest)
			w.Write(error)
			return
		}

		err = json.Unmarshal(bodyBytes, &updatedAuthor)

		if err != nil {
			error := apperror.GenerateError(400, err.Error())

			w.WriteHeader(http.StatusBadRequest)
			w.Write(error)
			return
		}

		err = updateAuthor(id, updatedAuthor)

		if err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		w.WriteHeader(http.StatusOK)
		return

	case http.MethodDelete:
		err = removeAuthor(id)

		if err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		w.WriteHeader(http.StatusOK)
		return

	case http.MethodOptions:
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
