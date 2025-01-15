package main

import (
	"github.com/feealc/tvshows-backend-go/database"
	"github.com/feealc/tvshows-backend-go/routes"
)

func main() {
	database.ConnectDataBase()

	routes.HandleRequests()
}
