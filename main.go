package main

import (
	"log"
	"net/http"
	"os"
	"zahid/movies/handler"
	"zahid/movies/injector"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	e := createMux()
	setupRouting(e)

	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + port))
}

func createMux() *echo.Echo {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
		},
		AllowCredentials: true,
	}))
	e.Use(middleware.Logger())

	return e
}

func setupRouting(e *echo.Echo) {
	noteHandler := injector.InjectMovieHandler()
	handler.InitMovieRouting(e, noteHandler)

	authHandler := injector.InjectAuthHandler()
	handler.InitAuthRouting(e, authHandler)
}
