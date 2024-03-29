package auth

import (
	"context"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/wilsonth122/money-tracker-api/pkg/model"
	u "github.com/wilsonth122/money-tracker-api/pkg/utils"
)

// JwtAuthentication - Authenticates the token in the header of a request
var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/api/user/new", "/api/user/login", "/api/stream/expenses"}
		requestPath := r.URL.Path

		// Check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization")

		// Token is missing, returns with error code 403 Unauthorized
		if tokenHeader == "" {
			u.RespondWithError(w, http.StatusForbidden, "Missing auth token")
			return
		}

		// The token normally comes in format `Bearer {token-body}`,
		// we check if the retrieved token matched this requirement
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			u.RespondWithError(w, http.StatusForbidden, "Invalid/Malformed auth token")
			return
		}

		// Grab the token part, what we are truly interested in
		tokenStr := splitted[1]
		token, err := ParseToken(tokenStr)
		if err != nil {
			u.RespondWithError(w, http.StatusForbidden, err.Error())
			return
		}

		// Useful for monitoring
		log.Println("User " + token.UserID)

		// Everything went well,
		// proceed with the request and set the caller to the user retrieved from the parsed token
		ctx := context.WithValue(r.Context(), "user", token.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// ParseToken - Parses and encrypted token string into a token
func ParseToken(tokenStr string) (model.Token, error) {
	tk := model.Token{}

	token, err := jwt.ParseWithClaims(tokenStr, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_PASSWORD")), nil
	})

	// Malformed token
	if err != nil {
		return tk, err
	}

	// Token is invalid, maybe not signed on this server
	if !token.Valid {
		return tk, errors.New("Token is not valid")
	}

	return tk, nil
}
