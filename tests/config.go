package tests

import (
	"fmt"
	"myblog/config"
	"myblog/models"
	"myblog/routes"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var router *gin.Engine
var dbOnce sync.Once

func setupRouter() *gin.Engine {
	if router == nil {
		router = routes.SetupRouter()
	}
	return router
}

func initTestDB() {
	dbOnce.Do(func() {
		// Load the .env file
		err := godotenv.Load("../.env")
		if err != nil {
			panic("Error loading .env file: " + err.Error())
		}

		// Load test environment variables
		dbHost := os.Getenv("TEST_DB_HOST")
		dbPort := os.Getenv("TEST_DB_PORT")
		dbUser := os.Getenv("TEST_DB_USER")
		dbPassword := os.Getenv("TEST_DB_PASSWORD")
		dbName := os.Getenv("TEST_DB_NAME")

		if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
			panic("One or more test database environment variables are missing")
		}

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			dbHost, dbUser, dbPassword, dbName, dbPort)

		// Initialize the database connection
		var dbErr error
		config.DB, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if dbErr != nil {
			panic("Failed to connect to test database: " + dbErr.Error())
		}

		// Migrate the models
		migrationErr := config.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.SearchStatistic{})
		if migrationErr != nil {
			panic("Failed to migrate models: " + migrationErr.Error())
		}
	})
}

func clearTables() {
	initTestDB()

	config.DB.Exec("DELETE FROM users")
	config.DB.Exec("DELETE FROM posts")
}
