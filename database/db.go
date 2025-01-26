package database

import (
	"fmt"
	"log"
	"os"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/feealc/tvshows-backend-go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error

	CLOUD_NAME string = "googlecloud"

	DB_HOST string
	DB_PORT string
	DB_USER string
	DB_PASS string
	DB_NAME string
)

func ConnectDataBase() {
	dsn := buildConnectionString()

	if os.Getenv("CLOUD_NAME") == CLOUD_NAME {
		log.Println("Connecting to Cloud SQL database")
		DB, err = gorm.Open(postgres.New(postgres.Config{
			DriverName: "cloudsqlpostgres",
			DSN:        dsn,
		}))
	} else {
		log.Println("Connecting to local database")
		DB, err = gorm.Open(postgres.Open(dsn))
	}

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

func buildConnectionString() string {
	DB_HOST = os.Getenv("DB_HOST")
	DB_PORT = os.Getenv("DB_PORT")
	DB_USER = os.Getenv("DB_USER")
	DB_PASS = os.Getenv("DB_PASS")
	DB_NAME = os.Getenv("DB_NAME")

	if DB_HOST_ENV := os.Getenv("DOCKER_DB_HOST"); DB_HOST_ENV != "" {
		DB_HOST = DB_HOST_ENV
	}

	if DB_HOST_ENV := os.Getenv("CLOUD_SQL_CONNECTION_NAME"); DB_HOST_ENV != "" {
		DB_HOST = DB_HOST_ENV
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		DB_HOST,
		DB_USER,
		DB_PASS,
		DB_NAME,
		DB_PORT)

	// log.Printf("dsn [%s]\n", dsn)
	return dsn
}
