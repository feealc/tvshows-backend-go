package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/feealc/tvshows-backend-go/controllers"
	"github.com/feealc/tvshows-backend-go/generic"
	"github.com/feealc/tvshows-backend-go/models"
	"github.com/feealc/tvshows-backend-go/tests/testutils"
	"github.com/stretchr/testify/assert"
)

var (
	TMDBID_CASTLE    int = 1419
	TMDBID_THEROOKIE int = 79744
	DEBUG            bool
	tvShowTest       models.TvShow
	episodesTest     []models.Episode
)

func GetEpisodeTest(id int) (models.Episode, error) {
	for _, episode := range episodesTest {
		if episode.Id == id {
			return episode, nil
		}
	}

	return models.Episode{}, fmt.Errorf("episode not found")
}

func UpdateEpisodeTest(episode models.Episode) error {
	for index, ep := range episodesTest {
		if ep.Id == episode.Id {
			episodesTest[index] = episode
			return nil
		}
	}

	return fmt.Errorf("error to update episode")
}

func DumpEpisodeTest() {
	for _, ep := range episodesTest {
		ep.Dump()
	}
}

func TestMain(m *testing.M) {
	// println("TesteMain()")
	if os.Getenv("DEBUG") == "true" {
		DEBUG = true
	}
	// fmt.Printf("DEBUG ENV [%s] DEBUG VAR [%t] \n", os.Getenv("DEBUG"), DEBUG)

	m.Run()
}

func TestHealth(t *testing.T) {
	r := testutils.SetUpTestRoutes(false)
	url := "/health"
	r.GET(url, controllers.Health)
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(t, err)
	r.ServeHTTP(w, req)

	type Response struct {
		DateTime string `json:"date_time"`
		Message  string `json:"message"`
	}
	resp := Response{DateTime: time.Now().Format("2006-01-02 15:04:05"), Message: "Ok"}
	respJson, err := json.Marshal(resp)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(respJson), w.Body.String())
}

func TestTruncateAll(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/truncate/all"
	r.DELETE(url, controllers.TruncateAll)
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	assert.Nil(t, err)
	r.ServeHTTP(w, req)

	// log.Println(w.Body.String())
	type Response struct {
		Message string `json:"message"`
	}
	resp := Response{Message: "All truncated"}
	respJson, err := json.Marshal(resp)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(respJson), w.Body.String())
}

// Tv Show

func TestTvShowCreate(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/tvshows/create"
	r.POST(url, controllers.TvShowCreate)
	w := httptest.NewRecorder()

	tvShowCreate := models.TvShow{
		TmdbId:    TMDBID_CASTLE,
		Name:      "Castle",
		Overview:  "",
		GroupType: 1,
		Status:    1,
	}
	tvShowJson, err := json.Marshal(tvShowCreate)
	assert.Nil(t, err)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(tvShowJson)))
	assert.Nil(t, err)
	r.ServeHTTP(w, req)
	// println(w.Body.String())

	var tvShowCreated models.TvShow
	err = json.Unmarshal(w.Body.Bytes(), &tvShowCreated)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, 1, tvShowCreated.Id)
	tvShowCreate.Id = tvShowCreated.Id
	testutils.CheckTvShow(t, tvShowCreate, tvShowCreated)

	tvShowTest = tvShowCreated
}

func TestTvShowCreateErrorBind(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/tvshows/create"
	r.POST(url, controllers.TvShowCreate)

	// TMDB ID
	type TvShowBindTmdbId struct {
		TmdbId string `json:"tmdb_id"`
	}
	tvShowTmdbId := TvShowBindTmdbId{
		TmdbId: "123",
	}
	testutils.CheckResponseError(r, t, url, tvShowTmdbId, http.StatusBadRequest, "json: cannot unmarshal string into Go struct field TvShow.tmdb_id of type int")

	// NAME
	type TvShowBindName struct {
		TmdbId int `json:"tmdb_id"`
		Name   int `json:"name"`
	}
	tvShowName := TvShowBindName{
		TmdbId: 123,
		Name:   0,
	}
	testutils.CheckResponseError(r, t, url, tvShowName, http.StatusBadRequest, "json: cannot unmarshal number into Go struct field TvShow.name of type string")

	// OVERVIEW
	type TvShowBindOverview struct {
		TmdbId   int    `json:"tmdb_id"`
		Name     string `json:"name"`
		Overview int    `json:"overview"`
	}
	tvShowOverview := TvShowBindOverview{
		TmdbId:   123,
		Name:     "Test",
		Overview: 2,
	}
	testutils.CheckResponseError(r, t, url, tvShowOverview, http.StatusBadRequest, "json: cannot unmarshal number into Go struct field TvShow.overview of type string")

	// GROUP
	type TvShowBindGroup struct {
		TmdbId    int    `json:"tmdb_id"`
		Name      string `json:"name"`
		Overview  string `json:"overview"`
		GroupType string `json:"group"`
	}
	tvShowGroup := TvShowBindGroup{
		TmdbId:    123,
		Name:      "Test",
		Overview:  "This is about",
		GroupType: "2",
	}
	testutils.CheckResponseError(r, t, url, tvShowGroup, http.StatusBadRequest, "json: cannot unmarshal string into Go struct field TvShow.group of type int")

	// STATUS
	type TvShowBindStatus struct {
		TmdbId    int    `json:"tmdb_id"`
		Name      string `json:"name"`
		Overview  string `json:"overview"`
		GroupType int    `json:"group"`
		Status    string `json:"status"`
	}
	tvShowStatus := TvShowBindStatus{
		TmdbId:    123,
		Name:      "Test",
		Overview:  "This is about",
		GroupType: 2,
		Status:    "0",
	}
	testutils.CheckResponseError(r, t, url, tvShowStatus, http.StatusBadRequest, "json: cannot unmarshal string into Go struct field TvShow.status of type int")
}

