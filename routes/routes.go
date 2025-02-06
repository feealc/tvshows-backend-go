package routes

import (
	"net/http"

	"github.com/feealc/tvshows-backend-go/controllers"
	"github.com/gin-gonic/gin"
)

func HandleRequests() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Ok",
		})
	})

	// TvShows
	r.GET("/tvshows", controllers.TvShowListAll)
	r.GET("/tvshows/:id", controllers.TvShowListById)
	r.POST("/tvshows/create", controllers.TvShowCreate)
	r.POST("/tvshows/create/batch", controllers.TvShowCreateBatch)
	r.PUT("/tvshows/:id", controllers.TvShowEdit)
	r.DELETE("/tvshows/:id", controllers.TvShowDelete)
	r.DELETE("/tvshows/truncate", controllers.TvShowTruncate)

	// Episodes
	r.GET("/episodes", controllers.EpisodeListAll)
	r.GET("/episodes/:tmdbid", controllers.EpisodeListByTmdbId)
	r.GET("/episodes/:tmdbid/:season", controllers.EpisodeListByTmdbIdAndSeason)
	r.GET("/episodes/:tmdbid/summary", controllers.EpisodeSummaryBySeason)
	r.POST("/episodes/create", controllers.EpisodeCreate)
	r.POST("/episodes/create/batch", controllers.EpisodeCreateBatch)
	r.PUT("/episodes/:tmdbid/:season/:episode", controllers.EpisodeEdit)
	r.PUT("/episodes/:tmdbid/:season/watched", controllers.EpisodeEditMarkWatched)
	r.PUT("/episodes/:tmdbid/:season/:episode/watched", controllers.EpisodeEditMarkWatched)
	r.DELETE("/episodes/:tmdbid", controllers.EpisodeDelete)
	r.DELETE("/episodes/:tmdbid/:season", controllers.EpisodeDelete)
	r.DELETE("/episodes/:tmdbid/:season/:episode", controllers.EpisodeDelete)
	r.DELETE("/episodes/truncate", controllers.EpisodeTruncate)

	r.Run(":8080")
}
