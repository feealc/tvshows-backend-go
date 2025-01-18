package database

import (
	"log"

	"github.com/feealc/tvshows-backend-go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func ConnectDataBase() {
	dsn := "host=localhost user=root password=root dbname=root port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn))

	if err != nil {
		log.Panic("Erro ao conectar com banco de dados")
	} else {
		log.Println("Database connected!")
	}

	DB.AutoMigrate(&models.TvShow{})
}
