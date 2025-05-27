package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"gopkg.in/validator.v2"
)

type TvShow struct {
	Id               int       `json:"id" gorm:"primaryKey;autoIncrement"`
	TmdbId           int       `json:"tmdb_id" gorm:"uniqueIndex" validate:"nonzero"`
	Name             string    `json:"name" gorm:"uniqueIndex" validate:"min=2,max=80"`
	Overview         string    `json:"overview"`
	GroupType        int       `json:"group" validate:"checkGroup"`
	Status           int       `json:"status" validate:"checkStatus"`
	UnwatchedSeason  int       `json:"unwatched_season"`
	UnwatchedEpisode int       `json:"unwatched_episode"`
	UnwatchedCount   int       `json:"unwatched_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (t *TvShow) TrimSpace() {
	t.Name = strings.TrimSpace(t.Name)
	t.Overview = strings.TrimSpace(t.Overview)
}

func (t *TvShow) Dump() {
	fmt.Printf("Id=%d TmdbId=%d Name=[%s] Overview=[%s] Group=%d Status=%d UnwSeason=%d UnwEp=%d UnwCount=%d Cr=%s Up=%s \n",
		t.Id,
		t.TmdbId,
		t.Name,
		t.Overview,
		t.GroupType,
		t.Status,
		t.UnwatchedSeason,
		t.UnwatchedEpisode,
		t.UnwatchedCount,
		t.CreatedAt.Format("2006-01-02 15:04:05"),
		t.UpdatedAt.Format("2006-01-02 15:04:05"),
	)
}

func (t *TvShow) DumpShort() {
	fmt.Printf("Id=%d TmdbId=%d Name=[%s] Group=%d Status=%d \n",
		t.Id,
		t.TmdbId,
		t.Name,
		t.GroupType,
		t.Status,
	)
}

// Validator

func ValidTvShow(tvShow *TvShow) error {
	tvShow.TrimSpace()
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
