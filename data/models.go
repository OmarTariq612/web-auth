package data

import (
	"database/sql"
	"time"
)

type Models struct {
	UserDao
	TokenDao
}

type UserDao interface {
	Insert(user *User) error
	GetByUsername(username string) (*User, error)
	GetByID(id int64) (*User, error)
}

type TokenDao interface {
	Insert(userID int64, ttl time.Duration) (*Token, error)
	GetUserIDFromToken(token string) (int64, error)
	DeleteAllForUser(userID int64) error
}

func NewModels(db *sql.DB) Models {
	return Models{
		UserDao:  UserModel{DB: db},
		TokenDao: TokenModel{DB: db},
	}
}
