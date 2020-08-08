package authentication

import (
	"fmt"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kielkow/Post-Service/shared/providers/apperror"
)

var mySigningKey = []byte(os.Getenv("JWT_APP_SECRET"))

// GenerateJWT func
func GenerateJWT(authorEmail string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = authorEmail
	claims["exp"] = time.Now().Add(time.Minute * 2880).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Print(err)
		return "", nil
	}

	return tokenString, nil
}

// IsAuthorized func
func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					error := apperror.GenerateError(500, "There was an error on JWT authenticaction")

					w.WriteHeader(http.StatusInternalServerError)
					w.Write(error)
				}

				return mySigningKey, nil
			})

			if err != nil {
				error := apperror.GenerateError(500, err.Error())

				w.WriteHeader(http.StatusInternalServerError)
				w.Write(error)
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			error := apperror.GenerateError(403, "Not Authorized")

			w.WriteHeader(http.StatusNonAuthoritativeInfo)
			w.Write(error)
		}
	})
}
