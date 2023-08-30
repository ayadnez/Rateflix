package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"zahid/movies/usecase"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type MovieHandler struct {
	movieUsecase usecase.MovieUsecase
}

func NewMovieHandler(movieUsecase usecase.MovieUsecase) MovieHandler {
	return MovieHandler{
		movieUsecase: movieUsecase,
	}
}

type showRequest struct {
	MovieId int `json:"movieId"`
}

type rateRequest struct {
	Rating  float64 `json:"rating" binding:"required"`
	MovieId int     `json:"movieId" binding:"required"`
}

func (r rateRequest) Validate() error {
	if r.Rating > 5 && r.Rating < 0 {
		return fmt.Errorf("rating invalid")
	}
	if r.MovieId < 0 {
		return fmt.Errorf("invalid movie id")
	}
	return nil
}

func (handler *MovieHandler) Show() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req showRequest
		var mid int
		if err := c.Bind(&req); err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusBadRequest, err)
		}
		if req.MovieId > 0 {
			mid = req.MovieId
		}
		data, err := handler.movieUsecase.Show(mid)
		if err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusUnprocessableEntity, err)
		}
		return c.JSON(http.StatusOK, echo.Map{
			"message": "fetched successfully",
			"data":    data,
		})
	}
}
func (handler *MovieHandler) Rate() echo.HandlerFunc {
	return func(c echo.Context) error {
		email, err := handler.GetEmailId(c)
		if err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusUnauthorized, err)
		}
		var req rateRequest
		if err := c.Bind(&req); err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusBadRequest, err)
		}
		if err := req.Validate(); err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusBadRequest, err)
		}
		err = handler.movieUsecase.Rate(email, req.Rating, req.MovieId)
		if err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusUnprocessableEntity, err)
		}
		return c.JSON(http.StatusOK, echo.Map{
			"message": "rated successfully",
		})
	}
}

func (handler *MovieHandler) Sync() echo.HandlerFunc {
	return func(c echo.Context) error {
		errs := handler.movieUsecase.Sync()
		return c.JSON(http.StatusOK, echo.Map{
			"failed_ids": errs,
		})
	}
}

func (handler *MovieHandler) GetEmailId(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorisation")
	t := strings.Split(authHeader, " ")
	if len(t) != 2 {
		return "", fmt.Errorf("invalid token")
	}
	token, err := handler.parseJwt(t[1])
	if err != nil {
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims["iss"].(string), nil
}

func (handler *MovieHandler) parseJwt(tokenString string) (*jwt.Token, error) {
	key := []byte(os.Getenv("JWT_SECRET_KEY"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		message := "JWT parsing failed"
		fmt.Println(message, err)

		return nil, err
	}
	if token.Claims.Valid() != nil {
		fmt.Println("Invalid JWT token:", token.Claims.Valid())
		return nil, errors.New("Invalid JWT token")
	}
	return token, nil
}
