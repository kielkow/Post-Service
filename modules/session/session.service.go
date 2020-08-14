package session

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kielkow/Post-Service/shared/http/cors"
	"github.com/kielkow/Post-Service/shared/providers/apperror"
	"github.com/kielkow/Post-Service/shared/providers/authentication"
	"github.com/kielkow/Post-Service/shared/providers/hasher"
)

const sessionBasePath = "session"

// SetupRoutes function
func SetupRoutes(apiBasePath string) {
	handleSession := http.HandlerFunc(postsHandler)

	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, sessionBasePath), cors.Middleware(handleSession))
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var session Session
		bodyBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			error := apperror.GenerateError(400, err.Error())

			w.WriteHeader(http.StatusBadRequest)
			w.Write(error)
			return
		}

		err = json.Unmarshal(bodyBytes, &session)

		if err != nil {
			error := apperror.GenerateError(400, err.Error())

			w.WriteHeader(http.StatusBadRequest)
			w.Write(error)
			return
		}

		passwordHashed, err := getPassword(session.Email)

		if err != nil {
			error := apperror.GenerateError(403, err.Error())

			w.WriteHeader(http.StatusNonAuthoritativeInfo)
			w.Write(error)
			return
		}

		match := hasher.CheckPasswordHash(session.Password, passwordHashed)

		if match == false {
			error := apperror.GenerateError(401, "Incorrect e-mail/password combination")

			w.WriteHeader(http.StatusNonAuthoritativeInfo)
			w.Write(error)
			return
		}

		token, err := authentication.GenerateJWT(session.Email)

		tokenJSON, err := json.Marshal(token)

		w.Header().Set("Content-Type", "application/json")
		w.Write(tokenJSON)
		return

	case http.MethodOptions:
		return
	}
}