func TestTvShowCreateErrorValidate(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/tvshows/create"
	r.POST(url, controllers.TvShowCreate)

	// TMDB ID
	type TvShowValidateTmdbId struct {
		TmdbId    int    `json:"tmdb_id"`
		Name      string `json:"name"`
		GroupType int    `json:"group"`
		Status    int    `json:"status"`
	}
	tvShowTmdbId := TvShowValidateTmdbId{
		TmdbId:    0,
		Name:      "Test Show",
		GroupType: 1,
		Status:    1,
	}
	testutils.CheckResponseError(r, t, url, tvShowTmdbId, http.StatusUnprocessableEntity, "TmdbId: zero value")

	// NAME
	type TvShowValidateName struct {
		TmdbId    int    `json:"tmdb_id"`
		Name      string `json:"name"`
		GroupType int    `json:"group"`
		Status    int    `json:"status"`
	}
	tvShowName := TvShowValidateName{
		TmdbId:    1,
		Name:      "",
		GroupType: 1,
		Status:    1,
	}
	testutils.CheckResponseError(r, t, url, tvShowName, http.StatusUnprocessableEntity, "Name: less than min")

	// GROUP
	type TvShowValidateGroup struct {
		TmdbId    int    `json:"tmdb_id"`
		Name      string `json:"name"`
		GroupType int    `json:"group"`
		Status    int    `json:"status"`
	}
	tvShowGroup := TvShowValidateGroup{
		TmdbId:    1,
		Name:      "Test",
		GroupType: 14,
		Status:    1,
	}
	testutils.CheckResponseError(r, t, url, tvShowGroup, http.StatusUnprocessableEntity, "GroupType: value must be 1, 2 or 3")

	// STATUS
	type TvShowValidateStatus struct {
		TmdbId    int    `json:"tmdb_id"`
		Name      string `json:"name"`
		GroupType int    `json:"group"`
		Status    int    `json:"status"`
	}
	tvShowStatus := TvShowValidateStatus{
		TmdbId:    1,
		Name:      "Test",
		GroupType: 1,
		Status:    0,
	}
	testutils.CheckResponseError(r, t, url, tvShowStatus, http.StatusUnprocessableEntity, "Status: value must be 1, 2, 3, 4 or 5")

	// TMDB ID / NAME / GROUP / STATUS
	type TvShowValidateAll struct {
		TmdbId    int    `json:"tmdb_id"`
		Name      string `json:"name"`
		GroupType int    `json:"group"`
		Status    int    `json:"status"`
	}
	tvShowAll := TvShowValidateAll{
		TmdbId:    0,
		Name:      "",
		GroupType: 0,
		Status:    6,
	}
	// sometimes the order of these error messages change
	// msg := "TmdbId: zero value, Name: less than min, GroupType: value must be 1, 2 or 3, Status: value must be 1, 2, 3, 4 or 5"
	msg := ""
	testutils.CheckResponseError(r, t, url, tvShowAll, http.StatusUnprocessableEntity, msg)
}

func TestTvShowCreateErrorAlreadyExist(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/tvshows/create"
	r.POST(url, controllers.TvShowCreate)

	// TMDB ID
	type TvShowValidateTmdbId struct {
		TmdbId    int    `json:"tmdb_id"`
		Name      string `json:"name"`
		GroupType int    `json:"group"`
		Status    int    `json:"status"`
	}
	tvShowExist := TvShowValidateTmdbId{
		TmdbId:    tvShowTest.TmdbId,
		Name:      "Test Show",
		GroupType: 1,
		Status:    1,
	}
	testutils.CheckResponseError(r, t, url, tvShowExist, http.StatusBadRequest, fmt.Sprintf("TvShow %s (TMDB ID %d) already exist", tvShowExist.Name, tvShowExist.TmdbId))
	// testutils.ListAllTvShows(DEBUG)
	testutils.CheckListAllTvShows(t, DEBUG, 1)
}

func TestTvShowCreateBatch(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/tvshows/create/batch"
	r.POST(url, controllers.TvShowCreateBatch)
	w := httptest.NewRecorder()

	tvShowCreate := models.TvShow{
		TmdbId:    TMDBID_THEROOKIE,
		Name:      "The Rookie",
		Overview:  "John Nolan, Lucy Chen",
		GroupType: 3,
		Status:    1,
	}
	var tvShowList []models.TvShow
	tvShowList = append(tvShowList, tvShowCreate)
	tvShowJson, err := json.Marshal(tvShowList)
	assert.Nil(t, err)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(tvShowJson)))
	assert.Nil(t, err)
	r.ServeHTTP(w, req)
	// println(w.Body.String(), w.Code)

	var tvShowListCreated []models.TvShow
	err = json.Unmarshal(w.Body.Bytes(), &tvShowListCreated)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, 2, tvShowListCreated[0].Id)
	tvShowCreate.Id = tvShowListCreated[0].Id
	testutils.CheckTvShow(t, tvShowCreate, tvShowListCreated[0])
}

func TestTvShowCreateBatchErrorBind(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/tvshows/create/batch"
	r.POST(url, controllers.TvShowCreateBatch)

	// ARRAY
	tvShowArray := models.TvShow{
		TmdbId: 1,
	}
	testutils.CheckResponseError(r, t, url, tvShowArray, http.StatusBadRequest, "json: cannot unmarshal object into Go value of type []models.TvShow")

	// TMDB ID
	type TvShowBindTmdbId struct {
		TmdbId string `json:"tmdb_id"`
	}
	tvShowTmdbId := TvShowBindTmdbId{
		TmdbId: "123",
	}
	testutils.CheckResponseError(r, t, url, []TvShowBindTmdbId{tvShowTmdbId}, http.StatusBadRequest, "json: cannot unmarshal string into Go struct field TvShow.tmdb_id of type int")

	// NAME
	type TvShowBindName struct {
		TmdbId int `json:"tmdb_id"`
		Name   int `json:"name"`
	}
	tvShowName := TvShowBindName{
		TmdbId: 123,
		Name:   0,
	}
	testutils.CheckResponseError(r, t, url, []TvShowBindName{tvShowName}, http.StatusBadRequest, "json: cannot unmarshal number into Go struct field TvShow.name of type string")

	// OVERVIEW
	type TvShowBindOverview struct {
		TmdbId   int    `json:"tmdb_id"`
		Name     string `json:"name"`
		Overview int    `json:"overview"`
	}
	tvShowOverview := TvShowBindOverview{
		TmdbId:   123,
		Name:     "Test",
		Overview: 2,
	}
	testutils.CheckResponseError(r, t, url, []TvShowBindOverview{tvShowOverview}, http.StatusBadRequest, "json: cannot unmarshal number into Go struct field TvShow.overview of type string")

	// GROUP
	type TvShowBindGroup struct {
		TmdbId    int    `json:"tmdb_id"`
		Name      string `json:"name"`
		Overview  string `json:"overview"`
		GroupType string `json:"group"`
	}
	tvShowGroup := TvShowBindGroup{
		TmdbId:    123,
		Name:      "Test",
		Overview:  "This is about",
		GroupType: "2",
	}
	testutils.CheckResponseError(r, t, url, []TvShowBindGroup{tvShowGroup}, http.StatusBadRequest, "json: cannot unmarshal string into Go struct field TvShow.group of type int")

	// STATUS
	type TvShowBindStatus struct {
		TmdbId    int    `json:"tmdb_id"`
		Name      string `json:"name"`
		Overview  string `json:"overview"`
		GroupType int    `json:"group"`
		Status    string `json:"status"`
	}
	tvShowStatus := TvShowBindStatus{
		TmdbId:    123,
		Name:      "Test",
		Overview:  "This is about",
		GroupType: 2,
		Status:    "0",
	}
	testutils.CheckResponseError(r, t, url, []TvShowBindStatus{tvShowStatus}, http.StatusBadRequest, "json: cannot unmarshal string into Go struct field TvShow.status of type int")
}

