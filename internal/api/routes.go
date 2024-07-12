package api

import (
	"github.com/bytebury/fun-banking/api/handler"
	"github.com/bytebury/fun-banking/api/middleware"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func Start() {
	// Load Templates
	router.LoadHTMLGlob("templates/**/*")
	// Load Static Files
	router.Static("/static", "public/")
	// Middleware
	// Setup Routes
	setupHomepageRoutes()
	setupUserRoutes()
	setupSessionRoutes()
	// Run the application
	router.Run()
}

func setupHomepageRoutes() {
	handler := handler.NewHomePageHandler()

	router.
		GET("", middleware.Authenticated(), handler.Homepage)
}

func setupSessionRoutes() {
	handler := &handler.SessionHandler{}

	router.
		GET("signin", handler.SignIn).
		POST("signin", handler.CreateSession)
}

func setupUserRoutes() {
	handler := handler.NewUserHandler()

	router.GET("signup", handler.SignUp)

	router.
		Group("users").
		GET("count", handler.Count).
		PUT("", handler.Create)
}
