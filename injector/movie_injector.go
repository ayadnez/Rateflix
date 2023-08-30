package injector

import (
	"zahid/movies/domain/repository"
	"zahid/movies/handler"
	"zahid/movies/infra"
	"zahid/movies/usecase"
)

func InjectMovieHandler() handler.MovieHandler {
	return handler.NewMovieHandler(InjectMovieUsecase())
}

func InjectMovieUsecase() usecase.MovieUsecase {
	movieRepository := InjectMovieRepository()
	userRepository := InjectUserRepository()

	return usecase.NewMovieUsecase(movieRepository, userRepository)
}

func InjectMovieRepository() repository.MovieRepository {
	sqlHandler := InjectDB()

	return infra.NewMovieRepository(sqlHandler)
}