func TestTvShowCreateBatchErrorValidate(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/tvshows/create/batch"
	r.POST(url, controllers.TvShowCreateBatch)

	// TMDB ID
	type TvShowValidateTmdbId struct {
		TmdbId    int    `json:"tmdb_id"`
		Name      string `json:"name"`
		GroupType int    `json:"group"`
		Status    int    `json:"status"`
	}
	tvShowTmdbId := TvShowValidateTmdbId{
		TmdbId:    0,
		Name:      "Test Show",
		GroupType: 1,
		Status:    1,
	}
	testutils.CheckResponseError(r, t, url, []TvShowValidateTmdbId{tvShowTmdbId}, http.StatusUnprocessableEntity, "TmdbId: zero value")

	// NAME
	type TvShowValidateName struct {
		TmdbId    int    `json:"tmdb_id"`
		Name      string `json:"name"`
		GroupType int    `json:"group"`
		Status    int    `json:"status"`
	}
	tvShowName := TvShowValidateName{
		TmdbId:    1,
		Name:      "",
		GroupType: 1,
		Status:    1,
	}
	testutils.CheckResponseError(r, t, url, []TvShowValidateName{tvShowName}, http.StatusUnprocessableEntity, "Name: less than min")

	// GROUP
	type TvShowValidateGroup struct {
		TmdbId    int    `json:"tmdb_id"`
		Name      string `json:"name"`
		GroupType int    `json:"group"`
		Status    int    `json:"status"`
	}
	tvShowGroup := TvShowValidateGroup{
		TmdbId:    1,
		Name:      "Test",
		GroupType: 14,
		Status:    1,
	}
	testutils.CheckResponseError(r, t, url, []TvShowValidateGroup{tvShowGroup}, http.StatusUnprocessableEntity, "GroupType: value must be 1, 2 or 3")

	// STATUS
	type TvShowValidateStatus struct {
		TmdbId    int    `json:"tmdb_id"`
		Name      string `json:"name"`
		GroupType int    `json:"group"`
		Status    int    `json:"status"`
	}
	tvShowStatus := TvShowValidateStatus{
		TmdbId:    1,
		Name:      "Test",
		GroupType: 1,
		Status:    0,
	}
	testutils.CheckResponseError(r, t, url, []TvShowValidateStatus{tvShowStatus}, http.StatusUnprocessableEntity, "Status: value must be 1, 2, 3, 4 or 5")

	// TMDB ID / NAME / GROUP / STATUS
	type TvShowValidateAll struct {
		TmdbId    int    `json:"tmdb_id"`
		Name      string `json:"name"`
		GroupType int    `json:"group"`
		Status    int    `json:"status"`
	}
	tvShowAll := TvShowValidateAll{
		TmdbId:    0,
		Name:      "",
		GroupType: 0,
		Status:    6,
	}
	// sometimes the order of these error messages change
	// msg := "TmdbId: zero value, Name: less than min, GroupType: value must be 1, 2 or 3, Status: value must be 1, 2, 3, 4 or 5"
	msg := ""
	testutils.CheckResponseError(r, t, url, []TvShowValidateAll{tvShowAll}, http.StatusUnprocessableEntity, msg)
}

func TestTvShowCreateBatchErrorAlreadyExist(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/tvshows/create/batch"
	r.POST(url, controllers.TvShowCreateBatch)

	// TMDB ID
	type TvShowValidateTmdbId struct {
		TmdbId    int    `json:"tmdb_id"`
		Name      string `json:"name"`
		GroupType int    `json:"group"`
		Status    int    `json:"status"`
	}
	tvShowExist := TvShowValidateTmdbId{
		TmdbId:    tvShowTest.TmdbId,
		Name:      "Test Show",
		GroupType: 1,
		Status:    1,
	}
	testutils.CheckResponseError(r, t, url, []TvShowValidateTmdbId{tvShowExist}, http.StatusBadRequest, fmt.Sprintf("TvShow %s (TMDB ID %d) already exist", tvShowExist.Name, tvShowExist.TmdbId))
	// testutils.ListAllTvShows(DEBUG)
	testutils.CheckListAllTvShows(t, DEBUG, 2)
}

func TestTvShowGetById(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/tvshows/:id"
	r.GET(url, controllers.TvShowListById)
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, strings.Replace(url, ":id", strconv.Itoa(tvShowTest.Id), 1), nil)
	assert.Nil(t, err)
	r.ServeHTTP(w, req)
	// println(w.Body.String())

	var tvShowListed models.TvShow
	err = json.Unmarshal(w.Body.Bytes(), &tvShowListed)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	// tvShowTest.Overview = "tbbt"
	testutils.CheckTvShow(t, tvShowListed, tvShowTest)
}

func TestTvShowEdit(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/tvshows/:id"
	r.PUT(url, controllers.TvShowEdit)
	w := httptest.NewRecorder()

	overview := "Test Updated Overview Edit"
	tvShowEdit := models.TvShow{
		Id:        tvShowTest.Id,
		TmdbId:    tvShowTest.TmdbId,
		Name:      tvShowTest.Name,
		Overview:  overview,
		GroupType: tvShowTest.GroupType,
		Status:    tvShowTest.Status,
	}
	tvShowJson, err := json.Marshal(tvShowEdit)
	assert.Nil(t, err)
	req, err := http.NewRequest(http.MethodPut, strings.Replace(url, ":id", strconv.Itoa(tvShowTest.Id), 1), strings.NewReader(string(tvShowJson)))
	assert.Nil(t, err)
	r.ServeHTTP(w, req)
	// println(w.Body.String())

	var tvShowUpdated models.TvShow
	err = json.Unmarshal(w.Body.Bytes(), &tvShowUpdated)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 1, tvShowUpdated.Id)
	tvShowTest.Overview = overview
	testutils.CheckTvShow(t, tvShowUpdated, tvShowTest)
}

// TODO: tv show edit bind json

// Episode

func TestEpisodeCreate(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/episodes/create"
	r.POST(url, controllers.EpisodeCreate)
	w := httptest.NewRecorder()

	episodeToCreate := models.Episode{
		TmdbId:  tvShowTest.TmdbId,
		Season:  1,
		Episode: 1,
		Name:    "Pilot",
		AirDate: 20090309,
		Watched: false,
	}
	episodeJson, err := json.Marshal(episodeToCreate)
	assert.Nil(t, err)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(episodeJson)))
	assert.Nil(t, err)
	r.ServeHTTP(w, req)
	// println(w.Body.String())

	var episodeCreated models.Episode
	err = json.Unmarshal(w.Body.Bytes(), &episodeCreated)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, 1, episodeCreated.Id)
	episodeToCreate.Id = episodeCreated.Id
	testutils.CheckEpisode(t, episodeCreated, episodeToCreate)

	episodesTest = append(episodesTest, episodeCreated)
}

