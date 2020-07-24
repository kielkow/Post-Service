package post

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		postsJSON, err := json.Marshal(postList)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(postsJSON)

	case http.MethodPost:
		var newPost Post
		bodyBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &newPost)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if newPost.ProductID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = insertPost(newPost)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	post, err := getPost(productID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if post == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		postJSON, err := json.Marshal(post)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(postJSON)

	case http.MethodPut:
		var updatedPost Post

		bodyBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &updatedPost)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if updatedPost.ProductID != productID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = updatePost(updatedPost)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return

	case http.MethodDelete:
		err = removePost(productID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
