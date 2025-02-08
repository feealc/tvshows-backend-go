package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/feealc/tvshows-backend-go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	sqlDB *sql.DB
	DB    *gorm.DB
	err   error

	DB_HOST string
	DB_PORT string
	DB_USER string
	DB_PASS string
	DB_NAME string
)

func ConnectDataBase() {
	dsn := buildConnectionString()

	sqlDB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Println(err.Error())
		log.Printf("dsn [%s]", dsn)
		log.Panic("Erro ao conectar com banco de dados usando SQL")
		return
	}

	DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		log.Println(err.Error())
		log.Printf("dsn [%s]", dsn)
		log.Panic("Erro ao conectar com banco de dados usando sqlDB e GORM")
		return
	}

	sqlDB = nil
	log.Println("Conectado com sucesso usando GORM")

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
