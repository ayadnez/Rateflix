package handler

import (
	"zahid/movies/middleware"

	"github.com/labstack/echo/v4"
)

func InitMovieRouting(e *echo.Echo, movieHandler MovieHandler) {
	// list of all movies or a specific movie if there's an id
	e.GET("/api/movies", movieHandler.Show()) //, middleware.IsAuthenticated)
	// rate a movie
	e.POST("/api/rate", movieHandler.Rate()) //, middleware.IsAuthenticated)
	// sync with moviedb
	e.POST("/api/sync", movieHandler.Sync()) //, middleware.IsAuthenticated)
}
func InitAuthRouting(e *echo.Echo, authHandler AuthHandler) {
	// e.GET("/api/auth", authHandler.Get())
	e.POST("/api/login", authHandler.Create())
	e.POST("/api/logout", authHandler.Delete(), middleware.IsAuthenticated)
	e.POST("/api/signup", authHandler.CreateUser())
}
