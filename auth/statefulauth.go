package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/OmarTariq612/web-auth/data"
)

type statefulAuth struct {
	userdao  data.UserDao
	tokendao data.TokenDao
}

func NewStatefulAuth(userdao data.UserDao, tokendao data.TokenDao) *statefulAuth {
	return &statefulAuth{
		userdao:  userdao,
		tokendao: tokendao,
	}
}

func (auth *statefulAuth) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(Name)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		userID, err := auth.tokendao.GetUserIDFromToken(cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		user, err := auth.userdao.GetByID(userID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), UserContextKey, user))

		next.ServeHTTP(w, r)
	})
}

func (auth *statefulAuth) GenerateToken(user any) (*Token, error) {
	token, err := auth.tokendao.Insert(user.(*data.User).ID, TTL)
	if err != nil {
		return nil, err
	}

	return &Token{
		Kind:  Cookie,
		Value: token.Plaintext,
	}, nil
}

var _ Auth = (*statefulAuth)(nil)
