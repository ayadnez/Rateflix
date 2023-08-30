package infra

import (
	"time"
	"zahid/movies/domain/model"
	"zahid/movies/domain/repository"
)

type UserRepository struct {
	SqlHandler
}

func NewUserRepository(sqlHandler SqlHandler) repository.UserRepository {
	return UserRepository{SqlHandler: sqlHandler}
}

func (userRepository UserRepository) Create(user *model.User) error {
	query := `INSERT INTO users (name, password, email, created, updated) VALUES (?, ?, ?, ?, ?);`

	now := time.Now()

	if _, err := userRepository.SqlHandler.Conn.Exec(query, user.Name, user.Password, user.Email, now, now); err != nil {
		return err
	}

	return nil
}

func (userRepository UserRepository) GetByEmail(email string) (*model.User, error) {
	query := `SELECT * FROM users WHERE email = ?;`

	var user model.User
	if err := userRepository.SqlHandler.Conn.Get(&user, query, email); err != nil {
		return nil, err
	}

	return &user, nil
}
