package api

import (
	"github.com/bytebury/fun-banking/internal/api/handler"
	"github.com/bytebury/fun-banking/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func Start() {
	// Load Templates
	router.LoadHTMLGlob("templates/**/*")
	// Load Static Files
	router.Static("/static", "public/")
	// Middleware
	router.Use(middleware.Audit())
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
		GET("", handler.Homepage)
}

func setupSessionRoutes() {
	handler := handler.NewSessionHandle()

	router.
		GET("signin", handler.SignIn).
		POST("signin", handler.CreateSession).
		DELETE("signout", handler.DestroySession)
}

func setupUserRoutes() {
	handler := handler.NewUserHandler()

	router.GET("signup", handler.SignUp)

	router.
		Group("users").
		GET("count", handler.Count).
		PUT("", handler.Create)
}
