package usecase

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"zahid/movies/domain/model"
	"zahid/movies/domain/repository"
	"zahid/movies/util"
)

type MovieUsecase interface {
	Sync() []error
	Rate(email string, rating float64, movieId int) error
	Show(movieId int) ([]model.Movie, error)
}

// Since the repository is only referenced from the usecase, create a lowercase struct
type movieUsecase struct {
	movieRepo repository.MovieRepository
	userRepo  repository.UserRepository
}

func NewMovieUsecase(movieRepo repository.MovieRepository, userRepo repository.UserRepository) MovieUsecase {
	return &movieUsecase{
		movieRepo: movieRepo,
		userRepo:  userRepo,
	}
}
func (mu *movieUsecase) Show(movieId int) ([]model.Movie, error) {
	if movieId > 0 {
		movie, err := mu.movieRepo.Show(movieId)
		if err != nil {
			return nil, err
		}
		movies := make([]model.Movie, 1)
		return append(movies, movie), nil
	}
	movies, err := mu.movieRepo.ShowAll()
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (usecase *movieUsecase) Rate(email string, rating float64, movieId int) error {
	user, err := usecase.userRepo.GetByEmail(email)
	if err != nil {
		return err
	}
	err = usecase.movieRepo.Rate(int(user.Id), rating, movieId)
	if err != nil {
		return err
	}
	rating, err = usecase.movieRepo.GetAverageRating(int(user.Id), movieId)
	if err != nil {
		return err
	}
	err = usecase.movieRepo.UpdateRating(movieId, rating)
	if err != nil {
		return err
	}
	return nil
}
func (usecase *movieUsecase) Sync() []error {
	var errors []error
	// http://files.tmdb.org/p/exports/movie_ids_05_15_2023.json.gz

	today := time.Now()
	formattedDate := today.Format("01_02_2006")
	url := util.FILES_TMDB_MOVIES_HOST + formattedDate + ".json.gz"
	// url := "https://download1654.mediafire.com/a4h0osb7sowg2GujgfCrpHaSH390PYzycUgS-TJ98YGnnzDGoniCr7JqJtJ7xP6rMZlkW0W4yLsgmzsKQhwrLHh9uS8J9cHTKnbwzpQUjV1p2sqOELsuU5OswmH0NUrTmnJz3sk57ZM3gG3RX4hmZz0h8mx7Lq_oozYlx9GsB2pH/pwhe8jrj44lx98b/x.json.gz"
	// make http request and get the data
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return append(errors, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("HTTP request failed with status code:", response.StatusCode)
		return append(errors, fmt.Errorf("status 200 not returned"))
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return append(errors, fmt.Errorf("status 200 not returned"))
	}
	// Create a bytes reader from the response body
	responseBuffer := bytes.NewBuffer(responseBody)
	// Create a gzip reader from the response buffer
	gzipReader, errs := gzip.NewReader(responseBuffer)
	if errs != nil {
		return append(errors, err)
	}
	defer gzipReader.Close()

	var movies []model.Movie
	scanner := bufio.NewScanner(gzipReader)
	for scanner.Scan() {
		line := scanner.Text()
		var movie model.Movie
		if err := parseLine(line, &movie); err != nil {
			// maybe log it
			fmt.Println("Error parsing line:", err)
			continue
		}
		movies = append(movies, movie)
	}

	// if err := scanner.Err(); err != nil {
	// 	errors = append(errors, err)
	// }

	for _, movie := range movies {
		err := usecase.movieRepo.Upsert(movie)
		if err != nil {
			errors = append(errors, fmt.Errorf("%d", movie.ID))
		}
	}

	return errors
}

func parseLine(line string, movie *model.Movie) error {
	err := json.Unmarshal([]byte(line), movie)
	if err != nil {
		return err
	}
	return nil
}
