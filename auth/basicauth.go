package auth

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/OmarTariq612/web-auth/data"
)

type basicAuth struct {
	userdao data.UserDao
}

func NewBasicAuth(userdao data.UserDao) *basicAuth {
	return &basicAuth{
		userdao: userdao,
	}
}

func (auth *basicAuth) extractAuthToken(r *http.Request) (*data.User, bool) {
	authTokenHeader := r.Header.Get("Authorization")
	if authTokenHeader == "" {
		return nil, false
	}
	params := strings.Split(authTokenHeader, " ")
	if params[0] != "Basic" {
		return nil, false
	}
	authParams, err := base64.StdEncoding.DecodeString(params[1])
	if err != nil {
		return nil, false
	}
	credentials := strings.Split(string(authParams), ":")
	username := credentials[0]
	password := credentials[1]
	user, err := auth.userdao.GetByUsername(username)
	if err != nil {
		return nil, false
	}
	if _, err := user.Password.Matches(password); err != nil {
		return nil, false
	}

	return user, true
}

func (auth *basicAuth) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// failed to authenticate
		user, ok := auth.extractAuthToken(r)
		if !ok {
			w.Header().Add("WWW-Authenticate", "Basic")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), UserContextKey, user))

		next.ServeHTTP(w, r)
	})
}

func (auth *basicAuth) GenerateToken(any) (*Token, error) {
	return &Token{
		Kind:  None,
		Value: "",
	}, nil
}

var _ Auth = (*basicAuth)(nil)
