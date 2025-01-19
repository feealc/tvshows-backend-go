package models

import (
	"errors"
	"reflect"
	"time"

	"gopkg.in/validator.v2"
)

type TvShow struct {
	TmdbId    int       `json:"tmdb_id" gorm:"primaryKey" validate:"nonzero"`
	Name      string    `json:"name" gorm:"unique" validate:"min=2,max=80"`
	GroupType int       `json:"group" validate:"checkGroup"`
	Status    int       `json:"status" validate:"checkStatus"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validator

func ValidTvShow(tvShow *TvShow) error {
	validator.SetValidationFunc("checkGroup", checkGroup)
	validator.SetValidationFunc("checkStatus", checkStatus)
	if err := validator.Validate(tvShow); err != nil {
		return err
	}
	return nil
}

func checkGroup(v interface{}, _ string) error {
	st := reflect.ValueOf(v)
	if st.Int() != 1 && st.Int() != 2 && st.Int() != 3 {
		return errors.New("value must be 1, 2 or 3")
	}
	return nil
}

func checkStatus(v interface{}, _ string) error {
	st := reflect.ValueOf(v)
	if st.Int() != 1 && st.Int() != 2 && st.Int() != 3 && st.Int() != 4 && st.Int() != 5 {
		return errors.New("value must be 1, 2, 3, 4 or 5")
	}
	return nil
}
