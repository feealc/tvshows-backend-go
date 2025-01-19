package routes

import (
	"github.com/feealc/tvshows-backend-go/controllers"
	"github.com/gin-gonic/gin"
)

func HandleRequests() {
	r := gin.Default()

	// TvShows
	r.GET("/tvshows", controllers.GetTvShows)
	r.GET("/tvshows/:id", controllers.GetTvShowId)
	r.POST("/tvshows/create", controllers.CreateTvShow)
	r.PUT("/tvshows/:id", controllers.EditTvShow)
	r.DELETE("/tvshows/:id", controllers.DeleteTvShow)
	// truncate

	// Episodes
	r.GET("/episodes", controllers.EpisodeListAll)
	// get tv show by id
	// get tv show by id e season
	r.POST("/episodes/create", controllers.EpisodeCreate)
	r.POST("/episodes/create/batch", controllers.EpisodeCreateBatch)
	// put watch episode (set watched date to now)
	r.DELETE("/episodes/truncate", controllers.EpisodeTruncate)

	r.Run()
}
