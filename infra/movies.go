package infra

import (
	"fmt"
	"time"
	"zahid/movies/domain/model"
	"zahid/movies/domain/repository"
)

type MovieRepository struct {
	SqlHandler
}

func NewMovieRepository(sqlHandler SqlHandler) repository.MovieRepository {
	return MovieRepository{SqlHandler: sqlHandler}
}

func (movieRepository MovieRepository) Upsert(m model.Movie) error {
	m.UpdatedAt = time.Now()
	query := `INSERT INTO movies (id, title, is_adult, popularity, created_at, updated_at) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE popularity = VALUES(popularity), updated_at = VALUES(updated_at);`
	if _, err := movieRepository.SqlHandler.Conn.Exec(query, m.ID, m.Title, m.IsAdult, m.Popularity, m.UpdatedAt, m.UpdatedAt); err != nil {
		return err
	}
	return nil
}

func (mr MovieRepository) Rate(userId int, rating float64, movieId int) error {
	query := `INSERT INTO rating_user_mapping (user_id, movie_id, rating, created_at, updated_at) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE user_id = VALUES(user_id), rating = VALUES(rating), updated_at = VALUES(updated_at);`
	if _, err := mr.SqlHandler.Conn.Exec(query, userId, rating, movieId, time.Now(), time.Now()); err != nil {
		return err
	}
	return nil
}

func (mr MovieRepository) GetAverageRating(userId int, movieId int) (float64, error) {
	var rating float64
	query := `SELECT AVG(*) FROM rating_user_mapping WHERE user_id = ? AND movie_id = ?;`
	rows, err := mr.SqlHandler.Conn.Queryx(query, userId, movieId)
	if err != nil {
		err = fmt.Errorf("failed to select rating: %w", err)
		return rating, err
	}
	defer rows.Close()

	err = rows.StructScan(&rating)
	if err != nil {
		err = fmt.Errorf("failed to scan rating: %w", err)
		return rating, err
	}
	return rating, nil
}

func (mr MovieRepository) UpdateRating(movieId int, rating float64) error {
	query := `UPDATE movies SET rating = ? WHERE id = ?;`
	_, err := mr.SqlHandler.Conn.Exec(query, rating, movieId)
	if err != nil {
		err = fmt.Errorf("failed to update rating: %w", err)
		return err
	}
	return nil
}

func (mr MovieRepository) ShowAll() ([]model.Movie, error) {
	var movies []model.Movie
	query := `
		SELECT id,title,is_adult,rating,popularity,created_at,updated_at
		FROM movies LIMIT 1000;
	`
	rows, err := mr.SqlHandler.Conn.Queryx(query)
	if err != nil {
		err = fmt.Errorf("failed to select movies: %w", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie model.Movie
		err := rows.StructScan(&movie)
		if err != nil {
			err = fmt.Errorf("failed to scan movie: %w", err)
			return nil, err
		}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		err = fmt.Errorf("failed to iterate over movies: %w", err)
		return nil, err
	}

	return movies, nil
}

func (mr MovieRepository) Show(movieId int) (model.Movie, error) {
	var mov model.Movie
	query := `
		SELECT id,title,is_adult,rating,popularity,created_at,updated_at
		FROM movies WHERE movie_id = ?;
	`
	rows, err := mr.SqlHandler.Conn.Queryx(query, movieId)
	if err != nil {
		err = fmt.Errorf("failed to select movie: %w", err)
		return mov, err
	}
	defer rows.Close()

	err = rows.StructScan(&mov)
	if err != nil {
		err = fmt.Errorf("failed to scan movie: %w", err)
		return mov, err
	}
	return mov, nil
}
