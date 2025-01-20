package routes

import (
	"github.com/feealc/tvshows-backend-go/controllers"
	"github.com/gin-gonic/gin"
)

func HandleRequests() {
	r := gin.Default()

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
	r.POST("/episodes/create", controllers.EpisodeCreate)
	r.POST("/episodes/create/batch", controllers.EpisodeCreateBatch)
	// put watch episode (set watched date to now)
	r.DELETE("/episodes/truncate", controllers.EpisodeTruncate)

	r.Run()
}
