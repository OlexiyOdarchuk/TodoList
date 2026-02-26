package routes

import (
	_ "todolist/docs"
	"todolist/internal/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, userHandler *UserHandler, todoHandler *TodoHandler, jwtManager *utils.JWTManager) {
	router.Use(LoggerMiddleware())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(corsConfig))

	router.Static("/assets", "./frontend/dist/assets")
	router.StaticFile("/icon.svg", "./frontend/dist/icon.svg")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")

	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
		auth.POST("/google", userHandler.GoogleLogin)
		auth.POST("/verify", userHandler.VerifyEmail)
	}
	protected := router.Group("/api")
	protected.Use(AuthMiddleware(jwtManager))
	{
		protected.POST("/todos", todoHandler.Create)
		protected.PUT("/todos/:id", todoHandler.Update)
		protected.GET("/todos", todoHandler.GetAll)
		protected.DELETE("/todos/:id", todoHandler.Delete)

		protected.GET("/user/me", userHandler.GetUser)
		protected.DELETE("/user/me", userHandler.DeleteUser)
		protected.PUT("/user/me/delete", userHandler.VerifyEmailDelete)
		protected.PATCH("/user/me", userHandler.UpdateUsername)
		protected.PUT("/user/me/password", userHandler.UpdatePassword)
		protected.POST("/user/me/email", userHandler.RequestEmailUpdate)
		protected.PUT("/user/me/email", userHandler.VerifyEmailUpdate)
	}
	router.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})
}
