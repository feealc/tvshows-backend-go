package models

import (
	"errors"
	"reflect"
	"strconv"
	"time"

	"gopkg.in/validator.v2"
)

type Episode struct {
	TvShowId    int       `json:"tvshow_id" gorm:"primaryKey" validate:"nonzero"`
	Season      int       `json:"season" gorm:"primaryKey" validate:"nonzero"`
	Episode     int       `json:"episode" gorm:"primaryKey" validate:"nonzero"`
	Name        string    `json:"name" validate:"min=2,max=80"`
	AirDate     int       `json:"air_date" validate:"checkDate"`
	Watched     bool      `json:"watched"`
	WatchedDate int       `json:"watched_date" validate:"checkDate"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validator

func ValidEpisode(episode *Episode) error {
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
