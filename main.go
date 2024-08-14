package main

import (
	"myblog/config"
	"myblog/routes"
)

func main() {
	// Initialize database connection
	config.ConnectDatabase()

	// Setup routes
	router := routes.SetupRouter()

	// Run the server
	router.Run(":8080")
}
