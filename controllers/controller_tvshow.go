package controllers

import (
	"fmt"
	"net/http"

	"github.com/feealc/tvshows-backend-go/database"
	"github.com/feealc/tvshows-backend-go/generic"
	"github.com/feealc/tvshows-backend-go/models"
	"github.com/gin-gonic/gin"
)

func TvShowListAll(c *gin.Context) {
	var tvShows []models.TvShow

	if result := database.DB.Order("name").Find(&tvShows); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	for index, tvShow := range tvShows {
		var episodes []models.Episode

		if result := database.DB.Where("tmdb_id = ? and watched = false", tvShow.TmdbId).Order(kEPISODE_ORDER_BY_TMDBID_SEASON_EPISODE).Find(&episodes); result.Error != nil {
			ResponseErrorInternalServerError(c, result.Error)
			return
		}

		if len(episodes) > 0 {
			ep := episodes[0]
			tvShow.UnwatchedSeason = ep.Season
			tvShow.UnwatchedEpisode = ep.Episode
			tvShow.UnwatchedCount = len(episodes) - 1
		}
		tvShows[index] = tvShow
	}

	c.JSON(http.StatusOK, tvShows)
}

func TvShowListAllUnwatchedEpisodes(c *gin.Context) {
	var tvShows []models.TvShow

	if result := database.DB.Order("name").Find(&tvShows); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	type TvShowEpisodes struct {
		TvShow   models.TvShow    `json:"tv_show"`
		Episodes []models.Episode `json:"episodes"`
	}
	var response []TvShowEpisodes

	for _, tvShow := range tvShows {
		var episodes []models.Episode

		if result := database.DB.Where("tmdb_id = ? and watched = false", tvShow.TmdbId).Order(kEPISODE_ORDER_BY_TMDBID_SEASON_EPISODE).Find(&episodes); result.Error != nil {
			ResponseErrorInternalServerError(c, result.Error)
			return
		}

		response = append(response, TvShowEpisodes{TvShow: tvShow, Episodes: episodes})
	}

	c.JSON(http.StatusOK, response)
}

func TvShowListById(c *gin.Context) {
	var tvShow models.TvShow
	paramId := c.Params.ByName("id")

	id, err := generic.CheckParamInt(paramId, kERROR_MESSAGE_ID)
	if err != nil {
		ResponseErrorBadRequest(c, err)
		return
	}

	if result := database.DB.Find(&tvShow, id); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	if tvShow.Id == 0 {
		ResponseErrorNotFound(c, models.TvShow{})
		return
	}

	c.JSON(http.StatusOK, tvShow)
}

func TvShowCreate(c *gin.Context) {
	var tvShow models.TvShow

	if err := c.ShouldBindJSON(&tvShow); err != nil {
		ResponseErrorBadRequest(c, err)
		return
	}

	if err := models.ValidTvShow(&tvShow); err != nil {
		ResponseErrorUnprocessableEntity(c, err)
		return
	}

	var tvShowExist models.TvShow
	if result := database.DB.Where(&models.TvShow{TmdbId: tvShow.TmdbId}).Find(&tvShowExist); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	if tvShowExist.Id > 0 {
		ResponseErrorBadRequest(c, fmt.Errorf("TvShow %s (TMDB ID %d) already exist", tvShow.Name, tvShow.TmdbId))
		return
	}

	if result := database.DB.Create(&tvShow); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	c.JSON(http.StatusCreated, tvShow)
}

func TvShowCreateBatch(c *gin.Context) {
	var tvShows []models.TvShow

	if err := c.ShouldBindJSON(&tvShows); err != nil {
		ResponseErrorBadRequest(c, err)
		return
	}

	for index, tvShow := range tvShows {
		if err := models.ValidTvShow(&tvShow); err != nil {
			ResponseErrorUnprocessableEntity(c, err)
			return
		}
		tvShows[index] = tvShow

		var tvShowExist models.TvShow
		if result := database.DB.Where(&models.TvShow{TmdbId: tvShow.TmdbId}).Find(&tvShowExist); result.Error != nil {
			ResponseErrorInternalServerError(c, result.Error)
			return
		}

		if tvShowExist.TmdbId > 0 {
			ResponseErrorBadRequest(c, fmt.Errorf("TvShow %s (TMDB ID %d) already exist", tvShow.Name, tvShow.TmdbId))
			return
		}
	}

	if result := database.DB.Create(&tvShows); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	c.JSON(http.StatusCreated, tvShows)
}

func TvShowEdit(c *gin.Context) {
	var tvShow models.TvShow
	paramId := c.Params.ByName("id")

	id, err := generic.CheckParamInt(paramId, kERROR_MESSAGE_ID)
	if err != nil {
		ResponseErrorBadRequest(c, err)
		return
	}

	if result := database.DB.Find(&tvShow, id); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	if tvShow.Id == 0 {
		ResponseErrorNotFound(c, models.TvShow{})
		return
	}

	if err := c.ShouldBindJSON(&tvShow); err != nil {
		ResponseErrorUnprocessableEntity(c, err)
		return
	}

	if err := models.ValidTvShow(&tvShow); err != nil {
		ResponseErrorUnprocessableEntity(c, err)
		return
	}

	if result := database.DB.Save(&tvShow); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	c.JSON(http.StatusOK, tvShow)
}

func TvShowDelete(c *gin.Context) {
	var tvShow models.TvShow
	paramId := c.Params.ByName("id")

	id, err := generic.CheckParamInt(paramId, kERROR_MESSAGE_ID)
	if err != nil {
		ResponseErrorBadRequest(c, err)
		return
	}

	if result := database.DB.Find(&tvShow, id); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	if tvShow.Id == 0 {
		ResponseErrorNotFound(c, models.TvShow{})
		return
	}

	if result := database.DB.Delete(&tvShow, id); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	var episodesToDelete []models.Episode
	if result := database.DB.Where(&models.Episode{TmdbId: tvShow.TmdbId}).Delete(&episodesToDelete); result.Error != nil {
		ResponseErrorInternalServerError(c, result.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "TvShow and episodes deleted successfully",
	})
}

func TvShowTruncate(c *gin.Context) {
	response, err := Truncate(c, models.TvShow{})
	if err != nil {
		ResponseErrorInternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}
