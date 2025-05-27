package controllers

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/feealc/tvshows-backend-go/database"
	"github.com/feealc/tvshows-backend-go/generic"
	"github.com/feealc/tvshows-backend-go/models"
	"github.com/gin-gonic/gin"
)

func EpisodeListAll(c *gin.Context) {
	var episodes []models.Episode

	if result := database.DB.Order(kEPISODE_ORDER_BY_TMDBID_SEASON_EPISODE).Find(&episodes); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	c.JSON(http.StatusOK, episodes)
}

func EpisodeListByTmdbId(c *gin.Context) {
	paramTmdbId := c.Params.ByName("tmdbid")

	tmdbId, err := generic.CheckParamInt(paramTmdbId, kERROR_MESSAGE_TMDBID)
	if err != nil {
		ResponseErrorBadRequest(c, err)
		return
	}

	var tvShowExist models.TvShow
	if result := database.DB.Where(&models.TvShow{TmdbId: tmdbId}).Find(&tvShowExist); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	if tvShowExist.Id == 0 {
		ResponseErrorNotFound(c, models.TvShow{})
		return
	}

	var episodes []models.Episode
	if result := database.DB.Where(&models.Episode{TmdbId: tvShowExist.TmdbId}).Order(kEPISODE_ORDER_BY_TMDBID_SEASON_EPISODE).Find(&episodes); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	c.JSON(http.StatusOK, episodes)
}

func EpisodeListByTmdbIdAndSeason(c *gin.Context) {
	paramTmdbId := c.Params.ByName("tmdbid")
	paramSeason := c.Params.ByName("season")

	tmdbId, err := generic.CheckParamInt(paramTmdbId, kERROR_MESSAGE_TMDBID)
	if err != nil {
		ResponseErrorBadRequest(c, err)
		return
	}

	season, err := generic.CheckParamInt(paramSeason, kERROR_MESSAGE_SEASON)
	if err != nil {
		ResponseErrorBadRequest(c, err)
		return
	}

	var tvShowExist models.TvShow
	if result := database.DB.Where(&models.TvShow{TmdbId: tmdbId}).Find(&tvShowExist); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	if tvShowExist.Id == 0 {
		ResponseErrorNotFound(c, models.TvShow{})
		return
	}

	var episodes []models.Episode
	if result := database.DB.Where(&models.Episode{TmdbId: tmdbId, Season: season}).Order(kEPISODE_ORDER_BY_TMDBID_SEASON_EPISODE).Find(&episodes); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	c.JSON(http.StatusOK, episodes)
}

