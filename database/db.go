package database

import (
	"fmt"
	"log"
	"os"

	"github.com/feealc/tvshows-backend-go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error

	DB_HOST     string = "localhost"
	DB_HOST_ENV string = ""
	DB_USER     string = "root"
	DB_PASSWORD string = "root"
	DB_NAME     string = "root"
	DB_PORT     string = "5432"
)

func ConnectDataBase() {
	DB_HOST_ENV = os.Getenv("DOCKER_DB_HOST")
	if DB_HOST_ENV != "" {
		DB_HOST = DB_HOST_ENV
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		DB_HOST,
		DB_USER,
		DB_PASSWORD,
		DB_NAME,
		DB_PORT)

	DB, err = gorm.Open(postgres.Open(dsn))

	if err != nil {
		log.Println(err.Error())
		log.Printf("dsn [%s]", dsn)
		log.Panic("Erro ao conectar com banco de dados")
	} else {
		log.Println("Database connected!")
	}

	DB.AutoMigrate(&models.TvShow{})
	DB.AutoMigrate(&models.Episode{})
}
