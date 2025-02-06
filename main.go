package main

import (
	"log"

	"github.com/feealc/tvshows-backend-go/database"
	"github.com/feealc/tvshows-backend-go/routes"
)

func main() {
	// database.ConnectDataBase()

	_, err, dsn := database.ConnectUnixSocket()
	if err != nil {
		log.Println(err.Error())
		log.Printf("dsn [%s]", dsn)
		log.Panic("Erro ao conectar com banco de dados usando ConnectUnixSocket")
	} else {
		log.Println("Conectado com sucesso usando ConnectUnixSocket")
	}

	routes.HandleRequests()
}
