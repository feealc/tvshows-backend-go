package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/feealc/tvshows-backend-go/database"
	"github.com/feealc/tvshows-backend-go/generic"
	"github.com/feealc/tvshows-backend-go/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	kEPISODE_ORDER_BY_TMDBID_SEASON_EPISODE = "tmdb_id, season, episode"
	kERROR_MESSAGE_ID                       = "id invalid"
	kERROR_MESSAGE_TMDBID                   = "tmdbId invalid"
	kERROR_MESSAGE_SEASON                   = "season invalid"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "Ok",
		"date_time": time.Now().String(),
	})
}

func RouteNotFound(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Route not found",
	})
}

func TvShowListAll(c *gin.Context) {
	var tvShows []models.TvShow
	database.DB.Order("name").Find(&tvShows)
	c.JSON(http.StatusOK, tvShows)
}

func TvShowListById(c *gin.Context) {
	var tvShow models.TvShow
	id := c.Params.ByName("id")
	database.DB.First(&tvShow, id)

	if tvShow.TmdbId == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "TvShow not found",
		})
		return
	}

	c.JSON(http.StatusOK, tvShow)
}

func TvShowCreate(c *gin.Context) {
	var tvShow models.TvShow

	if err := c.ShouldBindJSON(&tvShow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := models.ValidTvShow(&tvShow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	database.DB.Create(&tvShow)
	c.JSON(http.StatusCreated, tvShow)
}

func TvShowCreateBatch(c *gin.Context) {
	var tvShows []models.TvShow

	if err := c.ShouldBindJSON(&tvShows); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, tvShow := range tvShows {
		if err := models.ValidTvShow(&tvShow); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var tvShowExist models.TvShow
		database.DB.First(&tvShowExist, tvShow.TmdbId)

		if tvShowExist.TmdbId > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("TvShow %s (TMDB ID %d) already exist", tvShow.Name, tvShow.TmdbId),
			})
			return
		}
	}

	if result := database.DB.Create(&tvShows); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, tvShows)
}

func TvShowEdit(c *gin.Context) {
	var tvShow models.TvShow
	id := c.Params.ByName("id")
	database.DB.First(&tvShow, id)

	if tvShow.TmdbId == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "TvShow not found",
		})
		return
	}

	if err := c.ShouldBindJSON(&tvShow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := models.ValidTvShow(&tvShow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	database.DB.Save(&tvShow)
	c.JSON(http.StatusOK, tvShow)
}

func TvShowDelete(c *gin.Context) {
	var tvShow models.TvShow
	id := c.Params.ByName("id")
	database.DB.First(&tvShow, id)

	if tvShow.TmdbId == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "TvShow not found",
		})
		return
	}

	database.DB.Delete(&tvShow, id)
	c.JSON(http.StatusOK, gin.H{
		"message": "TvShow deleted successfully",
	})
}

func TvShowTruncate(c *gin.Context) {
	truncate(c, models.TvShow{})
}

//

func EpisodeListAll(c *gin.Context) {
	var episodes []models.Episode
	database.DB.Order(kEPISODE_ORDER_BY_TMDBID_SEASON_EPISODE).Find(&episodes)
	c.JSON(http.StatusOK, episodes)
}

func EpisodeListByTmdbId(c *gin.Context) {
	var tvShowExist models.TvShow
	tmdbId := c.Params.ByName("tmdbid")
	database.DB.First(&tvShowExist, tmdbId)

	if tvShowExist.TmdbId == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "TvShow not found",
		})
		return
	}

	var episodes []models.Episode
	database.DB.Where(&models.Episode{TmdbId: tvShowExist.TmdbId}).Order(kEPISODE_ORDER_BY_TMDBID_SEASON_EPISODE).Find(&episodes)

	c.JSON(http.StatusOK, episodes)
}

func EpisodeListByTmdbIdAndSeason(c *gin.Context) {
	var tvShowExist models.TvShow
	tmdbId := c.Params.ByName("tmdbid")
	season := c.Params.ByName("season")
	database.DB.First(&tvShowExist, tmdbId)

	if tvShowExist.TmdbId == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "TvShow not found",
		})
		return
	}

	seasonInt, err := strconv.Atoi(season)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "season invalid",
		})
		return
	}

	var episodes []models.Episode
	database.DB.Where(&models.Episode{TmdbId: tvShowExist.TmdbId, Season: seasonInt}).Order(kEPISODE_ORDER_BY_TMDBID_SEASON_EPISODE).Find(&episodes)

	c.JSON(http.StatusOK, episodes)
}

func EpisodeSummaryBySeason(c *gin.Context) {
	paramTmdbId := c.Params.ByName("tmdbid")

	tmdbId, err := generic.CheckParamInt(paramTmdbId, kERROR_MESSAGE_TMDBID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var tvShowExist models.TvShow
	database.DB.First(&tvShowExist, tmdbId)

	if tvShowExist.TmdbId == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "TvShow not found",
		})
		return
	}

	var maxSeason int
	database.DB.Model(&models.Episode{}).Select("max(season)").Group("tmdb_id").Where(&models.Episode{TmdbId: tvShowExist.TmdbId}).First(&maxSeason)
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

		database.DB.Model(&models.Episode{}).Select("count(*)").Group("tmdb_id").Where(&models.Episode{TmdbId: tvShowExist.TmdbId, Season: i}).First(&totalEpisodes)
		// fmt.Printf("total episodes %d \n", totalEpisodes)

		database.DB.Model(&models.Episode{}).Select("count(*)").Group("tmdb_id").Where(&models.Episode{TmdbId: tvShowExist.TmdbId, Season: i, Watched: true}).First(&totalEpisodesWatched)
		// fmt.Printf("total episodes watched %d \n", totalEpisodesWatched)

		responseSummary = append(responseSummary, SeasonSummary{Season: i, TotalEpisodes: totalEpisodes, TotalEpisodesWatched: totalEpisodesWatched})
	}

	c.JSON(http.StatusOK, responseSummary)
}

