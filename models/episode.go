package models

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gopkg.in/validator.v2"
)

type Episode struct {
	Id          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	TmdbId      int       `json:"tmdb_id" gorm:"index:idx_episode,unique" validate:"nonzero"`
	Season      int       `json:"season" gorm:"index:idx_episode,unique" validate:"nonzero"`
	Episode     int       `json:"episode" gorm:"index:idx_episode,unique" validate:"nonzero"`
	Name        string    `json:"name" validate:"min=2,max=80"`
	Overview    string    `json:"overview"`
	AirDate     int       `json:"air_date" validate:"checkDate"`
	Watched     bool      `json:"watched"`
	WatchedDate int       `json:"watched_date" validate:"checkDate"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (e *Episode) TrimSpace() {
	e.Name = strings.TrimSpace(e.Name)
	e.Overview = strings.TrimSpace(e.Overview)
}

// Validator

func ValidEpisode(episode *Episode) error {
	episode.TrimSpace()
	validator.SetValidationFunc("checkDate", checkDate)
	if err := validator.Validate(episode); err != nil {
		return err
	}
	return nil
}

func checkDate(v interface{}, _ string) error {
	st := reflect.ValueOf(v)
	value := int(st.Int())

	if value == 0 {
		return nil
	}

	if len(strconv.Itoa(value)) != 8 || value < 0 {
		return errors.New("date must be YYYYMMDD")
	}

	if _, err := time.Parse("20060102", strconv.Itoa(value)); err != nil {
		return err
	}

	return nil
}
