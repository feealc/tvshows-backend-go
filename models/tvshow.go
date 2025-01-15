package models

import "gorm.io/gorm"

type TvShow struct {
	gorm.Model
	Name      string `json:"name"`
	GroupType int    `json:"group"`
	Status    int    `json:"status"`
	TmdbId    int    `json:"tmdb_id"`
}
