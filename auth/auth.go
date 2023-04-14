package auth

import (
	"net/http"
)

type Auth interface {
	// middleware that authenticates users
	Authenticate(http.Handler) http.Handler
}

type contextKey string

const UserContextKey = contextKey("user")
