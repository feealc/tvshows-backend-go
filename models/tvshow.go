package models

import (
	"errors"
	"reflect"

	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

type TvShow struct {
	gorm.Model
	Name      string `json:"name" validate:"min=2,max=80"`
	GroupType int    `json:"group" validate:"checkGroup"`
	Status    int    `json:"status" validate:"checkStatus"`
	TmdbId    int    `json:"tmdb_id" validate:"nonzero"`
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
