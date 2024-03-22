package main

import (
	// "log"
	"mygram/configs"
	"mygram/models"
	"mygram/routes"
	"os"

	"github.com/gin-gonic/gin"
	// "github.com/joho/godotenv"
)

func main() {
	// Create a new gin instance
	r := gin.Default()

	// // Load .env file and Create a new connection to the database, uncomment to test locally
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	var db string = os.Getenv("DB")
	if db == "mysql" {
		config := models.ConfigMySQL{
			Host:      os.Getenv("DB_HOST"),
			Port:      os.Getenv("DB_PORT"),
			User:      os.Getenv("DB_USER"),
			Password:  os.Getenv("DB_PASSWORD"),
			DBName:    os.Getenv("DB_NAME"),
			DBCharset: os.Getenv("DB_CHARSET"),
		}
		configs.StartMySQL(config)
	} else if db == "postgres" {
		config := models.ConfigPostgres{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		}
		configs.StartPostgres(config)
	}

	// Load the routes
	routes.UserRoutes(r)
	routes.PhotoRoutes(r)
	routes.CommentRoutes(r)
	routes.SocialMediaRoutes(r)

	port := os.Getenv("APP_PORT")

	// Run the server
	r.Run(":" + port)
}
