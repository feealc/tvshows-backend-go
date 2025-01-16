package controllers

import (
	"net/http"

	"github.com/feealc/tvshows-backend-go/database"
	"github.com/feealc/tvshows-backend-go/models"
	"github.com/gin-gonic/gin"
)

func GetTvShows(c *gin.Context) {
	var tvShows []models.TvShow
	database.DB.Order("name").Find(&tvShows)
	c.JSON(http.StatusOK, tvShows)
}

func CreateTvShows(c *gin.Context) {
	var tvShow models.TvShow
	if err := c.ShouldBindJSON(&tvShow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	if err := models.ValidTvShow(&tvShow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	database.DB.Create(&tvShow)
	c.JSON(http.StatusCreated, tvShow)
}
