package auth

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

type Store struct {
	*sql.DB
}

type AuthStore interface {
	RegisterUser(user *User) error
	GetUserByUserName(username string) (*User, error)
}

func (s *Store) RegisterUser(user *User) error {
	_, err := s.Exec("INSERT INTO users (username, password_hash, created_at) VALUES($1, $2, $3)", user.Username, user.PasswordHash, user.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetUserByUserName(username string) (*User, error) {
	row := s.QueryRow("SELECT id, username, password_hash, created_at FROM users WHERE username = $1", username)

	user := User{}

	if err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}
