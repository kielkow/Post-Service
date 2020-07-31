package post

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/kielkow/Post-Service/shared/providers/apperror"
	"github.com/kielkow/Post-Service/modules/author"
	"github.com/kielkow/Post-Service/shared/http/cors"
)

const postsBasePath = "posts"

// SetupRoutes function
func SetupRoutes(apiBasePath string) {
	handlePosts := http.HandlerFunc(postsHandler)
	handlePost := http.HandlerFunc(postHandler)

	reportHandler := http.HandlerFunc(handlePostReport)

	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, postsBasePath), cors.Middleware(handlePosts))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, postsBasePath), cors.Middleware(handlePost))

	http.Handle(fmt.Sprintf("%s/%s/reports", apiBasePath, postsBasePath), cors.Middleware(reportHandler))
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		postList, err := getPostList()

		if err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		postsJSON, err := json.Marshal(postList)

		if err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(postsJSON)

	case http.MethodPost:
		var newPost CreatePost
		bodyBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			error := apperror.GenerateError(400, err.Error())

			w.WriteHeader(http.StatusBadRequest)
			w.Write(error)
			return
		}

		err = json.Unmarshal(bodyBytes, &newPost)

		if err != nil {
			error := apperror.GenerateError(400, err.Error())

			w.WriteHeader(http.StatusBadRequest)
			w.Write(error)
			return
		}

		authorExists, err := author.GetAuthor(newPost.AuthorID)

		if err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		if authorExists == nil {
			error := apperror.GenerateError(404, "Author not found")

			w.WriteHeader(http.StatusNotFound)
			w.Write(error)
			return
		}

		_, err = insertPost(newPost)

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

func postHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "posts/")
	id, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])

	if err != nil {
		error := apperror.GenerateError(500, err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(error)
		return
	}

	post, err := getPost(id)

	if err != nil {
		error := apperror.GenerateError(500, err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(error)
		return
	}

	if post == nil {
		error := apperror.GenerateError(404, "Post not found")

		w.WriteHeader(http.StatusNotFound)
		w.Write(error)
		return
	}

	switch r.Method {
	case http.MethodGet:
		postJSON, err := json.Marshal(post)

		if err != nil {
			error := apperror.GenerateError(500, err.Error())
	
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(postJSON)

	case http.MethodPut:
		var updatedPost UpdatePost

		bodyBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			error := apperror.GenerateError(400, err.Error())
	
			w.WriteHeader(http.StatusBadRequest)
			w.Write(error)
			return
		}

		err = json.Unmarshal(bodyBytes, &updatedPost)

		if err != nil {
			error := apperror.GenerateError(400, err.Error())
	
			w.WriteHeader(http.StatusBadRequest)
			w.Write(error)
			return
		}

		authorExists, err := author.GetAuthor(updatedPost.AuthorID)

		if err != nil {
			error := apperror.GenerateError(500, err.Error())
	
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		if authorExists == nil {
			error := apperror.GenerateError(404, "Author not found")
	
			w.WriteHeader(http.StatusNotFound)
			w.Write(error)
			return
		}

		err = updatePost(id, updatedPost)

		if err != nil {
			error := apperror.GenerateError(500, err.Error())
	
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		w.WriteHeader(http.StatusOK)
		return

	case http.MethodDelete:
		err = removePost(id)

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
