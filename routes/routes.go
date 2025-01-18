package routes

import (
	"github.com/feealc/tvshows-backend-go/controllers"
	"github.com/gin-gonic/gin"
)

func HandleRequests() {
	r := gin.Default()

	// GET
	r.GET("/tvshows", controllers.GetTvShows)

	// POST
	r.POST("/tvshows/create", controllers.CreateTvShows)

	// PUT
	r.PUT("/tvshows/:id", controllers.EditTvShow)

	// PATCH

	// DELETE
	r.DELETE("/tvshows/:id", controllers.DeleteTvShow)

	// r.GET("/:nome", controllers.Saudacao)
	// r.POST("/alunos", controllers.CriaNovoAluno)
	// r.GET("/alunos/:id", controllers.BuscaAlunoPorID)
	// r.PATCH("/alunos/:id", controllers.EditaAluno)
	// r.GET("/alunos/cpf/:cpf", controllers.BuscaAlunoPorCPF)
	r.Run()
}
