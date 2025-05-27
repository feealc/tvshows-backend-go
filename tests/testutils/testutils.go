package testutils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/feealc/tvshows-backend-go/controllers"
	"github.com/feealc/tvshows-backend-go/database"
	"github.com/feealc/tvshows-backend-go/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetUpTestRoutes(connectDb bool) *gin.Engine {
	if connectDb {
		database.ConnectDataBase()
	}
	gin.SetMode(gin.TestMode)
	// routes := gin.Default()
	routes := gin.New()
	routes.Use(gin.Recovery())
	return routes
}

func CheckResponseError(r *gin.Engine, t *testing.T, url string, body interface{}, statusCode int, errorMessage string) {
	tvShowJson, err := json.Marshal(body)
	assert.Nil(t, err)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(tvShowJson)))
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	type ResponseError struct {
		Error string `json:"error"`
	}

	resp := ResponseError{Error: errorMessage}
	respJson, err := json.Marshal(resp)
	assert.Nil(t, err)

	assert.Equal(t, statusCode, w.Code)
	if errorMessage != "" {
		assert.Equal(t, string(respJson), w.Body.String())
	}
}

func CheckTvShow(t *testing.T, tvshow models.TvShow, tvShowExpected models.TvShow) {
	assert.Equal(t, tvShowExpected.Id, tvshow.Id)
	assert.Equal(t, tvShowExpected.TmdbId, tvshow.TmdbId)
	assert.Equal(t, tvShowExpected.Name, tvshow.Name)
	assert.Equal(t, tvShowExpected.Overview, tvshow.Overview)
	assert.Equal(t, tvShowExpected.GroupType, tvshow.GroupType)
	assert.Equal(t, tvShowExpected.Status, tvshow.Status)
	assert.Equal(t, tvShowExpected.UnwatchedSeason, tvshow.UnwatchedSeason)
	assert.Equal(t, tvShowExpected.UnwatchedEpisode, tvshow.UnwatchedEpisode)
	assert.Equal(t, tvShowExpected.UnwatchedCount, tvshow.UnwatchedCount)
	// created at
	// updated at
	// if checkCreatedUpdated {
	// 	assert.Equal(t, tvShowExpected.CreatedAt, tvShowExpected.UpdatedAt)
	// }

	// if tvShowExpected.Id > 0 {

	// }
}

func CheckEpisode(t *testing.T, episode models.Episode, episodeExpected models.Episode) {
	assert.Equal(t, episodeExpected.Id, episode.Id)
	assert.Equal(t, episodeExpected.TmdbId, episode.TmdbId)
	assert.Equal(t, episodeExpected.Season, episode.Season)
	assert.Equal(t, episodeExpected.Episode, episode.Episode)
	assert.Equal(t, episodeExpected.Name, episode.Name)
	assert.Equal(t, episodeExpected.Overview, episode.Overview)
	assert.Equal(t, episodeExpected.AirDate, episode.AirDate)
	assert.Equal(t, episodeExpected.Watched, episode.Watched)
	assert.Equal(t, episodeExpected.WatchedDate, episode.WatchedDate)
	// created at
	// updated at
	// if checkCreatedUpdated {
	// 	assert.Equal(t, tvShowExpected.CreatedAt, tvShowExpected.UpdatedAt)
	// }

	// if tvShowExpected.Id > 0 {

	// }
}

//

func ListAllTvShows(debug bool) (int, error) {
	r := SetUpTestRoutes(true)
	url := "/tvshows"
	r.GET(url, controllers.TvShowListAll)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	r.ServeHTTP(w, req)

	var tvShowsResp []models.TvShow
	aux, _ := io.ReadAll(w.Body)
	err := json.Unmarshal(aux, &tvShowsResp)
	if err != nil {
		if debug {
			println("erro unmarshal")
			println(err.Error())
		}
		return 0, err
	}
	total := len(tvShowsResp)
	if debug {
		fmt.Printf("ListAllTvShows() - total [%d] \n", total)
		for _, tmp := range tvShowsResp {
			tmp.DumpShort()
		}
		// println(w.Body)
	}
	return total, nil
}

func CheckListAllTvShows(t *testing.T, debug bool, totalExpected int) {
	total, err := ListAllTvShows(debug)
	assert.Nil(t, err)
	assert.Equal(t, totalExpected, total)
}

func ListAllEpisodes(debug bool) (int, error) {
	r := SetUpTestRoutes(true)
	url := "/episodes"
	r.GET(url, controllers.EpisodeListAll)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	r.ServeHTTP(w, req)

	var episodesResp []models.Episode
	aux, _ := io.ReadAll(w.Body)
	err := json.Unmarshal(aux, &episodesResp)
	if err != nil {
		if debug {
			println("erro unmarshal")
			println(err.Error())
		}
		return 0, err
	}
	total := len(episodesResp)
	if debug {
		fmt.Printf("ListAllEpisodes() - total [%d] \n", total)
		for _, tmp := range episodesResp {
			tmp.DumpShort()
		}
		// println(w.Body)
	}
	return total, nil
}

func CheckListAllEpisodes(t *testing.T, debug bool, totalExpected int) {
	total, err := ListAllEpisodes(debug)
	assert.Nil(t, err)
	assert.Equal(t, totalExpected, total)
}
