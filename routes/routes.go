package routes

import (
	"github.com/feealc/tvshows-backend-go/controllers"
	"github.com/gin-gonic/gin"
)

func HandleRequests() {
	r := gin.Default()

	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Health
			v1.GET("/health", controllers.Health)

			// TvShows
			v1.GET("/tvshows", controllers.TvShowListAll)
			v1.GET("/tvshows/episodes", controllers.TvShowListAllUnwatchedEpisodes)
			v1.GET("/tvshows/:id", controllers.TvShowListById)
			v1.POST("/tvshows/create", controllers.TvShowCreate)
			v1.POST("/tvshows/create/batch", controllers.TvShowCreateBatch)
			v1.PUT("/tvshows/:id", controllers.TvShowEdit)
			v1.DELETE("/tvshows/:id", controllers.TvShowDelete)
			v1.DELETE("/tvshows/truncate", controllers.TvShowTruncate)

			// Episodes
			v1.GET("/episodes", controllers.EpisodeListAll)
			v1.GET("/episodes/:tmdbid", controllers.EpisodeListByTmdbId)
			v1.GET("/episodes/:tmdbid/:season", controllers.EpisodeListByTmdbIdAndSeason)
			v1.GET("/episodes/summary/:id", controllers.EpisodeSummaryBySeason)
			v1.POST("/episodes/create", controllers.EpisodeCreate)
			v1.POST("/episodes/create/batch", controllers.EpisodeCreateBatch)
			v1.PUT("/episodes/edit/:id", controllers.EpisodeEdit)
			v1.PUT("/episodes/watched/:id", controllers.EpisodeEditMarkWatched)
			v1.PUT("/episodes/watched/season/:tmdbid/:season", controllers.EpisodeEditMarkWatched)
			v1.DELETE("/episodes/delete/:id", controllers.EpisodeDelete)
			v1.DELETE("/episodes/delete/tvshow/:tmdbid", controllers.EpisodeDelete)
			v1.DELETE("/episodes/delete/season/:tmdbid/:season", controllers.EpisodeDelete)
			v1.DELETE("/episodes/truncate", controllers.EpisodeTruncate)

			//
			v1.DELETE("/truncate/all", controllers.TruncateAll)
		}
	}

	r.NoRoute(controllers.RouteNotFound)

	r.Run(":8080")
}
