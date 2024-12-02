package main

import (
	"fitnesshub/database"
	"fitnesshub/routes"
)

func main() {
	// Initialize database connection
	database.ConnectDatabase()

	// Setup routes
	router := routes.SetupRouter()

	// Run the server
	router.Run(":8080")
}
