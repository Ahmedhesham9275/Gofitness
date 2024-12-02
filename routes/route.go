package routes

import (
	"fitnesshub/controllers"
	"fitnesshub/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/login", controllers.Login)
	router.POST("/register", controllers.Register)

	protected := router.Group("/api")
	protected.GET("/package", controllers.GetPackages)
	protected.GET("/package/:id", controllers.GetPackage)
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/package", controllers.CreatePackage)
		protected.PUT("/package/:id", controllers.UpdatePackage)
		protected.DELETE("/package/:id", controllers.DeletePackage)
	}

	return router
}
