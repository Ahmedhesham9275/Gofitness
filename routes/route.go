package routes

import (
	"myblog/controllers"
	"myblog/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/login", controllers.Login)
	router.POST("/register", controllers.Register)
	router.POST("/search_words", controllers.SearchWords)

	protected := router.Group("/api")
	protected.GET("/posts", controllers.GetPosts)
	protected.GET("/posts/:id", controllers.GetPost)
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/posts", controllers.CreatePost)
		protected.PUT("/posts/:id", controllers.UpdatePost)
		protected.DELETE("/posts/:id", controllers.DeletePost)
	}

	return router
}
