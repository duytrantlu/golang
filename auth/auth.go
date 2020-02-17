package auth

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	u "app/utils"
	"os"
	"strings"
	"app/models"
)

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/api/user/new", "/api/user/login"} // List of endpoint don't require auth
		requestUrl := r.URL.Path // current Url

		// check if current url have need to auth
		for _, value := range notAuth {
			if value == requestUrl {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string] interface{})
		tokenHeader := r.Header.Get("Authentication")

		if tokenHeader == "" { // Token is missing, return code 403 authorized
			response = u.Message(false, "Missing token in header")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Response(w, response)
			return
		}

		tokenSpilted := strings.Split(tokenHeader, " ")

		if len(tokenSpilted) != 2 { // Because token have format `Bear {token-body}`
			response = u.Message(false, "Invalid format token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Response(w, response)
			return
		}

		tokenPart := tokenSpilted[1]
		tokenStruct := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart,tokenStruct, func(token *jwt.Token) (i interface{}, err error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil {
			response = u.Message(false, "Malformed authentication token ")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Response(w, response)
			return
		}

		if !token.Valid {
			response = u.Message(false, "Invalid token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Response(w, response)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		fmt.Printf("User %s", tokenStruct.UserId)
		ctx := context.WithValue(r.Context(), "user", tokenStruct.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