func TestEpisodeCreateErrorBind(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/episodes/create"
	r.POST(url, controllers.EpisodeCreate)

	// TMDB ID
	type EpisodeBindTmdbId struct {
		TmdbId string `json:"tmdb_id"`
	}
	episodeTmdbId := EpisodeBindTmdbId{
		TmdbId: "123",
	}
	testutils.CheckResponseError(r, t, url, episodeTmdbId, http.StatusBadRequest, "json: cannot unmarshal string into Go struct field Episode.tmdb_id of type int")

	// SEASON
	type EpisodeBindSeason struct {
		TmdbId int    `json:"tmdb_id"`
		Season string `json:"season"`
	}
	episodeSeason := EpisodeBindSeason{
		TmdbId: 123,
		Season: "2",
	}
	testutils.CheckResponseError(r, t, url, episodeSeason, http.StatusBadRequest, "json: cannot unmarshal string into Go struct field Episode.season of type int")

	// EPISODE
	type EpisodeBindEpisode struct {
		TmdbId  int    `json:"tmdb_id"`
		Season  int    `json:"season"`
		Episode string `json:"episode"`
	}
	episodeEpisode := EpisodeBindEpisode{
		TmdbId:  123,
		Season:  2,
		Episode: "1",
	}
	testutils.CheckResponseError(r, t, url, episodeEpisode, http.StatusBadRequest, "json: cannot unmarshal string into Go struct field Episode.episode of type int")

	// NAME
	type EpisodeBindName struct {
		TmdbId  int `json:"tmdb_id"`
		Season  int `json:"season"`
		Episode int `json:"episode"`
		Name    int `json:"name"`
	}
	episodeName := EpisodeBindName{
		TmdbId:  123,
		Season:  1,
		Episode: 1,
		Name:    0,
	}
	testutils.CheckResponseError(r, t, url, episodeName, http.StatusBadRequest, "json: cannot unmarshal number into Go struct field Episode.name of type string")

	// OVERVIEW
	type EpisodeBindOverview struct {
		TmdbId   int    `json:"tmdb_id"`
		Season   int    `json:"season"`
		Episode  int    `json:"episode"`
		Name     string `json:"name"`
		Overview int    `json:"overview"`
	}
	episodeOverview := EpisodeBindOverview{
		TmdbId:   123,
		Season:   1,
		Episode:  1,
		Name:     "Test",
		Overview: 2,
	}
	testutils.CheckResponseError(r, t, url, episodeOverview, http.StatusBadRequest, "json: cannot unmarshal number into Go struct field Episode.overview of type string")

	// AIR DATE
	type EpisodeBindAirDate struct {
		TmdbId   int    `json:"tmdb_id"`
		Season   int    `json:"season"`
		Episode  int    `json:"episode"`
		Name     string `json:"name"`
		Overview string `json:"overview"`
		AirDate  string `json:"air_date"`
	}
	episodeAirDate := EpisodeBindAirDate{
		TmdbId:   123,
		Season:   1,
		Episode:  1,
		Name:     "Test",
		Overview: "About",
		AirDate:  "20250101",
	}
	testutils.CheckResponseError(r, t, url, episodeAirDate, http.StatusBadRequest, "json: cannot unmarshal string into Go struct field Episode.air_date of type int")

	// WATCHED
	type EpisodeBindWatched struct {
		TmdbId   int    `json:"tmdb_id"`
		Season   int    `json:"season"`
		Episode  int    `json:"episode"`
		Name     string `json:"name"`
		Overview string `json:"overview"`
		AirDate  int    `json:"air_date"`
		Watched  int    `json:"watched"`
	}
	episodeWatched := EpisodeBindWatched{
		TmdbId:   123,
		Season:   1,
		Episode:  1,
		Name:     "Test",
		Overview: "About",
		AirDate:  20250101,
		Watched:  0,
	}
	testutils.CheckResponseError(r, t, url, episodeWatched, http.StatusBadRequest, "json: cannot unmarshal number into Go struct field Episode.watched of type bool")

	// WATCHED DATE
	type EpisodeBindWatchedDate struct {
		TmdbId      int    `json:"tmdb_id"`
		Season      int    `json:"season"`
		Episode     int    `json:"episode"`
		Name        string `json:"name"`
		Overview    string `json:"overview"`
		AirDate     int    `json:"air_date"`
		Watched     bool   `json:"watched"`
		WatchedDate string `json:"watched_date"`
	}
	episodeWatchedDate := EpisodeBindWatchedDate{
		TmdbId:      123,
		Season:      1,
		Episode:     1,
		Name:        "Test",
		Overview:    "About",
		AirDate:     20250101,
		Watched:     false,
		WatchedDate: "20241231",
	}
	testutils.CheckResponseError(r, t, url, episodeWatchedDate, http.StatusBadRequest, "json: cannot unmarshal string into Go struct field Episode.watched_date of type int")
}

