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
		// cleanup, err2 := pgxv5.RegisterDriver("cloudsql-postgres")
		// if err2 != nil {
		// 	panic(err2)
		// }
		// cleanup will stop the driver from retrieving ephemeral certificates
		// Don't call cleanup until you're done with your database connections
		// defer cleanup()

		log.Println("Connecting to Cloud SQL database")
		// DB, err = gorm.Open(postgres.New(postgres.Config{
		// 	DriverName: "cloudsql-postgres",
		// 	DSN:        dsn,
		// }))
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

	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		DB_HOST,
		DB_USER,
		DB_PASS,
		DB_NAME,
		DB_PORT)

	// dsn := fmt.Sprintf("user=%s password=%s database=%s host=%s",
	// 	DB_USER,
	// 	DB_PASS,
	// 	DB_NAME,
	// 	DB_HOST)

	// log.Printf("dsn [%s]\n", dsn)
	return dsn
}

func ConnectUnixSocket() (*sql.DB, error, string) {
	// mustGetenv := func(k string) string {
	// 	v := os.Getenv(k)
	// 	if v == "" {
	// 		log.Fatalf("Fatal Error in connect_unix.go: %s environment variable not set.\n", k)
	// 	}
	// 	return v
	// }
	// Note: Saving credentials in environment variables is convenient, but not
	// secure - consider a more secure solution such as
	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
	// keep secrets safe.
	// var (
	// 	dbUser         = mustGetenv("DB_USER")              // e.g. 'my-db-user'
	// 	dbPwd          = mustGetenv("DB_PASS")              // e.g. 'my-db-password'
	// 	unixSocketPath = mustGetenv("INSTANCE_UNIX_SOCKET") // e.g. '/cloudsql/project:region:instance'
	// 	dbName         = mustGetenv("DB_NAME")              // e.g. 'my-database'
	// )

	// dbURI := fmt.Sprintf("user=%s password=%s database=%s host=%s",
	// 	dbUser, dbPwd, dbName, unixSocketPath)

	dsn := buildConnectionString()

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("pgx", dsn)
	if err != nil {
		// return nil, fmt.Errorf("sql.Open: %w", err)
		return nil, err, dsn
	}

	// ...

	return dbPool, nil, dsn
}
