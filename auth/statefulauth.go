package auth

import (
	"net/http"

	"github.com/OmarTariq612/web-auth/data"
)

type statefulAuth struct {
	userdao data.UserDao
}

func NewStatefulAuth(userdao data.UserDao) *statefulAuth {
	return &statefulAuth{
		userdao: userdao,
	}
}

func (auth *statefulAuth) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO:
		next.ServeHTTP(w, r)
	})
}

func (auth *statefulAuth) GenerateAuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO:
		next.ServeHTTP(w, r)
	})
}

var _ Auth = (*statefulAuth)(nil)