func TestEpisodeCreateErrorValidade(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/episodes/create"
	r.POST(url, controllers.EpisodeCreate)

	// TMDB ID
	type EpisodeValidadeTmdbId struct {
		TmdbId  int    `json:"tmdb_id"`
		Season  int    `json:"season"`
		Episode int    `json:"episode"`
		Name    string `json:"name"`
	}
	episodeTmdbId := EpisodeValidadeTmdbId{
		TmdbId:  0,
		Season:  1,
		Episode: 1,
		Name:    "Test",
	}
	testutils.CheckResponseError(r, t, url, episodeTmdbId, http.StatusUnprocessableEntity, "TmdbId: zero value")

	// SEASON
	type EpisodeValidadeSeason struct {
		TmdbId  int    `json:"tmdb_id"`
		Season  int    `json:"season"`
		Episode int    `json:"episode"`
		Name    string `json:"name"`
	}
	episodeSeason := EpisodeValidadeSeason{
		TmdbId:  10,
		Season:  0,
		Episode: 1,
		Name:    "Test",
	}
	testutils.CheckResponseError(r, t, url, episodeSeason, http.StatusUnprocessableEntity, "Season: zero value")

	// EPISODE
	type EpisodeValidadeEpisode struct {
		TmdbId  int    `json:"tmdb_id"`
		Season  int    `json:"season"`
		Episode int    `json:"episode"`
		Name    string `json:"name"`
	}
	episodeEpisode := EpisodeValidadeEpisode{
		TmdbId:  1,
		Season:  1,
		Episode: 0,
		Name:    "Test",
	}
	testutils.CheckResponseError(r, t, url, episodeEpisode, http.StatusUnprocessableEntity, "Episode: zero value")

	// NAME
	type EpisodeValidadeName struct {
		TmdbId  int    `json:"tmdb_id"`
		Season  int    `json:"season"`
		Episode int    `json:"episode"`
		Name    string `json:"name"`
	}
	episodeName := EpisodeValidadeName{
		TmdbId:  100,
		Season:  1,
		Episode: 1,
		Name:    "",
	}
	testutils.CheckResponseError(r, t, url, episodeName, http.StatusUnprocessableEntity, "Name: less than min")

	// AIR DATE - len != 8
	type EpisodeValidadeAirDateLen struct {
		TmdbId  int    `json:"tmdb_id"`
		Season  int    `json:"season"`
		Episode int    `json:"episode"`
		Name    string `json:"name"`
		AirDate int    `json:"air_date"`
	}
	episodeAirDateLen := EpisodeValidadeAirDateLen{
		TmdbId:  100,
		Season:  1,
		Episode: 1,
		Name:    "Test",
		AirDate: 202501,
	}
	testutils.CheckResponseError(r, t, url, episodeAirDateLen, http.StatusUnprocessableEntity, "AirDate: date must be YYYYMMDD")

	// AIR DATE - value < 0
	type EpisodeValidadeAirDateValueNegative struct {
		TmdbId  int    `json:"tmdb_id"`
		Season  int    `json:"season"`
		Episode int    `json:"episode"`
		Name    string `json:"name"`
		AirDate int    `json:"air_date"`
	}
	episodeAirDateValueNegative := EpisodeValidadeAirDateValueNegative{
		TmdbId:  100,
		Season:  1,
		Episode: 1,
		Name:    "Test",
		AirDate: -1,
	}
	testutils.CheckResponseError(r, t, url, episodeAirDateValueNegative, http.StatusUnprocessableEntity, "AirDate: date must be YYYYMMDD")

	// AIR DATE - invalid date (invalid month)
	type EpisodeValidadeAirDateInvalidDateInvalidMonth struct {
		TmdbId  int    `json:"tmdb_id"`
		Season  int    `json:"season"`
		Episode int    `json:"episode"`
		Name    string `json:"name"`
		AirDate int    `json:"air_date"`
	}
	episodeAirDateInvalidDateInvalidMonth := EpisodeValidadeAirDateInvalidDateInvalidMonth{
		TmdbId:  100,
		Season:  1,
		Episode: 1,
		Name:    "Test",
		AirDate: 20251301,
	}
	testutils.CheckResponseError(r, t, url, episodeAirDateInvalidDateInvalidMonth, http.StatusUnprocessableEntity, fmt.Sprintf("AirDate: parsing time \"%d\": month out of range", episodeAirDateInvalidDateInvalidMonth.AirDate))

	// AIR DATE - invalid date (invalid day)
	type EpisodeValidadeAirDateInvalidDateInvalidDay struct {
		TmdbId  int    `json:"tmdb_id"`
		Season  int    `json:"season"`
		Episode int    `json:"episode"`
		Name    string `json:"name"`
		AirDate int    `json:"air_date"`
	}
	episodeAirDateInvalidDateInvalidDay := EpisodeValidadeAirDateInvalidDateInvalidDay{
		TmdbId:  100,
		Season:  1,
		Episode: 1,
		Name:    "Test",
		AirDate: 20250431,
	}
	testutils.CheckResponseError(r, t, url, episodeAirDateInvalidDateInvalidDay, http.StatusUnprocessableEntity, fmt.Sprintf("AirDate: parsing time \"%d\": day out of range", episodeAirDateInvalidDateInvalidDay.AirDate))

	// WATCHED DATE - len != 8
	type EpisodeValidadeWatchedDateLen struct {
		TmdbId      int    `json:"tmdb_id"`
		Season      int    `json:"season"`
		Episode     int    `json:"episode"`
		Name        string `json:"name"`
		WatchedDate int    `json:"watched_date"`
	}
	episodeWatchedDateLen := EpisodeValidadeWatchedDateLen{
		TmdbId:      100,
		Season:      1,
		Episode:     1,
		Name:        "Test",
		WatchedDate: 2025102,
	}
	testutils.CheckResponseError(r, t, url, episodeWatchedDateLen, http.StatusUnprocessableEntity, "WatchedDate: date must be YYYYMMDD")

	// WATCHED DATE - value < 0
	type EpisodeValidadeWatchedDateValueNegative struct {
		TmdbId      int    `json:"tmdb_id"`
		Season      int    `json:"season"`
		Episode     int    `json:"episode"`
		Name        string `json:"name"`
		WatchedDate int    `json:"watched_date"`
	}
	episodeWatchedDateValueNegative := EpisodeValidadeWatchedDateValueNegative{
		TmdbId:      100,
		Season:      1,
		Episode:     1,
		Name:        "Test",
		WatchedDate: -2025,
	}
	testutils.CheckResponseError(r, t, url, episodeWatchedDateValueNegative, http.StatusUnprocessableEntity, "WatchedDate: date must be YYYYMMDD")

	// WATCHED DATE - invalid date (invalid month)
	type EpisodeValidadeWatchedDateInvalidDateInvalidMonth struct {
		TmdbId      int    `json:"tmdb_id"`
		Season      int    `json:"season"`
		Episode     int    `json:"episode"`
		Name        string `json:"name"`
		WatchedDate int    `json:"watched_date"`
	}
	episodeWatchedDateInvalidDateInvalidMonth := EpisodeValidadeWatchedDateInvalidDateInvalidMonth{
		TmdbId:      100,
		Season:      1,
		Episode:     1,
		Name:        "Test",
		WatchedDate: 20250001,
	}
	testutils.CheckResponseError(r, t, url, episodeWatchedDateInvalidDateInvalidMonth, http.StatusUnprocessableEntity, fmt.Sprintf("WatchedDate: parsing time \"%d\": month out of range", episodeWatchedDateInvalidDateInvalidMonth.WatchedDate))

	// WATCHED DATE - invalid date (invalid day)
	type EpisodeValidadeWatchedDateInvalidDateInvalidDay struct {
		TmdbId      int    `json:"tmdb_id"`
		Season      int    `json:"season"`
		Episode     int    `json:"episode"`
		Name        string `json:"name"`
		WatchedDate int    `json:"watched_date"`
	}
	episodeWatchedDateInvalidDateInvalidDay := EpisodeValidadeWatchedDateInvalidDateInvalidDay{
		TmdbId:      100,
		Season:      1,
		Episode:     1,
		Name:        "Test",
		WatchedDate: 20250230,
	}
	testutils.CheckResponseError(r, t, url, episodeWatchedDateInvalidDateInvalidDay, http.StatusUnprocessableEntity, fmt.Sprintf("WatchedDate: parsing time \"%d\": day out of range", episodeWatchedDateInvalidDateInvalidDay.WatchedDate))
}

func TestEpisodeCreateErrorTvShowNotFound(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/episodes/create"
	r.POST(url, controllers.EpisodeCreate)

	episodeToCreate := models.Episode{
		TmdbId:  421412421,
		Season:  1,
		Episode: 1,
		Name:    "Test",
	}

	testutils.CheckResponseError(r, t, url, episodeToCreate, http.StatusNotFound, "TvShow not found")
}

func TestEpisodeCreateErrorAlreadyExist(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/episodes/create"
	r.POST(url, controllers.EpisodeCreate)

	episodeToCreate := models.Episode{
		TmdbId:  tvShowTest.TmdbId,
		Season:  1,
		Episode: 1,
		Name:    "Test",
	}

	testutils.CheckResponseError(r, t, url, episodeToCreate, http.StatusBadRequest, fmt.Sprintf("episode %dx%02d already exist for %s", episodeToCreate.Season, episodeToCreate.Episode, tvShowTest.Name))
	// testutils.ListAllEpisodes(DEBUG)
	testutils.CheckListAllEpisodes(t, DEBUG, 1)
}

func TestEpisodeCreateBatch(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/episodes/create/batch"
	r.POST(url, controllers.EpisodeCreateBatch)
	w := httptest.NewRecorder()

	// Castle
	episodeToCreate1 := models.Episode{TmdbId: TMDBID_CASTLE, Season: 1, Episode: 2, Name: "Nanny McDead", AirDate: 20090316, Watched: false}
	episodeToCreate2 := models.Episode{TmdbId: TMDBID_CASTLE, Season: 2, Episode: 1, Name: "Deep in Death", AirDate: 20090921, Watched: false}
	// The Rookie
	episodeToCreate3 := models.Episode{TmdbId: TMDBID_THEROOKIE, Season: 1, Episode: 1, Name: "Pilot", AirDate: 20181016, Watched: false}

	var episodesToCreate []models.Episode
	episodesToCreate = append(episodesToCreate, episodeToCreate1)
	episodesToCreate = append(episodesToCreate, episodeToCreate2)
	episodesToCreate = append(episodesToCreate, episodeToCreate3)
	episodeJson, err := json.Marshal(episodesToCreate)
	assert.Nil(t, err)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(episodeJson)))
	assert.Nil(t, err)
	r.ServeHTTP(w, req)
	// println(w.Body.String(), w.Code)

	var episodesCreated []models.Episode
	err = json.Unmarshal(w.Body.Bytes(), &episodesCreated)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, w.Code)

	for index, episode := range episodesCreated {
		// fmt.Println(index, episode)
		if index == 0 {
			episodeToCreate1.Id = 2
			testutils.CheckEpisode(t, episodeToCreate1, episode)
		} else if index == 1 {
			episodeToCreate2.Id = 3
			testutils.CheckEpisode(t, episodeToCreate2, episode)
		} else if index == 2 {
			episodeToCreate3.Id = 4
			testutils.CheckEpisode(t, episodeToCreate3, episode)
		}
	}

	episodesTest = append(episodesTest, episodesCreated...)
	assert.Equal(t, 4, len(episodesTest))
}

