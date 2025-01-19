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
	r.PUT("/tvshows/:id", controllers.TvShowEdit)
	r.DELETE("/tvshows/:id", controllers.TvShowDelete)
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
