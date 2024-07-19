package api

import (
	"html/template"
	"time"

	"github.com/bytebury/fun-banking/internal/api/handler"
	"github.com/bytebury/fun-banking/internal/api/middleware"
	"github.com/bytebury/fun-banking/internal/utils"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func Start() {
	router.SetFuncMap(template.FuncMap{
		"currency": func(amount float64) string { return utils.FormatCurrency(amount) },
		"sub":      func(a, b int) int { return a - b },
		"add":      func(a, b int) int { return a + b },
		"mul":      func(a, b int) int { return a * b },
		"mulfloat": func(a, b float64) float64 { return a * b },
		"datetime": func(dateTime time.Time) string { return dateTime.Format("July 02, 2006 at 3:04 PM") },
	})
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
	setupAccountRoutes()
	setupTransactionRoutes()
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
		GET(":username/:slug", middleware.Auth(), handler.ViewBank).
		GET("settings", middleware.Auth(), handler.Settings).
		PATCH("", middleware.Auth(), handler.UpdateBank)
}

func setupCustomerRoutes() {
	handler := handler.NewCustomerHandler()

	router.
		Group("customers").
		GET("modal-create", middleware.Auth(), handler.OpenCreateModal).
		GET(":id", middleware.Auth(), handler.GetCustomer).
		GET(":id/settings", middleware.Auth(), handler.OpenSettings).
		PATCH(":id", middleware.Auth(), handler.Update).
		PUT("", middleware.Auth(), handler.CreateCustomer)
}

func setupAccountRoutes() {
	handler := handler.NewAccountHandler()

	router.
		Group("accounts").
		GET(":id", middleware.Auth(), handler.Get).
		GET(":id/transactions", middleware.Auth(), handler.GetTransactions).
		GET(":id/settings", middleware.Auth(), handler.OpenSettings).
		PATCH(":id", middleware.Auth(), handler.Update)
}

func setupTransactionRoutes() {
	handler := handler.NewTransactionHandler()

	router.
		Group("transactions").
		PUT("", middleware.Auth(), handler.Create)
}