func TestEpisodeCreateBatchErrorBind(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/episodes/create/batch"
	r.POST(url, controllers.EpisodeCreateBatch)

	// ARRAY
	episodeArray := models.Episode{
		TmdbId: 1,
	}
	testutils.CheckResponseError(r, t, url, episodeArray, http.StatusBadRequest, "json: cannot unmarshal object into Go value of type []models.Episode")

	// TMDB ID
	type EpisodeBindTmdbId struct {
		TmdbId string `json:"tmdb_id"`
	}
	episodeTmdbId := EpisodeBindTmdbId{
		TmdbId: "123",
	}
	testutils.CheckResponseError(r, t, url, []EpisodeBindTmdbId{episodeTmdbId}, http.StatusBadRequest, "json: cannot unmarshal string into Go struct field Episode.tmdb_id of type int")

	// SEASON
	type EpisodeBindSeason struct {
		TmdbId int    `json:"tmdb_id"`
		Season string `json:"season"`
	}
	episodeSeason := EpisodeBindSeason{
		TmdbId: 123,
		Season: "2",
	}
	testutils.CheckResponseError(r, t, url, []EpisodeBindSeason{episodeSeason}, http.StatusBadRequest, "json: cannot unmarshal string into Go struct field Episode.season of type int")

	// EPISODE
	type EpisodeBindEpisode struct {
		TmdbId  int    `json:"tmdb_id"`
		Season  int    `json:"season"`
		Episode string `json:"episode"`
	}
	episodeEpisode := EpisodeBindEpisode{
		TmdbId:  123,
		Season:  2,
		Episode: "1",
	}
	testutils.CheckResponseError(r, t, url, []EpisodeBindEpisode{episodeEpisode}, http.StatusBadRequest, "json: cannot unmarshal string into Go struct field Episode.episode of type int")

	// NAME
	type EpisodeBindName struct {
		TmdbId  int `json:"tmdb_id"`
		Season  int `json:"season"`
		Episode int `json:"episode"`
		Name    int `json:"name"`
	}
	episodeName := EpisodeBindName{
		TmdbId:  123,
		Season:  1,
		Episode: 1,
		Name:    0,
	}
	testutils.CheckResponseError(r, t, url, []EpisodeBindName{episodeName}, http.StatusBadRequest, "json: cannot unmarshal number into Go struct field Episode.name of type string")

	// OVERVIEW
	type EpisodeBindOverview struct {
		TmdbId   int    `json:"tmdb_id"`
		Season   int    `json:"season"`
		Episode  int    `json:"episode"`
		Name     string `json:"name"`
		Overview int    `json:"overview"`
	}
	episodeOverview := EpisodeBindOverview{
		TmdbId:   123,
		Season:   1,
		Episode:  1,
		Name:     "Test",
		Overview: 2,
	}
	testutils.CheckResponseError(r, t, url, []EpisodeBindOverview{episodeOverview}, http.StatusBadRequest, "json: cannot unmarshal number into Go struct field Episode.overview of type string")

	// AIR DATE
	type EpisodeBindAirDate struct {
		TmdbId   int    `json:"tmdb_id"`
		Season   int    `json:"season"`
		Episode  int    `json:"episode"`
		Name     string `json:"name"`
		Overview string `json:"overview"`
		AirDate  string `json:"air_date"`
	}
	episodeAirDate := EpisodeBindAirDate{
		TmdbId:   123,
		Season:   1,
		Episode:  1,
		Name:     "Test",
		Overview: "About",
		AirDate:  "20250101",
	}
	testutils.CheckResponseError(r, t, url, []EpisodeBindAirDate{episodeAirDate}, http.StatusBadRequest, "json: cannot unmarshal string into Go struct field Episode.air_date of type int")

	// WATCHED
	type EpisodeBindWatched struct {
		TmdbId   int    `json:"tmdb_id"`
		Season   int    `json:"season"`
		Episode  int    `json:"episode"`
		Name     string `json:"name"`
		Overview string `json:"overview"`
		AirDate  int    `json:"air_date"`
		Watched  int    `json:"watched"`
	}
	episodeWatched := EpisodeBindWatched{
		TmdbId:   123,
		Season:   1,
		Episode:  1,
		Name:     "Test",
		Overview: "About",
		AirDate:  20250101,
		Watched:  0,
	}
	testutils.CheckResponseError(r, t, url, []EpisodeBindWatched{episodeWatched}, http.StatusBadRequest, "json: cannot unmarshal number into Go struct field Episode.watched of type bool")

	// WATCHED DATE
	type EpisodeBindWatchedDate struct {
		TmdbId      int    `json:"tmdb_id"`
		Season      int    `json:"season"`
		Episode     int    `json:"episode"`
		Name        string `json:"name"`
		Overview    string `json:"overview"`
		AirDate     int    `json:"air_date"`
		Watched     bool   `json:"watched"`
		WatchedDate string `json:"watched_date"`
	}
	episodeWatchedDate := EpisodeBindWatchedDate{
		TmdbId:      123,
		Season:      1,
		Episode:     1,
		Name:        "Test",
		Overview:    "About",
		AirDate:     20250101,
		Watched:     false,
		WatchedDate: "20241231",
	}
	testutils.CheckResponseError(r, t, url, []EpisodeBindWatchedDate{episodeWatchedDate}, http.StatusBadRequest, "json: cannot unmarshal string into Go struct field Episode.watched_date of type int")
}

