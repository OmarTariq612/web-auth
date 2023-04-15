package data

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
}

func generateToken(userID int64, ttl time.Duration) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)

	if _, err := rand.Read(randomBytes); err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

type TokenModel struct {
	DB *sql.DB
}

func (m TokenModel) Insert(userID int64, ttl time.Duration) (*Token, error) {
	token, err := generateToken(userID, ttl)
	if err != nil {
		return nil, err
	}

	query := `
	INSERT INTO tokens (user_id, token_hash, expiry)
	VALUES (?, ?, ?)`

	args := []any{userID, token.Hash, token.Expiry}

	_, err = m.DB.Exec(query, args...)
	return token, err
}

func (m TokenModel) GetUserIDFromToken(tokenPlaintext string) (int64, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
	SELECT user_id FROM tokens
	WHERE token_hash = ? AND expiry > ?`

	args := []any{tokenHash[:], time.Now()}

	var userID int64

	if err := m.DB.QueryRow(query, args...).Scan(&userID); err != nil {
		return 0, err
	}

	return userID, nil
}

func (m TokenModel) DeleteAllForUser(userID int64) error {
	query := `
	DELETE FROM tokens WHERE user_id = ?`

	_, err := m.DB.Exec(query, userID)
	return err
}

var _ TokenDao = (*TokenModel)(nil)