func EpisodeCreate(c *gin.Context) {
	var episode models.Episode

	if err := c.ShouldBindJSON(&episode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := models.ValidEpisode(&episode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var tvShow models.TvShow
	database.DB.First(&tvShow, episode.TmdbId)

	if tvShow.TmdbId == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "TvShow not found",
		})
		return
	}

	if ret := database.DB.Create(&episode); ret.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ret.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, episode)
}

func EpisodeCreateBatch(c *gin.Context) {
	var episodes []models.Episode

	if err := c.ShouldBindJSON(&episodes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, episode := range episodes {
		if err := models.ValidEpisode(&episode); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var tvShow models.TvShow
		database.DB.First(&tvShow, episode.TmdbId)

		if tvShow.TmdbId == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("TvShow %d not found", episode.TmdbId),
			})
			return
		}

		var episodeExist models.Episode
		rows := database.DB.Where(&models.Episode{TmdbId: episode.TmdbId, Season: episode.Season, Episode: episode.Episode}).First(&episodeExist).RowsAffected

		if rows > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Episode %dx%02d already exist for %s", episode.Season, episode.Episode, tvShow.Name),
			})
			return
		}
	}

	if result := database.DB.Create(&episodes); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, episodes)
}

func EpisodeEdit(c *gin.Context) {
	paramId := c.Params.ByName("id")

	id, err := generic.CheckParamInt(paramId, kERROR_MESSAGE_ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var episodeUpdate models.Episode
	result := database.DB.Find(&episodeUpdate, id)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if episodeUpdate.TmdbId == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Episode not found",
		})
		return
	}

	if err := c.ShouldBindJSON(&episodeUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := models.ValidEpisode(&episodeUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	database.DB.Save(&episodeUpdate)
	c.JSON(http.StatusOK, episodeUpdate)
}

func EpisodeEditMarkWatched(c *gin.Context) {
	paramId := c.Params.ByName("id")
	paramTmdbId := c.Params.ByName("tmdbid")
	paramSeason := c.Params.ByName("season")

	var id, tmdbId, season int

	for index, value := range []string{paramId, paramTmdbId, paramSeason} {
		var message string
		var err error

		switch index {
		case 0:
			message = kERROR_MESSAGE_ID
			id, err = generic.CheckParamInt(value, message)
		case 1:
			message = kERROR_MESSAGE_TMDBID
			tmdbId, err = generic.CheckParamInt(value, message)
		case 2:
			message = kERROR_MESSAGE_SEASON
			season, err = generic.CheckParamInt(value, message)
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	if paramId != "" {
		var episodeUpdate models.Episode
		result := database.DB.Find(&episodeUpdate, id)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": result.Error.Error(),
			})
			return
		}

		if episodeUpdate.Id == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Episode not found",
			})
			return
		}

		episodeUpdate.Watched = !episodeUpdate.Watched
		if episodeUpdate.Watched {
			episodeUpdate.WatchedDate = generic.GetCurrentDate()
		} else {
			episodeUpdate.WatchedDate = 0
		}

		result = database.DB.Save(&episodeUpdate)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, episodeUpdate)
	} else {
		var episodesToUpdate []models.Episode
		result := database.DB.Where(&models.Episode{TmdbId: tmdbId, Season: season}).Find(&episodesToUpdate)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": result.Error.Error(),
			})
			return
		}

		if len(episodesToUpdate) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Episodes not found",
			})
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

		result = database.DB.Save(&episodesToUpdate)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": result.Error.Error(),
			})
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
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var episode models.Episode
		database.DB.Find(&episode, id)

		if episode.Id == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "episode not found",
			})
			return
		}

		database.DB.Delete(&episode, id)

		c.JSON(http.StatusOK, gin.H{
			"message": "Episode deleted",
		})
		return
	} else {
		var result *gorm.DB
		var tmdbId, season int

		tmdbId, err = generic.CheckParamInt(paramTmdbId, kERROR_MESSAGE_TMDBID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		season, err = generic.CheckParamInt(paramSeason, kERROR_MESSAGE_SEASON)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var episodes []models.Episode
		result = database.DB.Where(&models.Episode{TmdbId: tmdbId, Season: season}).Delete(&episodes)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": result.Error.Error(),
			})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Episodes not found",
			})
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
	if result := database.DB.Where("tmdb_id is not null").Delete(&models.Episode{}); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Episodes truncated",
	})
}

//

func truncate(c *gin.Context, table interface{}) {
	mode := c.Query("mode")
	// mode := c.DefaultQuery("mode", "drop and create")
	response := gin.H{
		"message": "TvShows truncated",
	}

	if mode == "delete" {
		if result := database.DB.Where("id is not null").Delete(&table); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": result.Error.Error(),
			})
			return
		}
	} else {
		if err := database.DB.Migrator().DropTable(&table); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := database.DB.Migrator().CreateTable(&table); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		response["mode"] = "drop and create"
	}

	c.JSON(http.StatusOK, response)
}