func TestEpisodeCreateBatchErrorValidate(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/episodes/create/batch"
	r.POST(url, controllers.EpisodeCreateBatch)

	// TMDB ID
	type EpisodeValidadeTmdbId struct {
		TmdbId  int    `json:"tmdb_id"`
		Season  int    `json:"season"`
		Episode int    `json:"episode"`
		Name    string `json:"name"`
	}
	episodeTmdbId := EpisodeValidadeTmdbId{
		TmdbId:  0,
		Season:  1,
		Episode: 1,
		Name:    "Test",
	}
	testutils.CheckResponseError(r, t, url, []EpisodeValidadeTmdbId{episodeTmdbId}, http.StatusUnprocessableEntity, "TmdbId: zero value")

	// SEASON
	type EpisodeValidadeSeason struct {
		TmdbId  int    `json:"tmdb_id"`
		Season  int    `json:"season"`
		Episode int    `json:"episode"`
		Name    string `json:"name"`
	}
	episodeSeason := EpisodeValidadeSeason{
		TmdbId:  10,
		Season:  0,
		Episode: 1,
		Name:    "Test",
	}
	testutils.CheckResponseError(r, t, url, []EpisodeValidadeSeason{episodeSeason}, http.StatusUnprocessableEntity, "Season: zero value")

	// EPISODE
	type EpisodeValidadeEpisode struct {
		TmdbId  int    `json:"tmdb_id"`
		Season  int    `json:"season"`
		Episode int    `json:"episode"`
		Name    string `json:"name"`
	}
	episodeEpisode := EpisodeValidadeEpisode{
		TmdbId:  1,
		Season:  1,
		Episode: 0,
		Name:    "Test",
	}
	testutils.CheckResponseError(r, t, url, []EpisodeValidadeEpisode{episodeEpisode}, http.StatusUnprocessableEntity, "Episode: zero value")

	// NAME
	type EpisodeValidadeName struct {
		TmdbId  int    `json:"tmdb_id"`
		Season  int    `json:"season"`
		Episode int    `json:"episode"`
		Name    string `json:"name"`
	}
	episodeName := EpisodeValidadeName{
		TmdbId:  100,
		Season:  1,
		Episode: 1,
		Name:    "",
	}
	testutils.CheckResponseError(r, t, url, []EpisodeValidadeName{episodeName}, http.StatusUnprocessableEntity, "Name: less than min")

	// AIR DATE - len != 8
	type EpisodeValidadeAirDateLen struct {
		TmdbId  int    `json:"tmdb_id"`
		Season  int    `json:"season"`
		Episode int    `json:"episode"`
		Name    string `json:"name"`
		AirDate int    `json:"air_date"`
	}
	episodeAirDateLen := EpisodeValidadeAirDateLen{
		TmdbId:  100,
		Season:  1,
		Episode: 1,
		Name:    "Test",
		AirDate: 202501,
	}
	testutils.CheckResponseError(r, t, url, []EpisodeValidadeAirDateLen{episodeAirDateLen}, http.StatusUnprocessableEntity, "AirDate: date must be YYYYMMDD")

	// AIR DATE - value < 0
	type EpisodeValidadeAirDateValueNegative struct {
		TmdbId  int    `json:"tmdb_id"`
		Season  int    `json:"season"`
		Episode int    `json:"episode"`
		Name    string `json:"name"`
		AirDate int    `json:"air_date"`
	}
	episodeAirDateValueNegative := EpisodeValidadeAirDateValueNegative{
		TmdbId:  100,
		Season:  1,
		Episode: 1,
		Name:    "Test",
		AirDate: -1,
	}
	testutils.CheckResponseError(r, t, url, []EpisodeValidadeAirDateValueNegative{episodeAirDateValueNegative}, http.StatusUnprocessableEntity, "AirDate: date must be YYYYMMDD")

	// AIR DATE - invalid date (invalid month)
	type EpisodeValidadeAirDateInvalidDateInvalidMonth struct {
		TmdbId  int    `json:"tmdb_id"`
		Season  int    `json:"season"`
		Episode int    `json:"episode"`
		Name    string `json:"name"`
		AirDate int    `json:"air_date"`
	}
	episodeAirDateInvalidDateInvalidMonth := EpisodeValidadeAirDateInvalidDateInvalidMonth{
		TmdbId:  100,
		Season:  1,
		Episode: 1,
		Name:    "Test",
		AirDate: 20251301,
	}
	testutils.CheckResponseError(r, t, url, []EpisodeValidadeAirDateInvalidDateInvalidMonth{episodeAirDateInvalidDateInvalidMonth}, http.StatusUnprocessableEntity, fmt.Sprintf("AirDate: parsing time \"%d\": month out of range", episodeAirDateInvalidDateInvalidMonth.AirDate))

	// AIR DATE - invalid date (invalid day)
	type EpisodeValidadeAirDateInvalidDateInvalidDay struct {
		TmdbId  int    `json:"tmdb_id"`
		Season  int    `json:"season"`
		Episode int    `json:"episode"`
		Name    string `json:"name"`
		AirDate int    `json:"air_date"`
	}
	episodeAirDateInvalidDateInvalidDay := EpisodeValidadeAirDateInvalidDateInvalidDay{
		TmdbId:  100,
		Season:  1,
		Episode: 1,
		Name:    "Test",
		AirDate: 20250431,
	}
	testutils.CheckResponseError(r, t, url, []EpisodeValidadeAirDateInvalidDateInvalidDay{episodeAirDateInvalidDateInvalidDay}, http.StatusUnprocessableEntity, fmt.Sprintf("AirDate: parsing time \"%d\": day out of range", episodeAirDateInvalidDateInvalidDay.AirDate))

	// WATCHED DATE - len != 8
	type EpisodeValidadeWatchedDateLen struct {
		TmdbId      int    `json:"tmdb_id"`
		Season      int    `json:"season"`
		Episode     int    `json:"episode"`
		Name        string `json:"name"`
		WatchedDate int    `json:"watched_date"`
	}
	episodeWatchedDateLen := EpisodeValidadeWatchedDateLen{
		TmdbId:      100,
		Season:      1,
		Episode:     1,
		Name:        "Test",
		WatchedDate: 2025102,
	}
	testutils.CheckResponseError(r, t, url, []EpisodeValidadeWatchedDateLen{episodeWatchedDateLen}, http.StatusUnprocessableEntity, "WatchedDate: date must be YYYYMMDD")

	// WATCHED DATE - value < 0
	type EpisodeValidadeWatchedDateValueNegative struct {
		TmdbId      int    `json:"tmdb_id"`
		Season      int    `json:"season"`
		Episode     int    `json:"episode"`
		Name        string `json:"name"`
		WatchedDate int    `json:"watched_date"`
	}
	episodeWatchedDateValueNegative := EpisodeValidadeWatchedDateValueNegative{
		TmdbId:      100,
		Season:      1,
		Episode:     1,
		Name:        "Test",
		WatchedDate: -2025,
	}
	testutils.CheckResponseError(r, t, url, []EpisodeValidadeWatchedDateValueNegative{episodeWatchedDateValueNegative}, http.StatusUnprocessableEntity, "WatchedDate: date must be YYYYMMDD")

	// WATCHED DATE - invalid date (invalid month)
	type EpisodeValidadeWatchedDateInvalidDateInvalidMonth struct {
		TmdbId      int    `json:"tmdb_id"`
		Season      int    `json:"season"`
		Episode     int    `json:"episode"`
		Name        string `json:"name"`
		WatchedDate int    `json:"watched_date"`
	}
	episodeWatchedDateInvalidDateInvalidMonth := EpisodeValidadeWatchedDateInvalidDateInvalidMonth{
		TmdbId:      100,
		Season:      1,
		Episode:     1,
		Name:        "Test",
		WatchedDate: 20250001,
	}
	testutils.CheckResponseError(r, t, url, []EpisodeValidadeWatchedDateInvalidDateInvalidMonth{episodeWatchedDateInvalidDateInvalidMonth}, http.StatusUnprocessableEntity, fmt.Sprintf("WatchedDate: parsing time \"%d\": month out of range", episodeWatchedDateInvalidDateInvalidMonth.WatchedDate))

	// WATCHED DATE - invalid date (invalid day)
	type EpisodeValidadeWatchedDateInvalidDateInvalidDay struct {
		TmdbId      int    `json:"tmdb_id"`
		Season      int    `json:"season"`
		Episode     int    `json:"episode"`
		Name        string `json:"name"`
		WatchedDate int    `json:"watched_date"`
	}
	episodeWatchedDateInvalidDateInvalidDay := EpisodeValidadeWatchedDateInvalidDateInvalidDay{
		TmdbId:      100,
		Season:      1,
		Episode:     1,
		Name:        "Test",
		WatchedDate: 20250230,
	}
	testutils.CheckResponseError(r, t, url, []EpisodeValidadeWatchedDateInvalidDateInvalidDay{episodeWatchedDateInvalidDateInvalidDay}, http.StatusUnprocessableEntity, fmt.Sprintf("WatchedDate: parsing time \"%d\": day out of range", episodeWatchedDateInvalidDateInvalidDay.WatchedDate))
}

