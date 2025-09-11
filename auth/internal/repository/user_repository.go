package repository

import (
	"auth/internal/model"
	"crypto/sha256"
	"encoding/hex"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	FindByUsername(username string) (*model.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	db.MustExec(`CREATE TABLE IF NOT EXISTS users (
		id INT PRIMARY KEY,
		username TEXT,
		password_hash TEXT
	)`)

	password := sha256.Sum256([]byte("password"))
	passwordHash := hex.EncodeToString(password[:])

	db.NamedExec(`INSERT INTO users (id, username, password_hash) VALUES (:id, :username, :password_hash) ON CONFLICT DO NOTHING`, &model.User{
		ID:       1,
		Username: "david",
		Password: passwordHash,
	})
	db.NamedExec(`INSERT INTO users (id, username, password_hash) VALUES (:id, :username, :password_hash) ON CONFLICT DO NOTHING`, &model.User{
		ID:       2,
		Username: "sam",
		Password: passwordHash,
	})

	return &userRepository{db: db}
}

func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var u model.User
	err := r.db.Get(&u, "SELECT id, username, password_hash FROM users WHERE username=$1", username)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
