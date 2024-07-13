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
	setupBankRoutes()
	setupCustomerRoutes()
	// Run the application
	router.Run()
}

func setupHomepageRoutes() {
	handler := handler.NewHomePageHandler()

	router.
		GET("", handler.Homepage)
}

func setupSessionRoutes() {
	handler := handler.NewSessionHandler()

	router.
		GET("signin", middleware.NoAuth(), handler.SignIn).
		POST("signin", middleware.NoAuth(), handler.CreateSession).
		DELETE("signout", handler.DestroySession)
}

func setupUserRoutes() {
	handler := handler.NewUserHandler()

	router.GET("signup", middleware.NoAuth(), handler.SignUp)

	router.
		Group("users").
		PUT("", handler.Create)
}

func setupBankRoutes() {
	handler := handler.NewBankHandler()

	router.
		Group("banks").
		GET("", middleware.Auth(), handler.MyBanks).
		GET("modal-create", middleware.Auth(), handler.CreateModal).
		PUT("", middleware.Auth(), handler.CreateBank).
		GET(":username/:slug", middleware.Auth(), handler.ViewBank)
}

func setupCustomerRoutes() {
	handler := handler.NewCustomerHandler()

	router.
		Group("customers").
		GET("modal-create", middleware.Auth(), handler.OpenCreateModal).
		PUT("", middleware.Auth(), handler.CreateCustomer)
}
