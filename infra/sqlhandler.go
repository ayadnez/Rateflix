package infra

import (
	"errors"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

type SqlHandler struct {
	Conn *sqlx.DB
}

func NewSqlHandler() *SqlHandler {
	// Data Source Name
	dsn, err := dsn()
	if err != nil {
		panic(err.Error)
	}

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic(err.Error)
	}
	if err := db.Ping(); err != nil {
		panic(err.Error)
	}

	return &SqlHandler{Conn: db}
}

func ConnectDB() (*sqlx.DB, error) {
	handler := NewSqlHandler()

	return handler.Conn, nil
}

func dsn() (string, error) {
	user := os.Getenv("DB_USER")
	if user == "" {
		return "", errors.New("$DB_USER is not set")
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		return "", errors.New("$DB_PASSWORD is not set")
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		return "", errors.New("$DB_PORT is not set")
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		return "", errors.New("$DB_HOST is not set")
	}

	name := os.Getenv("DB_NAME")
	if name == "" {
		return "", errors.New("$DB_NAME is not set")
	}

	options := "charset=utf8mb4&parseTime=True&loc=Local"

	// "user:password@host:port/dbname?options"
	return fmt.Sprintf("%s:%s@(%s:%s)/%s?%s",
		user, password, host, port, name, options), nil
}
