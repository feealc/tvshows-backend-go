package generic

import (
	"errors"
	"strconv"
	"time"
)

func CheckParamsInt(paramTmdbId, paramSeason, paramEpisode string) error {
	if paramTmdbId != "" {
		_, err := strconv.Atoi(paramTmdbId)
		if err != nil {
			return errors.New("tmdbId invalid")
		}
	}

	if paramSeason != "" {
		_, err := strconv.Atoi(paramSeason)
		if err != nil {
			return errors.New("season invalid")
		}
	}

	if paramEpisode != "" {
		_, err := strconv.Atoi(paramEpisode)
		if err != nil {
			return errors.New("episode invalid")
		}
	}

	return nil
}

func GetCurrentDate() int {
	dateInt, _ := strconv.Atoi(time.Now().Format("20060102"))
	return dateInt
}
