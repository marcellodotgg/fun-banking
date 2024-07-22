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
		"number":   func(amount int64) string { return utils.FormatNumber(amount) },
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
	setupNotificationRoutes()
	setupControlPanelRoutes()
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
		DELETE("signout", middleware.UserAuth(), handler.DestroySession).
		POST("sessions/customer", middleware.NoAuth(), handler.CreateCustomerSession).
		DELETE("sessions/customer", middleware.CustomerAuth(), handler.DestroyCustomerSession)
}

func setupUserRoutes() {
	handler := handler.NewUserHandler()

	router.GET("signup", middleware.NoAuth(), handler.SignUp)
	router.GET("settings", middleware.UserAuth(), handler.Settings)

	router.
		Group("users").
		PUT("", handler.Create).
		PATCH("", middleware.UserAuth(), handler.Update)
}

func setupNotificationRoutes() {
	handler := handler.NewUserHandler()

	router.
		Group("notifications").
		GET("", middleware.UserAuth(), handler.Notifications).
		GET("pending", middleware.UserAuth(), handler.PendingTransactions).
		POST("has-pending", middleware.UserAuth(), handler.HasPendingTransactions)
}

func setupBankRoutes() {
	handler := handler.NewBankHandler()

	router.
		Group("banks").
		GET("", middleware.UserAuth(), handler.MyBanks).
		PUT("", middleware.UserAuth(), handler.CreateBank).
		POST("create", middleware.UserAuth(), handler.OpenCreateModal).
		GET(":id", middleware.UserAuth(), handler.ViewBank).
		PATCH(":id", middleware.UserAuth(), handler.UpdateBank).
		POST(":id/settings", middleware.UserAuth(), handler.OpenSettingsModal).
		GET(":id/customers", middleware.AnyAuth(), handler.CustomerSearch).
		GET(":id/customers-filter", middleware.UserAuth(), handler.FilterCustomers).
		POST(":id/create-customer", middleware.UserAuth(), handler.OpenCreateCustomerModal).
		PUT(":id/create-customer", middleware.UserAuth(), handler.CreateCustomer)
}

func setupCustomerRoutes() {
	handler := handler.NewCustomerHandler()

	router.
		Group("customers").
		GET(":id", middleware.AnyAuth(), handler.GetCustomer).
		PATCH(":id", middleware.UserAuth(), handler.Update).
		POST(":id/settings", middleware.UserAuth(), handler.OpenSettingsModal)
}

func setupAccountRoutes() {
	handler := handler.NewAccountHandler()

	router.
		Group("accounts").
		GET(":id", middleware.AnyAuth(), handler.Get).
		PATCH(":id", middleware.UserAuth(), handler.Update).
		GET(":id/transactions", middleware.AnyAuth(), handler.GetTransactions).
		POST(":id/settings", middleware.UserAuth(), handler.OpenSettingsModal).
		GET(":id/cash-flow", middleware.AnyAuth(), handler.CashFlow).
		POST(":id/withdraw-or-deposit", middleware.AnyAuth(), handler.OpenWithdrawOrDepositModal).
		PUT(":id/withdraw-or-deposit", middleware.AnyAuth(), handler.WithdrawOrDeposit).
		GET(":id/send-money", middleware.AnyAuth(), handler.OpenSendMoneyModal).
		PUT(":id/send-money", middleware.AnyAuth(), handler.SendMoney)
}

func setupTransactionRoutes() {
	handler := handler.NewTransactionHandler()

	router.
		Group("transactions").
		// TODO creation should actually be on the account
		PUT("", middleware.AnyAuth(), handler.Create).
		PATCH(":id/approve", middleware.UserAuth(), handler.Approve).
		PATCH(":id/decline", middleware.UserAuth(), handler.Decline)
}

func setupActionRoutes() {
	handler := handler.NewActionHandler()

	router.Group("actions").
		GET("open-app-drawer", handler.OpenAppDrawer)
}

func setupControlPanelRoutes() {
	handler := handler.NewControlPanelHandler()

	router.
		Group("control-panel", middleware.UserAuth(), middleware.AdminOnly()).
		GET("", handler.AppInsights).
		GET("users", handler.Users).
		GET("announcements", handler.Announcements).
		GET("polls", handler.Polls)
}
