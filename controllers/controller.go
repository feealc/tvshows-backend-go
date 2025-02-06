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
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "Ok",
		"date_time": time.Now().String(),
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
	if result := database.DB.Where("tmdb_id is not null").Delete(&models.TvShow{}); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "TvShows truncated",
	})
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

	err := generic.CheckParamsInt(paramTmdbId, "", "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var tvShowExist models.TvShow
	database.DB.First(&tvShowExist, paramTmdbId)

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
	paramTmdbId := c.Params.ByName("tmdbid")
	paramSeason := c.Params.ByName("season")
	paramEpisode := c.Params.ByName("episode")

	err := generic.CheckParamsInt(paramTmdbId, paramSeason, paramEpisode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	tmdbId, _ := strconv.Atoi(paramTmdbId)
	season, _ := strconv.Atoi(paramSeason)
	episode, _ := strconv.Atoi(paramEpisode)

	var episodeUpdate models.Episode
	result := database.DB.Where(&models.Episode{TmdbId: tmdbId, Season: season, Episode: episode}).First(&episodeUpdate)

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
	paramTmdbId := c.Params.ByName("tmdbid")
	paramSeason := c.Params.ByName("season")
	paramEpisode := c.Params.ByName("episode")

	err := generic.CheckParamsInt(paramTmdbId, paramSeason, paramEpisode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	tmdbId, _ := strconv.Atoi(paramTmdbId)
	season, _ := strconv.Atoi(paramSeason)
	episode, _ := strconv.Atoi(paramEpisode)

	if paramEpisode != "" {
		var episodeUpdate models.Episode
		result := database.DB.Where(&models.Episode{TmdbId: tmdbId, Season: season, Episode: episode}).First(&episodeUpdate)

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

		episodeUpdate.Watched = !episodeUpdate.Watched
		if episodeUpdate.Watched {
			episodeUpdate.WatchedDate = generic.GetCurrentDate()
		} else {
			episodeUpdate.WatchedDate = 0
		}

		database.DB.Save(&episodeUpdate)
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
			fmt.Println(episode)
			episode.Watched = !episode.Watched
			if episode.Watched {
				episode.WatchedDate = generic.GetCurrentDate()
			} else {
				episode.WatchedDate = 0
			}
			episodesToUpdate[index] = episode
		}

		database.DB.Save(&episodesToUpdate)
		c.JSON(http.StatusOK, episodesToUpdate)
	}
}

func EpisodeDelete(c *gin.Context) {
	paramTmdbId := c.Params.ByName("tmdbid")
	paramSeason := c.Params.ByName("season")
	paramEpisode := c.Params.ByName("episode")

	var err error
	var tmdbId, season, episode int

	tmdbId, err = strconv.Atoi(paramTmdbId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "tmdbId invalid",
		})
		return
	}

	if paramSeason != "" {
		season, err = strconv.Atoi(paramSeason)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "season invalid",
			})
			return
		}
	}

	if paramEpisode != "" {
		episode, err = strconv.Atoi(paramEpisode)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "episode invalid",
			})
			return
		}
	}

	var result *gorm.DB
	if paramTmdbId != "" && paramSeason == "" && paramEpisode == "" {
		result = database.DB.Where(&models.Episode{TmdbId: tmdbId}).Delete(&models.Episode{})
	} else if paramTmdbId != "" && paramSeason != "" && paramEpisode == "" {
		result = database.DB.Where(&models.Episode{TmdbId: tmdbId, Season: season}).Delete(&models.Episode{})
	} else if paramTmdbId != "" && paramSeason != "" && paramEpisode != "" {
		result = database.DB.Where(&models.Episode{TmdbId: tmdbId, Season: season, Episode: episode}).Delete(&models.Episode{})
	}

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
