package routes

import (
	"github.com/feealc/tvshows-backend-go/controllers"
	"github.com/gin-gonic/gin"
)

func HandleRequests() {
	r := gin.Default()

	// TvShows
	r.GET("/tvshows", controllers.GetTvShows)
	r.POST("/tvshows/create", controllers.CreateTvShows)
	r.PUT("/tvshows/:id", controllers.EditTvShow)
	r.DELETE("/tvshows/:id", controllers.DeleteTvShow)

	r.Run()
}
