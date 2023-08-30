package model

import (
	"time"
)

type Movie struct {
	ID         int     `json:"id"`
	Title      string  `json:"original_title"`
	IsAdult    bool    `json:"adult"`
	Popularity float64 `json:"popularity"`
	Video      bool    `json:"video"`
	Rating     float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
