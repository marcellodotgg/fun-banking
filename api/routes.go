package api

import (
	"github.com/bytebury/fun-banking/api/handler"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func Start() {
	// Load Templates
	router.LoadHTMLGlob("public/templates/**/*")
	// Load Static Files
	router.Static("/static", "public/static")
	// Middleware
	// Setup Routes
	setupHomepageRoutes()
	setupUserRoutes()
	// Run the application
	router.Run()
}

func setupHomepageRoutes() {
	handler := handler.NewHomePageHandler()

	router.
		GET("", handler.Homepage).
		GET("signin", handler.SignIn).
		POST("signin", handler.CreateSession).
		GET("signup", handler.SignUp)
}

func setupUserRoutes() {
	handler := handler.NewUserHandler()

	router.
		Group("users").
		GET("count", handler.Count)
}
