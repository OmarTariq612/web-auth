package auth

import (
	"net/http"
	"time"
)

type Kind int

const (
	None Kind = iota
	Cookie
	Body
)

type Token struct {
	Kind  Kind
	Value string
}

type Auth interface {
	// middleware that authenticates users
	Authenticate(http.Handler) http.Handler
	GenerateToken(any) (*Token, error)
}

type contextKey string

const UserContextKey = contextKey("user")

const Name = "token"
const TTL = 2 * time.Hour
