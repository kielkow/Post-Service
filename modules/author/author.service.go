package author

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kielkow/Post-Service/shared/http/cors"
	"github.com/kielkow/Post-Service/shared/providers/apperror"
	"github.com/kielkow/Post-Service/shared/providers/authentication"
	"github.com/kielkow/Post-Service/shared/providers/hasher"
	"github.com/kielkow/Post-Service/shared/providers/storage"
)

const authorsBasePath = "authors"

// ReceiptDirectory uploads
var ReceiptDirectory string = filepath.Join("uploads")

// SetupRoutes function
func SetupRoutes(apiBasePath string) {
	handleCreate := http.HandlerFunc(authorCreate)
	handleAuthors := http.HandlerFunc(authorsHandler)
	handleAuthor := http.HandlerFunc(authorHandler)
	reportHandler := http.HandlerFunc(handleAuthorReport)

	http.Handle(fmt.Sprintf("%s/%s/create", apiBasePath, authorsBasePath), cors.Middleware(handleCreate))
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, authorsBasePath), cors.Middleware(authentication.IsAuthorized(handleAuthors)))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, authorsBasePath), cors.Middleware(authentication.IsAuthorized(handleAuthor)))
	http.Handle(fmt.Sprintf("%s/%s/reports", apiBasePath, authorsBasePath), cors.Middleware(authentication.IsAuthorized(reportHandler)))
}

func authorCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
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

		authorExists, err := getAuthorByEmail(newAuthor.Email)

		if err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		if authorExists != nil {
			error := apperror.GenerateError(404, "Author already exists")

			w.WriteHeader(http.StatusBadRequest)
			w.Write(error)
			return
		}

		hashedPassword, err := hasher.HashPassword(newAuthor.Password)

		if err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		newAuthor.Password = hashedPassword

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
		var updatedAuthor UpdateAuthor

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

		if updatedAuthor.Password != "" {
			if updatedAuthor.ConfirmPassword != updatedAuthor.Password {
				error := apperror.GenerateError(400, "Confirm password must be equal like password")

				w.WriteHeader(http.StatusInternalServerError)
				w.Write(error)
				return
			}

			hashedPassword, err := hasher.HashPassword(updatedAuthor.Password)

			if err != nil {
				error := apperror.GenerateError(500, err.Error())

				w.WriteHeader(http.StatusInternalServerError)
				w.Write(error)
				return
			}

			updatedAuthor.Password = hashedPassword

			err = updateAuthor(id, updatedAuthor)
		} else {
			err = updateAuthor(id, updatedAuthor)
		}

		if err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		w.WriteHeader(http.StatusOK)
		return

	case http.MethodPatch:
		if author.Avatar != nil {
			err = os.Remove(filepath.Join(ReceiptDirectory, *author.Avatar))

			if err != nil {
				error := apperror.GenerateError(500, err.Error())

				w.WriteHeader(http.StatusInternalServerError)
				w.Write(error)
				return
			}

			err := removeAvatar(*author.Avatar)

			if err != nil {
				error := apperror.GenerateError(500, err.Error())

				w.WriteHeader(http.StatusInternalServerError)
				w.Write(error)
				return
			}

			err = storage.DeleteFile(*author.Avatar)

			if err != nil {
				error := apperror.GenerateError(500, err.Error())

				w.WriteHeader(http.StatusInternalServerError)
				w.Write(error)
				return
			}
		}

		r.ParseMultipartForm(5 << 20) // 5Mb
		file, handler, err := r.FormFile("avatar")

		bytes := make([]byte, 10)

		if _, err := rand.Read(bytes); err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		hashedName := hex.EncodeToString(bytes) + "-" + handler.Filename

		if err != nil {
			error := apperror.GenerateError(400, err.Error())

			w.WriteHeader(http.StatusBadRequest)
			w.Write(error)
			return
		}

		defer file.Close()

		f, err := os.OpenFile(filepath.Join(ReceiptDirectory, hashedName), os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		defer f.Close()

		io.Copy(f, file)

		err = storage.UploadFile(filepath.Join(ReceiptDirectory, hashedName), hashedName)

		if err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		newAvatar := CreateAuthorAvatar{id, hashedName}

		_, err = createAvatar(newAvatar)

		if err != nil {
			error := apperror.GenerateError(500, err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(error)
			return
		}

		w.WriteHeader(http.StatusCreated)
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
