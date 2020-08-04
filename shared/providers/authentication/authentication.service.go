package authentication

import (
	"fmt"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte(os.Getenv("JWT_APP_SECRET"))

// GenerateJWT func
func GenerateJWT(authorName string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = authorName
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
	return nil
}
