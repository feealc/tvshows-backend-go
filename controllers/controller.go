package controllers

import (
	"net/http"
	"time"

	"github.com/feealc/tvshows-backend-go/database"
	"github.com/feealc/tvshows-backend-go/generic"
	"github.com/gin-gonic/gin"
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

func Truncate(c *gin.Context, table interface{}) {
	mode := c.Query("mode")
	// mode := c.DefaultQuery("mode", "drop and create")

	name := generic.GetStructName(table)
	response := gin.H{
		"message": name + " truncated",
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
