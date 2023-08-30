package repository

import "zahid/movies/domain/model"

type MovieRepository interface {
	Rate(userId int, rating float64, movieId int) error
	GetAverageRating(userId int, movieId int) (float64, error)
	UpdateRating(movieId int, rating float64) error
	Upsert(movie model.Movie) error
	ShowAll() ([]model.Movie, error)
	Show(movieId int) (model.Movie, error)
}
