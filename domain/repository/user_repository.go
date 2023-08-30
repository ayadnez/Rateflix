package repository

import "zahid/movies/domain/model"

type UserRepository interface {
	GetByEmail(email string) (*model.User, error)
	Create(user *model.User) error
}
