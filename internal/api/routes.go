package api

import (
	"html/template"
	"time"

	"github.com/bytebury/fun-banking/internal/api/handler"
	"github.com/bytebury/fun-banking/internal/api/middleware"
	"github.com/bytebury/fun-banking/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var router = gin.Default()

func Start() {
	router.SetFuncMap(template.FuncMap{
		"titleize": func(text string) string { return cases.Title(language.AmericanEnglish).String(text) },
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
	router.Use(middleware.Audit(), middleware.CustomerAudit())
	// Setup Routes
	setupHomepageRoutes()
	setupActionRoutes()
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
		GET("", handler.Homepage).
		GET(":username/:slug", middleware.NoAuth(), handler.BankSignIn)
}

func setupSessionRoutes() {
	handler := handler.NewSessionHandler()

	router.
		GET("signin", middleware.NoAuth(), handler.SignIn).
		POST("signin", middleware.NoAuth(), handler.CreateSession).
		DELETE("signout", handler.DestroySession).
		POST("sessions/customer", middleware.NoAuth(), handler.CreateCustomerSession).
		DELETE("sessions/customer", handler.DestroyCustomerSession)
}

func setupUserRoutes() {
	handler := handler.NewUserHandler()

	router.GET("signup", middleware.NoAuth(), handler.SignUp)
	router.GET("settings", middleware.Auth(), handler.Settings)

	router.
		Group("users").
		PUT("", handler.Create).
		PATCH("", middleware.Auth(), handler.Update)
}

func setupBankRoutes() {
	handler := handler.NewBankHandler()

	router.
		Group("banks").
		GET("", middleware.Auth(), handler.MyBanks).
		PUT("", middleware.Auth(), handler.CreateBank).
		POST("create", middleware.Auth(), handler.OpenCreateModal).
		GET(":id", middleware.Auth(), handler.ViewBank).
		PATCH(":id", middleware.Auth(), handler.UpdateBank).
		POST(":id/settings", middleware.Auth(), handler.OpenSettingsModal).
		GET(":id/customers", middleware.Auth(), handler.CustomerSearch).
		GET(":id/customers-filter", middleware.Auth(), handler.FilterCustomers).
		POST(":id/create-customer", middleware.Auth(), handler.OpenCreateCustomerModal).
		PUT(":id/create-customer", middleware.Auth(), handler.CreateCustomer)
}

func setupCustomerRoutes() {
	handler := handler.NewCustomerHandler()

	router.
		Group("customers").
		GET(":id", middleware.Auth(), handler.GetCustomer).
		PATCH(":id", middleware.Auth(), handler.Update).
		POST(":id/settings", middleware.Auth(), handler.OpenSettingsModal)
}

func setupAccountRoutes() {
	handler := handler.NewAccountHandler()

	router.
		Group("accounts").
		GET(":id", middleware.Auth(), handler.Get).
		PATCH(":id", middleware.Auth(), handler.Update).
		GET(":id/transactions", middleware.Auth(), handler.GetTransactions).
		POST(":id/settings", middleware.Auth(), handler.OpenSettingsModal).
		GET(":id/cash-flow", middleware.Auth(), handler.CashFlow).
		POST(":id/withdraw-or-deposit", middleware.Auth(), handler.OpenWithdrawOrDepositModal).
		PUT(":id/withdraw-or-deposit", middleware.Auth(), handler.WithdrawOrDeposit).
		GET(":id/send-money", middleware.Auth(), handler.OpenSendMoneyModal).
		PUT(":id/send-money", middleware.Auth(), handler.SendMoney)
}

func setupTransactionRoutes() {
	handler := handler.NewTransactionHandler()

	router.
		Group("transactions").
		// TODO creation should actually be on the account
		PUT("", middleware.Auth(), handler.Create)
}

func setupActionRoutes() {
	handler := handler.NewActionHandler()

	router.Group("actions").
		GET("open-app-drawer", handler.OpenAppDrawer)
}
