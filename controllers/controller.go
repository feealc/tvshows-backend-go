package controllers

import (
	"fmt"
	"net/http"

	"github.com/feealc/tvshows-backend-go/database"
	"github.com/feealc/tvshows-backend-go/models"
	"github.com/gin-gonic/gin"
)

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
	database.DB.Order("tmdb_id, season, episode").Find(&episodes)
	c.JSON(http.StatusOK, episodes)
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
				"error": fmt.Sprintf("Episode %dx%02d already exist", episode.Season, episode.Episode),
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
