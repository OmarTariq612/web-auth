package auth

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/OmarTariq612/web-auth/data"
	"github.com/golang-jwt/jwt/v4"
)

type jwtObject struct {
	jwt.RegisteredClaims
	UserID int64 `json:"user_id"`
}

func signJWT(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	return tokenString, err
}

func verifyJWTCustom(tokenString string, claims jwt.Claims) (verifiedClaims jwt.Claims, err error, expired bool) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	if err != nil {
		// expired (signature is valid but token is expired)
		if er, ok := err.(*jwt.ValidationError); ok && !er.Is(jwt.ErrSignatureInvalid) && er.Is(jwt.ErrTokenExpired) {
			return nil, err, true
		}
		return nil, err, false // invalid or anything else
	}
	return token.Claims, nil, false
}

type statelessAuth struct {
	userdao data.UserDao
}

func NewStatelessAuth(userdao data.UserDao) *statelessAuth {
	return &statelessAuth{
		userdao: userdao,
	}
}

func (auth *statelessAuth) extractAuthToken(r *http.Request) (*data.User, bool) {
	authTokenHeader := r.Header.Get("Authorization")
	if authTokenHeader == "" {
		return nil, false
	}
	params := strings.Split(authTokenHeader, " ")
	if params[0] != "Bearer" {
		return nil, false
	}
	authToken := params[1]
	var obj jwtObject
	// TODO: differintiate between expired tokens and invalid tokens
	_, err, _ := verifyJWTCustom(authToken, &obj)
	if err != nil {
		return nil, false
	}
	user, err := auth.userdao.GetByID(obj.UserID)
	if err != nil {
		return nil, false
	}

	return user, true
}

func (auth *statelessAuth) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.extractAuthToken(r)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), UserContextKey, user))

		next.ServeHTTP(w, r)
	})
}

func (auth *statelessAuth) GenerateToken(user any) (*Token, error) {
	token, err := signJWT(jwtObject{RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(TTL))}, UserID: user.(*data.User).ID})
	if err != nil {
		return nil, err
	}

	return &Token{
		Kind:  Body,
		Value: token,
	}, nil
}

var _ Auth = (*statelessAuth)(nil)
