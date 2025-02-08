package controllers

import (
	"fmt"
	"net/http"

	"github.com/feealc/tvshows-backend-go/database"
	"github.com/feealc/tvshows-backend-go/generic"
	"github.com/feealc/tvshows-backend-go/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	var tvShowExist models.TvShow
	paramTmdbId := c.Params.ByName("tmdbid")

	tmdbId, err := generic.CheckParamInt(paramTmdbId, kERROR_MESSAGE_TMDBID)
	if err != nil {
		ResponseError(c, err, 0)
		return
	}

	if result := database.DB.First(&tvShowExist, tmdbId); result.Error != nil {
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
	var tvShowExist models.TvShow
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

	if result := database.DB.First(&tvShowExist, tmdbId); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	if tvShowExist.Id == 0 {
		ResponseErrorNotFound(c, models.TvShow{})
		return
	}

	var episodes []models.Episode
	if result := database.DB.Where(&models.Episode{TmdbId: tvShowExist.TmdbId, Season: season}).Order(kEPISODE_ORDER_BY_TMDBID_SEASON_EPISODE).Find(&episodes); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	c.JSON(http.StatusOK, episodes)
}

func EpisodeSummaryBySeason(c *gin.Context) {
	paramTmdbId := c.Params.ByName("tmdbid")

	tmdbId, err := generic.CheckParamInt(paramTmdbId, kERROR_MESSAGE_TMDBID)
	if err != nil {
		ResponseErrorBadRequest(c, err)
		return
	}

	var tvShowExist models.TvShow
	if result := database.DB.First(&tvShowExist, tmdbId); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	if tvShowExist.Id == 0 {
		ResponseErrorNotFound(c, models.TvShow{})
		return
	}

	var maxSeason int
	if result := database.DB.Model(&models.Episode{}).Select("max(season)").Group("tmdb_id").Where(&models.Episode{TmdbId: tvShowExist.TmdbId}).First(&maxSeason); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}
	// fmt.Printf("max season [%d]\n", maxSeason)

	type SeasonSummary struct {
		Season               int `json:"season"`
		TotalEpisodes        int `json:"total_episodes"`
		TotalEpisodesWatched int `json:"total_episodes_watched"`
	}

	var responseSummary []SeasonSummary
	var totalEpisodes, totalEpisodesWatched int
	for i := 1; i <= maxSeason; i++ {
		// fmt.Printf("season %d \n", i)

		if result := database.DB.Model(&models.Episode{}).Select("count(*)").Group("tmdb_id").Where(&models.Episode{TmdbId: tvShowExist.TmdbId, Season: i}).First(&totalEpisodes); result.Error != nil {
			ResponseErrorInternalServerError(c, result.Error)
			return
		}
		// fmt.Printf("total episodes %d \n", totalEpisodes)

		if result := database.DB.Model(&models.Episode{}).Select("count(*)").Group("tmdb_id").Where(&models.Episode{TmdbId: tvShowExist.TmdbId, Season: i, Watched: true}).First(&totalEpisodesWatched); result.Error != nil {
			ResponseErrorInternalServerError(c, result.Error)
			return
		}
		// fmt.Printf("total episodes watched %d \n", totalEpisodesWatched)

		responseSummary = append(responseSummary, SeasonSummary{Season: i, TotalEpisodes: totalEpisodes, TotalEpisodesWatched: totalEpisodesWatched})
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
		ResponseErrorBadRequest(c, err)
		return
	}

	var tvShow models.TvShow
	if result := database.DB.Find(&tvShow, episode.TmdbId); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	if tvShow.Id == 0 {
		ResponseErrorNotFound(c, models.TvShow{})
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

	for _, episode := range episodes {
		if err := models.ValidEpisode(&episode); err != nil {
			ResponseErrorBadRequest(c, err)
			return
		}

		var tvShow models.TvShow
		if result := database.DB.First(&tvShow, episode.TmdbId); result.Error != nil {
			ResponseErrorInternalServerError(c, result.Error)
			return
		}

		if tvShow.Id == 0 {
			ResponseErrorNotFound(c, fmt.Errorf("TvShow %d not found", episode.TmdbId))
			return
		}

		var episodeExist models.Episode
		result := database.DB.Where(&models.Episode{TmdbId: episode.TmdbId, Season: episode.Season, Episode: episode.Episode}).First(&episodeExist)

		if result.Error != nil {
			ResponseErrorInternalServerError(c, result.Error)
		}

		if result.RowsAffected > 0 {
			ResponseErrorBadRequest(c, fmt.Errorf("episode %dx%02d already exist for %s", episode.Season, episode.Episode, tvShow.Name))
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
		ResponseErrorBadRequest(c, err)
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
		var episodesToUpdate []models.Episode
		if result := database.DB.Where(&models.Episode{TmdbId: tmdbId, Season: season}).Find(&episodesToUpdate); result.Error != nil {
			ResponseErrorInternalServerError(c, result.Error)
			return
		}

		if len(episodesToUpdate) == 0 {
			ResponseErrorNotFound(c, models.Episode{})
			return
		}

		for index, episode := range episodesToUpdate {
			// fmt.Println(episode)
			episode.Watched = !episode.Watched
			if episode.Watched {
				episode.WatchedDate = generic.GetCurrentDate()
			} else {
				episode.WatchedDate = 0
			}
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
		var result *gorm.DB
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
		result = database.DB.Where(&models.Episode{TmdbId: tmdbId, Season: season}).Delete(&episodes)

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
	Truncate(c, models.Episode{})
}