func EpisodeSummaryBySeason(c *gin.Context) {
	paramId := c.Params.ByName("id")

	id, err := generic.CheckParamInt(paramId, kERROR_MESSAGE_ID)
	if err != nil {
		ResponseErrorBadRequest(c, err)
		return
	}

	var tvShowExist models.TvShow
	if result := database.DB.Find(&tvShowExist, id); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	if tvShowExist.Id == 0 {
		ResponseErrorNotFound(c, models.TvShow{})
		return
	}

	var episodes []models.Episode
	if result := database.DB.Where(&models.Episode{TmdbId: tvShowExist.TmdbId}).Order(kEPISODE_ORDER_BY_TMDBID_SEASON_EPISODE).Find(&episodes); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	// if no episodes found, return empty slice
	if len(episodes) == 0 {
		c.JSON(http.StatusOK, episodes)
		return
	}

	episodesBySeasons := make(map[int][]models.Episode)

	for _, episode := range episodes {
		season := episode.Season

		aux := episodesBySeasons[season]
		aux = append(aux, episode)
		episodesBySeasons[season] = aux
		// log.Printf("append %dx%d", episode.Season, episode.Episode)
	}

	// sort map by key (key = season number)
	keys := make([]int, 0, len(episodesBySeasons))
	for k := range episodesBySeasons {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// var maxSeason int
	// if len(keys) > 0 {
	// 	maxSeason = keys[len(keys)-1]
	// }
	// fmt.Printf("max season [%d]\n", maxSeason)

	// for _, season := range keys {
	// 	log.Println("==============================")
	// 	for _, ep := range episodesBySeasons[season] {
	// 		log.Printf("ep %dx%d => %s \n", season, ep.Episode, ep.Name)
	// 	}
	// }

	type SeasonSummary struct {
		Season               int `json:"season"`
		TotalEpisodes        int `json:"total_episodes"`
		TotalEpisodesWatched int `json:"total_episodes_watched"`
	}

	var responseSummary []SeasonSummary
	var totalEpisodes, totalEpisodesWatched int
	for _, season := range keys {
		totalEpisodes = 0
		totalEpisodesWatched = 0
		for _, ep := range episodesBySeasons[season] {
			totalEpisodes += 1
			if ep.Watched {
				totalEpisodesWatched += 1
			}
		}
		responseSummary = append(responseSummary, SeasonSummary{Season: season, TotalEpisodes: totalEpisodes, TotalEpisodesWatched: totalEpisodesWatched})
	}

	c.JSON(http.StatusOK, responseSummary)
}

func EpisodeCreate(c *gin.Context) {
	var episode models.Episode

	if err := c.ShouldBindJSON(&episode); err != nil {
		ResponseErrorBadRequest(c, err)
		return
	}

	if err := models.ValidEpisode(&episode); err != nil {
		ResponseErrorUnprocessableEntity(c, err)
		return
	}

	var tvShowExist models.TvShow
	if result := database.DB.Where(&models.TvShow{TmdbId: episode.TmdbId}).Find(&tvShowExist); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	if tvShowExist.Id == 0 {
		ResponseErrorNotFound(c, models.TvShow{})
		return
	}

	var episodeExist models.Episode
	if result := database.DB.Where(&models.Episode{TmdbId: episode.TmdbId, Season: episode.Season, Episode: episode.Episode}).Find(&episodeExist); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	if episodeExist.Id > 0 {
		ResponseErrorBadRequest(c, fmt.Errorf("episode %dx%02d already exist for %s", episode.Season, episode.Episode, tvShowExist.Name))
		return
	}

	if result := database.DB.Create(&episode); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	c.JSON(http.StatusCreated, episode)
}

func EpisodeCreateBatch(c *gin.Context) {
	var episodes []models.Episode

	if err := c.ShouldBindJSON(&episodes); err != nil {
		ResponseErrorBadRequest(c, err)
		return
	}

	for index, episode := range episodes {
		if err := models.ValidEpisode(&episode); err != nil {
			ResponseErrorUnprocessableEntity(c, err)
			return
		}
		episodes[index] = episode

		var tvShowExist models.TvShow
		if result := database.DB.Where(&models.TvShow{TmdbId: episode.TmdbId}).Find(&tvShowExist); result.Error != nil {
			ResponseErrorInternalServerError(c, result.Error)
			return
		}

		if tvShowExist.Id == 0 {
			ResponseErrorNotFound(c, models.TvShow{})
			return
		}

		var episodeExist models.Episode
		if result := database.DB.Where(&models.Episode{TmdbId: episode.TmdbId, Season: episode.Season, Episode: episode.Episode}).Find(&episodeExist); result.Error != nil {
			ResponseErrorInternalServerError(c, result.Error)
			return
		}

		if episodeExist.Id > 0 {
			ResponseErrorBadRequest(c, fmt.Errorf("episode %dx%02d already exist for %s", episode.Season, episode.Episode, tvShowExist.Name))
			return
		}
	}

	if result := database.DB.Create(&episodes); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	c.JSON(http.StatusCreated, episodes)
}

func EpisodeEdit(c *gin.Context) {
	paramId := c.Params.ByName("id")

	id, err := generic.CheckParamInt(paramId, kERROR_MESSAGE_ID)
	if err != nil {
		ResponseErrorBadRequest(c, err)
		return
	}

	var episodeUpdate models.Episode
	if result := database.DB.Find(&episodeUpdate, id); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	if episodeUpdate.Id == 0 {
		ResponseErrorNotFound(c, models.Episode{})
		return
	}

	if err := c.ShouldBindJSON(&episodeUpdate); err != nil {
		ResponseErrorBadRequest(c, err)
		return
	}

	if err := models.ValidEpisode(&episodeUpdate); err != nil {
		ResponseErrorUnprocessableEntity(c, err)
		return
	}

	if result := database.DB.Save(&episodeUpdate); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	c.JSON(http.StatusOK, episodeUpdate)
}

func EpisodeEditMarkWatched(c *gin.Context) {
	paramId := c.Params.ByName("id")
	paramTmdbId := c.Params.ByName("tmdbid")
	paramSeason := c.Params.ByName("season")

	id, err := generic.CheckParamInt(paramId, kERROR_MESSAGE_ID)
	if err != nil {
		ResponseErrorBadRequest(c, err)
		return
	}

	tmdbId, err := generic.CheckParamInt(paramTmdbId, kERROR_MESSAGE_TMDBID)
	if err != nil {
		ResponseErrorBadRequest(c, err)
		return
	}

	season, err := generic.CheckParamInt(paramSeason, kERROR_MESSAGE_SEASON)
	if err != nil {
		ResponseErrorBadRequest(c, err)
		return
	}

	if paramId != "" {
		var episodeUpdate models.Episode
		if result := database.DB.Find(&episodeUpdate, id); result.Error != nil {
			ResponseErrorInternalServerError(c, result.Error)
			return
		}

		if episodeUpdate.Id == 0 {
			ResponseErrorNotFound(c, models.Episode{})
			return
		}

		episodeUpdate.Watched = !episodeUpdate.Watched
		if episodeUpdate.Watched {
			episodeUpdate.WatchedDate = generic.GetCurrentDate()
		} else {
			episodeUpdate.WatchedDate = 0
		}

		if result := database.DB.Save(&episodeUpdate); result.Error != nil {
			ResponseErrorInternalServerError(c, result.Error)
			return
		}

		c.JSON(http.StatusOK, episodeUpdate)
	} else {
		var tvShowExist models.TvShow
		if result := database.DB.Where(&models.TvShow{TmdbId: tmdbId}).Find(&tvShowExist); result.Error != nil {
			ResponseErrorInternalServerError(c, result.Error)
			return
		}

		if tvShowExist.Id == 0 {
			ResponseErrorNotFound(c, models.TvShow{})
			return
		}

		var episodesToUpdate []models.Episode
		if result := database.DB.Where(&models.Episode{TmdbId: tmdbId, Season: season}).Find(&episodesToUpdate); result.Error != nil {
			ResponseErrorInternalServerError(c, result.Error)
			return
		}

		if len(episodesToUpdate) == 0 {
			ResponseError(c, fmt.Errorf("episodes not found for season %d", season), http.StatusNotFound)
			return
		}

		for index, episode := range episodesToUpdate {
			// fmt.Println(episode)
			episode.Watched = true
			episode.WatchedDate = generic.GetCurrentDate()
			episodesToUpdate[index] = episode
		}

		if result := database.DB.Save(&episodesToUpdate); result.Error != nil {
			ResponseErrorInternalServerError(c, result.Error)
			return
		}

		c.JSON(http.StatusOK, episodesToUpdate)
	}
}

func EpisodeDelete(c *gin.Context) {
	paramId := c.Params.ByName("id")
	paramTmdbId := c.Params.ByName("tmdbid")
	paramSeason := c.Params.ByName("season")

	var err error

	if paramId != "" {
		var id int
		id, err = generic.CheckParamInt(paramId, kERROR_MESSAGE_ID)
		if err != nil {
			ResponseErrorBadRequest(c, err)
			return
		}

		var episode models.Episode
		if result := database.DB.Find(&episode, id); result.Error != nil {
			ResponseErrorInternalServerError(c, result.Error)
			return
		}

		if episode.Id == 0 {
			ResponseErrorNotFound(c, models.Episode{})
			return
		}

		if result := database.DB.Delete(&episode, id); result.Error != nil {
			ResponseErrorInternalServerError(c, result.Error)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Episode deleted",
		})
		return
	} else {
		var tmdbId, season int

		tmdbId, err = generic.CheckParamInt(paramTmdbId, kERROR_MESSAGE_TMDBID)
		if err != nil {
			ResponseErrorBadRequest(c, err)
			return
		}

		season, err = generic.CheckParamInt(paramSeason, kERROR_MESSAGE_SEASON)
		if err != nil {
			ResponseErrorBadRequest(c, err)
			return
		}

		var episodes []models.Episode
		result := database.DB.Where(&models.Episode{TmdbId: tmdbId, Season: season}).Delete(&episodes)

		if result.Error != nil {
			ResponseErrorInternalServerError(c, result.Error)
			return
		}

		if result.RowsAffected == 0 {
			ResponseErrorNotFound(c, models.Episode{})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Episodes deleted",
			"rows":    result.RowsAffected,
		})
		return
	}
}

func EpisodeTruncate(c *gin.Context) {
	response, err := Truncate(c, models.Episode{})
	if err != nil {
		ResponseErrorInternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}
