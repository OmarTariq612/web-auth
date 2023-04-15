package data

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64    `json:"id"`
	Username string   `json:"username"`
	Password password `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintext
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintext string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintext)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users (username, password_hash)
	VALUES (?, ?)`

	args := []any{user.Username, user.Password.hash}

	if _, err := m.DB.Exec(query, args...); err != nil {
		return err
	}

	return nil
}

func (m UserModel) GetByUsername(username string) (*User, error) {
	query := `
	SELECT id, username, password_hash
	FROM users
	WHERE username = ?`

	var user User

	if err := m.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password.hash); err != nil {
		return nil, err
	}

	return &user, nil
}

func (m UserModel) GetByID(id int64) (*User, error) {
	query := `
	SELECT id, username, password_hash
	FROM users
	WHERE id = ?`

	var user User

	if err := m.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password.hash); err != nil {
		return nil, err
	}

	return &user, nil
}

var _ UserDao = (*UserModel)(nil)
