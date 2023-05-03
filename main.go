package main

import (
	"github.com/gin-gonic/gin"
	docs "github.com/reemployed/reemployed/docs"
	"github.com/reemployed/reemployed/handlers"
	"github.com/reemployed/reemployed/repositories"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	repo := repositories.NewFileUserRepository("users.json")
	uc := handlers.NewUserController(repo)
	ac := handlers.NewAuthController(repo)

	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	//Initialize the user controller and add the routes to the router
	v1userRoutes := r.Group("/api/v1/users")
	v1userRoutes.GET("/", uc.GetUserList)
	v1userRoutes.GET("/:id", uc.GetUserByID)
	v1userRoutes.POST("/", uc.CreateUser)
	v1userRoutes.PUT("/:id", uc.UpdateUser)
	v1userRoutes.DELETE("/:id", uc.DeleteUser)

	// Initialize the auth controller and add the login route to the router
	authRoutes := r.Group("/api/v1/auth")
	{
		authRoutes.POST("/login", ac.Login)
	}

	// Serve Swagger documentation

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run() // listen and serve on 0.0.0.0:8080 by default
}