func TestEpisodeCreateBatchErrorTvShowNotFound(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/episodes/create/batch"
	r.POST(url, controllers.EpisodeCreateBatch)

	episodeToCreate := models.Episode{
		TmdbId:  421412421,
		Season:  1,
		Episode: 1,
		Name:    "Test",
	}

	testutils.CheckResponseError(r, t, url, []models.Episode{episodeToCreate}, http.StatusNotFound, "TvShow not found")
}

func TestEpisodeCreateBatchErrorAlreadyExist(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/episodes/create/batch"
	r.POST(url, controllers.EpisodeCreateBatch)

	episodeToCreate := models.Episode{
		TmdbId:  tvShowTest.TmdbId,
		Season:  1,
		Episode: 1,
		Name:    "Test",
	}

	testutils.CheckResponseError(r, t, url, []models.Episode{episodeToCreate}, http.StatusBadRequest, fmt.Sprintf("episode %dx%02d already exist for %s", episodeToCreate.Season, episodeToCreate.Episode, tvShowTest.Name))
	// testutils.ListAllEpisodes(DEBUG)
	testutils.CheckListAllEpisodes(t, DEBUG, 4)
}

func TestEpisodeEdit(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/episodes/edit/:id"
	r.PUT(url, controllers.EpisodeEdit)
	w := httptest.NewRecorder()

	id := 4
	episodeToUpdate, err := GetEpisodeTest(id)
	assert.Nil(t, err)
	episodeToUpdate.Overview = "The Rookie overview test update"
	// episodeToUpdate.Dump()

	episodeJson, err := json.Marshal(episodeToUpdate)
	assert.Nil(t, err)
	req, err := http.NewRequest(http.MethodPut, strings.Replace(url, ":id", strconv.Itoa(tvShowTest.Id), 1), strings.NewReader(string(episodeJson)))
	assert.Nil(t, err)
	r.ServeHTTP(w, req)
	// println(w.Body.String())

	var episodeUpdated models.Episode
	err = json.Unmarshal(w.Body.Bytes(), &episodeUpdated)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, id, episodeUpdated.Id)
	testutils.CheckEpisode(t, episodeUpdated, episodeToUpdate)

	err = UpdateEpisodeTest(episodeUpdated)
	assert.Nil(t, err)
}

// TODO: episode edit bind json

func TestEpisodeMarkAsWatched(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/episodes/watched/:id"
	r.PUT(url, controllers.EpisodeEditMarkWatched)
	w := httptest.NewRecorder()

	// DumpEpisodeTest()

	id := 4
	episodeToWatch, err := GetEpisodeTest(id)
	assert.Nil(t, err)
	// episodeToWatch.Dump()

	req, err := http.NewRequest(http.MethodPut, strings.Replace(url, ":id", strconv.Itoa(episodeToWatch.Id), 1), nil)
	assert.Nil(t, err)
	r.ServeHTTP(w, req)
	// println(w.Body.String())

	var episodeWatched models.Episode
	err = json.Unmarshal(w.Body.Bytes(), &episodeWatched)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	episodeToWatch.Watched = true
	episodeToWatch.WatchedDate = generic.GetCurrentDate()
	testutils.CheckEpisode(t, episodeWatched, episodeToWatch)

	err = UpdateEpisodeTest(episodeWatched)
	assert.Nil(t, err)
}

func TestEpisodeMarkSeasonAsWatched(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/episodes/watched/season/:tmdbid/:season"
	r.PUT(url, controllers.EpisodeEditMarkWatched)
	w := httptest.NewRecorder()

	episodeToWatch, err := GetEpisodeTest(1)
	assert.Nil(t, err)
	// episodeToWatch.Dump()

	req, err := http.NewRequest(http.MethodPut, strings.Replace(strings.Replace(url, ":season", strconv.Itoa(episodeToWatch.Id), 1), ":tmdbid", strconv.Itoa(episodeToWatch.TmdbId), 1), nil)
	assert.Nil(t, err)
	r.ServeHTTP(w, req)
	// println(w.Body.String())

	var episodesSeasonWatched []models.Episode
	err = json.Unmarshal(w.Body.Bytes(), &episodesSeasonWatched)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	for _, ep := range episodesSeasonWatched {
		episode, err := GetEpisodeTest(ep.Id)
		assert.Nil(t, err)
		episode.Watched = true
		episode.WatchedDate = generic.GetCurrentDate()
		testutils.CheckEpisode(t, ep, episode)
		err = UpdateEpisodeTest(episode)
		assert.Nil(t, err)
	}

	// testutils.ListAllEpisodes(DEBUG)
	// DumpEpisodeTest()
}

func TestEpisodeSummary(t *testing.T) {
	r := testutils.SetUpTestRoutes(true)
	url := "/episodes/summary/:id"
	r.GET(url, controllers.EpisodeSummaryBySeason)
	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, strings.Replace(url, ":id", strconv.Itoa(tvShowTest.Id), 1), nil)
	assert.Nil(t, err)
	r.ServeHTTP(w, req)
	// println(w.Body.String())

	type SeasonSummary struct {
		Season               int `json:"season"`
		TotalEpisodes        int `json:"total_episodes"`
		TotalEpisodesWatched int `json:"total_episodes_watched"`
	}

	var summary []SeasonSummary
	err = json.Unmarshal(w.Body.Bytes(), &summary)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 2, len(summary))
	var index int
	// season 1
	index = 0
	assert.Equal(t, 1, summary[index].Season)
	assert.Equal(t, 2, summary[index].TotalEpisodes)
	assert.Equal(t, 2, summary[index].TotalEpisodesWatched)
	// season 2
	index = 1
	assert.Equal(t, 2, summary[index].Season)
	assert.Equal(t, 1, summary[index].TotalEpisodes)
	assert.Equal(t, 0, summary[index].TotalEpisodesWatched)
}

func TestEpisodeSummaryErrors(t *testing.T) {

}

// func TestTvShowTruncate(t *testing.T) {
// 	r := SetUpTestRoutes(true)
// 	r.DELETE("/tvshows/truncate", controllers.TvShowTruncate)
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest(http.MethodDelete, "/tvshows/truncate", nil)
// 	r.ServeHTTP(w, req)

// 	log.Println(w.Body.String())
// 	assert.Equal(t, http.StatusOK, w.Code)
// }

// func TestEpisodeTruncate(t *testing.T) {

// }

// FUNCOES INTERNAS QUE CRIEI
// generic.GetCurrentDate()
