package api

import (
	"html/template"
	"os"
	"time"

	"github.com/gin-contrib/gzip"

	"github.com/bytebury/fun-banking/internal/api/handler"
	"github.com/bytebury/fun-banking/internal/api/middleware"
	"github.com/bytebury/fun-banking/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var router = gin.Default()

func Start() {
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.SetFuncMap(template.FuncMap{
		"html":     func(text string) template.HTML { return template.HTML(text) },
		"titleize": func(text string) string { return cases.Title(language.AmericanEnglish).String(text) },
		"number":   func(amount int64) string { return utils.FormatNumber(amount) },
		"currency": func(amount float64) string { return utils.FormatCurrency(amount) },
		"sub":      func(a, b int) int { return a - b },
		"add":      func(a, b int) int { return a + b },
		"mul":      func(a, b int) int { return a * b },
		"mulfloat": func(a, b float64) float64 { return a * b },
		"datetime": func(dateTime time.Time) string { return dateTime.Format("January 02, 2006 at 3:04 PM") },
		"date":     func(date time.Time) string { return date.Format("January 02, 2006") },
	})
	// Load Templates
	router.LoadHTMLGlob("templates/**/*")
	// Load Static Files
	router.Static("/static", "public/")
	// Middleware
	router.Use(middleware.Audit(), middleware.CustomerAudit(), middleware.PreferencesAudit())
	// Setup Routes
	setupHomepageRoutes()
	setupAppDrawerRoutes()
	setupUserRoutes()
	setupSessionRoutes()
	setupBankRoutes()
	setupCustomerRoutes()
	setupAccountRoutes()
	setupTransactionRoutes()
	setupNotificationRoutes()
	setupControlPanelRoutes()
	setupAnnouncementRoutes()
	// Run the application
	router.Run()
}

func setupHomepageRoutes() {
	handler := handler.NewHomePageHandler()

	router.
		GET("", handler.Homepage).
		GET("terms", handler.TermsOfService).
		GET("privacy", handler.PrivacyPolicy).
		GET("verify-account", middleware.NoAuth(), handler.VerifyEmail).
		POST("verify-account", middleware.NoAuth(), handler.ResendVerifyEmail).
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
	router.GET("preferences", middleware.UserAuth(), handler.Preferences)
	router.PATCH("preferences", middleware.UserAuth(), handler.UpdatePreferences)
	router.GET("forgot", handler.ForgotPassword)
	router.POST("forgot", handler.SendForgotPasswordEmail)
	router.GET("reset-password", handler.ResetPassword)
	router.POST("reset-password", handler.UpdatePassword)

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
	bank := handler.NewBankHandler()

	router.
		Group("banks").
		GET("", middleware.UserAuth(), bank.MyBanks).
		PUT("", middleware.UserAuth(), bank.CreateBank).
		POST("create", middleware.UserAuth(), bank.OpenCreateModal).
		GET(":id", middleware.UserAuth(), bank.ViewBank).
		PATCH(":id", middleware.UserAuth(), bank.UpdateBank).
		DELETE(":id", middleware.UserAuth(), bank.Delete).
		POST(":id/settings", middleware.UserAuth(), bank.OpenSettingsModal).
		GET(":id/customers", middleware.AnyAuth(), bank.CustomerSearch).
		GET(":id/customers-filter", middleware.UserAuth(), bank.FilterCustomers).
		POST(":id/create-customer", middleware.UserAuth(), bank.OpenCreateCustomerModal).
		PUT(":id/create-customer", middleware.UserAuth(), bank.CreateCustomer)
}

func setupCustomerRoutes() {
	handler := handler.NewCustomerHandler()

	router.
		Group("customers").
		GET(":id", middleware.AnyAuth(), handler.GetCustomer).
		PATCH(":id", middleware.UserAuth(), handler.Update).
		DELETE(":id", middleware.UserAuth(), handler.Delete).
		POST(":id/settings", middleware.UserAuth(), handler.OpenSettingsModal)
}

func setupAccountRoutes() {
	account := handler.NewAccountHandler()

	router.
		Group("accounts").
		GET(":id", middleware.AnyAuth(), account.Get).
		PATCH(":id", middleware.UserAuth(), account.Update).
		GET(":id/transactions", middleware.AnyAuth(), account.GetTransactions).
		POST(":id/settings", middleware.UserAuth(), account.OpenSettingsModal).
		GET(":id/cash-flow", middleware.AnyAuth(), account.CashFlow).
		POST(":id/withdraw-or-deposit", middleware.AnyAuth(), account.OpenWithdrawOrDepositModal).
		PUT(":id/withdraw-or-deposit", middleware.AnyAuth(), account.WithdrawOrDeposit).
		GET(":id/send-money", middleware.AnyAuth(), account.OpenSendMoneyModal).
		PUT(":id/send-money", middleware.AnyAuth(), account.SendMoney).
		GET(":id/statements", middleware.AnyAuth(), account.Statements).
		POST(":id/auto-pay", middleware.UserAuth(), account.OpenAutoPayModal).
		GET(":id/auto-pay", middleware.UserAuth(), account.AutoPay).
		PUT(":id/auto-pay", middleware.UserAuth(), account.CreateAutoPay).
		PATCH(":id/auto-pay/:auto_pay_id", middleware.UserAuth(), account.UpdateAutoPay)
}

func setupTransactionRoutes() {
	handler := handler.NewTransactionHandler()

	router.
		Group("transactions").
		PUT("", middleware.AnyAuth(), handler.Create).
		PATCH(":id/approve", middleware.UserAuth(), handler.Approve).
		PATCH(":id/decline", middleware.UserAuth(), handler.Decline).
		GET("open-bulk-transfer", middleware.UserAuth(), handler.OpenBulkTransferModal).
		PUT("bulk", middleware.UserAuth(), handler.BulkTransfer)
}

func setupAppDrawerRoutes() {
	appDrawer := handler.NewAppDrawerHandler()

	router.Group("app-drawer").
		POST("open", appDrawer.Open)
}

func setupControlPanelRoutes() {
	controlPanel := handler.NewControlPanelHandler()
	announcements := handler.NewAnnouncementHandler()

	router.
		Group("control-panel", middleware.UserAuth(), middleware.AdminOnly()).
		GET("", controlPanel.AppInsights).
		GET("users", controlPanel.GetUsers).
		GET("users/:id", controlPanel.OpenUserModal).
		GET("users/search", controlPanel.SearchUsers).
		GET("announcements", announcements.Dashboard).
		GET("announcements/:id", announcements.Edit).
		PUT("announcements", announcements.Create).
		PATCH("announcements/:id", announcements.Update).
		DELETE("announcements/:id", announcements.Destroy).
		GET("polls", controlPanel.Polls)
}

func setupAnnouncementRoutes() {
	announcements := handler.NewAnnouncementHandler()

	router.Group("announcements").
		GET("", middleware.UserAuth(), announcements.FindAll).
		GET(":id", middleware.UserAuth(), announcements.FindByID).
		POST("recent", middleware.UserAuth(), announcements.RecentAnnouncements)
}
