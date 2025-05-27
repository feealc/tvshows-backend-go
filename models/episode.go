package models

import (
	"errors"
	"fmt"
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

func (e *Episode) Dump() {
	fmt.Printf("Id=%d TmdbId=%d T%dE%d Name=[%s] Overview=[%s] AirDate=%d Watched=%t WatchedDate=%d Cr=%s Up=%s \n",
		e.Id,
		e.TmdbId,
		e.Season,
		e.Episode,
		e.Name,
		e.Overview,
		e.AirDate,
		e.Watched,
		e.WatchedDate,
		e.CreatedAt.Format("2006-01-02 15:04:05"),
		e.UpdatedAt.Format("2006-01-02 15:04:05"),
	)
}

func (e *Episode) DumpShort() {
	fmt.Printf("Id=%d TmdbId=%d T%dE%d Name=[%s] AirDate=%d Watched=%t WatchedDate=%d \n",
		e.Id,
		e.TmdbId,
		e.Season,
		e.Episode,
		e.Name,
		e.AirDate,
		e.Watched,
		e.WatchedDate,
	)
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
