package data

import "database/sql"

type Models struct {
	UserDao
}

type UserDao interface {
	Insert(user *User) error
	GetByUsername(username string) (*User, error)
}

func NewModels(db *sql.DB) Models {
	return Models{
		UserDao: UserModel{DB: db},
	}
}
