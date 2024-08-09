package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/pagination"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/bytebury/fun-banking/internal/utils"
	"github.com/gin-gonic/gin"
)

type accountHandler struct {
	pageObject
	Account            domain.Account
	PagingInfo         pagination.PagingInfo[domain.Transaction]
	StatementPeriod    string
	LastTwelveMonths   [][]string
	accountService     service.AccountService
	transactionService service.TransactionService
	customerService    service.CustomerService
	userService        service.UserService
}

func NewAccountHandler() accountHandler {
	return accountHandler{
		Account:            domain.Account{},
		PagingInfo:         pagination.PagingInfo[domain.Transaction]{},
		StatementPeriod:    "",
		LastTwelveMonths:   make([][]string, 0),
		accountService:     service.NewAccountService(),
		transactionService: service.NewTransactionService(),
		customerService:    service.NewCustomerService(),
		userService:        service.NewUserService(),
	}
}

func (h accountHandler) Get(c *gin.Context) {
	h.Reset(c)

	accountID := c.Param("id")
	customerID, _ := strconv.Atoi(c.GetString("customer_id"))
	userID := c.GetString("user_id")

	if err := h.accountService.FindByID(accountID, &h.Account); err != nil {
		c.HTML(http.StatusNotFound, "not-found", h)
		return
	}

	if h.Account.CustomerID != customerID && !h.isOwner(accountID, userID) {
		c.HTML(http.StatusForbidden, "forbidden", h)
		return
	}

	h.SignedIn = userID != ""

	c.HTML(http.StatusOK, "account", h)
}

func (h accountHandler) OpenSettingsModal(c *gin.Context) {
	h.Reset(c)

	h.ModalType = "account_settings"
	h.accountService.FindByID(c.Param("id"), &h.Account)
	h.Form.Data["name"] = h.Account.Name

	c.HTML(http.StatusOK, "modal", h)
}

func (h accountHandler) OpenWithdrawOrDepositModal(c *gin.Context) {
	h.Reset(c)
	h.ModalType = "withdraw_or_deposit_modal"
	h.accountService.FindByID(c.Param("id"), &h.Account)

	c.HTML(http.StatusOK, "modal", h)
}

func (h accountHandler) WithdrawOrDeposit(c *gin.Context) {
	h.Reset(c)

	accountID, _ := strconv.Atoi(c.Param("id"))
	amount, _ := utils.GetDollarAmount(h.Form.Data["amount"])
	userID, _ := utils.ConvertToIntPointer(c.GetString("user_id"))

	if h.Form.Data["type"] == "withdraw" {
		amount = amount * -1
	}

	transaction := domain.Transaction{
		AccountID:   accountID,
		Amount:      amount,
		Description: h.Form.Data["description"],
		UserID:      userID,
	}

	if err := h.transactionService.Create(&transaction); err != nil {
		if strings.Contains(err.Error(), "cannot be 0") {
			h.Form.Errors["general"] = "Please fix the fields marked with errors"
			h.Form.Errors["amount"] = "Amount cannot be 0"
			c.HTML(http.StatusUnprocessableEntity, "withdraw_or_deposit_form", h)
			return
		}
		if strings.Contains(err.Error(), "greater than") {
			h.Form.Errors["general"] = "Please fix the fields marked with errors"
			h.Form.Errors["amount"] = "Amount cannot be greater than 25,000,000"
			c.HTML(http.StatusUnprocessableEntity, "withdraw_or_deposit_form", h)
			return
		}
		h.Form.Errors["general"] = "Something happened trying to create that transaction"
		c.HTML(http.StatusUnprocessableEntity, "withdraw_or_deposit_form", h)
		return
	}

	h.accountService.FindByID(c.Param("id"), &h.Account)

	c.Header("HX-Trigger", "closeModal")
	c.HTML(http.StatusOK, "account_oob", h)
}

func (h accountHandler) CashFlow(c *gin.Context) {
	var cashflow service.Cashflow

	if err := h.accountService.CashFlow(c.Param("id"), &cashflow); err != nil {
		c.HTML(http.StatusOK, "chart_deposits_vs_withdrawals", cashflow)
		return
	}

	c.HTML(http.StatusOK, "chart_deposits_vs_withdrawals", cashflow)
}

func (h accountHandler) Update(c *gin.Context) {
	h.Reset(c)

	if !h.isOwner(c.Param("id"), c.GetString("user_id")) {
		h.Form.Errors["general"] = "You don't have access to do that"
		c.HTML(http.StatusForbidden, "account_settings_form", h)
		return
	}

	h.Account.Name = h.Form.Data["name"]
	if err := h.accountService.Update(c.Param("id"), &h.Account); err != nil {
		h.Form.Errors["general"] = "Something happened trying to update your account"
		c.HTML(http.StatusUnprocessableEntity, "account_settings_form", h)
		return
	}

	c.Header("HX-Trigger", "closeModal")
	c.HTML(http.StatusOK, "account_settings_oob", h)
}

func (h accountHandler) GetTransactions(c *gin.Context) {
	h.Reset(c)

	accountID := c.Param("id")
	pageNumber, _ := strconv.Atoi(c.Query("page"))

	if pageNumber < 1 {
		pageNumber = 1
	}

	pagingInfo := pagination.PagingInfo[domain.Transaction]{
		ItemsPerPage: 8,
		PageNumber:   pageNumber,
		TotalItems:   0,
		Items:        nil,
	}

	if err := h.accountService.Transactions(accountID, &pagingInfo); err != nil {
		c.HTML(http.StatusNotFound, "not-found", h)
		return
	}

	c.HTML(http.StatusOK, "transactions", struct {
		PagingInfo pagination.PagingInfo[domain.Transaction]
		AccountID  string
	}{pagingInfo, accountID})
}

func (h accountHandler) OpenSendMoneyModal(c *gin.Context) {
	h.Reset(c)
	h.ModalType = "send_money_modal"

	if err := h.accountService.FindByID(c.Param("id"), &h.Account); err != nil {
		c.HTML(http.StatusNotFound, "not-found", nil)
		return
	}

	c.HTML(http.StatusOK, "modal", h)
}

func (h accountHandler) SendMoney(c *gin.Context) {
	accountID := c.Param("id")

	h.Reset(c)
	recipientID := h.Form.Data["recipient"]

	if err := h.accountService.FindByID(accountID, &h.Account); err != nil {
		c.HTML(http.StatusNotFound, "not-found", nil)
		return
	}

	var recipient domain.Customer
	if err := h.customerService.FindByID(recipientID, &recipient); err != nil {
		c.HTML(http.StatusNotFound, "not-found", nil)
		return
	}

	amount, _ := utils.GetDollarAmount(h.Form.Data["amount"])

	transaction := domain.Transaction{
		Amount:      amount,
		Description: h.Form.Data["description"],
	}

	if err := h.transactionService.SendMoney(h.Account, recipient, &transaction); err != nil {
		if strings.Contains(err.Error(), "not enough money") {
			h.Form.Errors["general"] = "You do not have enough money"
			c.HTML(http.StatusUnprocessableEntity, "send_money_form", h)
			return
		}

		if strings.Contains(err.Error(), "cannot be 0") {
			h.Form.Errors["general"] = "Amount cannot be 0"
			c.HTML(http.StatusUnprocessableEntity, "send_money_form", h)
			return
		}

		h.Form.Errors["general"] = "Something went wrong sending money"
		c.HTML(http.StatusUnprocessableEntity, "send_money_form", h)
		return
	}

	c.Header("HX-Redirect", "/accounts/"+accountID)
}

func (h accountHandler) Statements(c *gin.Context) {
	h.Reset(c)

	accountID := c.Param("id")

	if err := h.accountService.FindByID(accountID, &h.Account); err != nil {
		c.HTML(http.StatusNotFound, "not-found", h)
		return
	}

	if !h.isOwner(accountID, c.GetString("user_id")) && !h.isCustomer(accountID, c.GetString("customer_id")) {
		c.HTML(http.StatusForbidden, "forbidden", h)
		return
	}

	pageNumber, err := strconv.Atoi(c.Query("page"))

	if err != nil {
		pageNumber = 1
	}

	h.PagingInfo = pagination.PagingInfo[domain.Transaction]{
		ItemsPerPage: 10,
		PageNumber:   pageNumber,
		TotalItems:   0,
		Items:        nil,
	}

	h.StatementPeriod = c.Query("period")
	if h.StatementPeriod == "" {
		h.StatementPeriod = fmt.Sprintf("%d-%s", time.Now().Year(), utils.ConvertMonthToNumeric(time.Now().Month()))
	}

	h.accountService.TransactionsByPeriod(accountID, h.StatementPeriod, &h.PagingInfo)

	h.LastTwelveMonths = utils.LastTwelveMonths()

	c.HTML(http.StatusOK, "statements", h)
}

func (h accountHandler) isOwner(accountID, userID string) bool {
	var account domain.Account
	if err := h.accountService.FindByID(accountID, &account); err != nil {
		return false
	}

	var user domain.User
	if err := h.userService.FindByID(userID, &user); err != nil {
		return false
	}

	return account.Customer.Bank.UserID == user.ID || user.IsAdmin()
}

func (h accountHandler) isCustomer(accountID, customerID string) bool {
	var account domain.Account
	if err := h.accountService.FindByID(accountID, &account); err != nil {
		return false
	}

	return strconv.Itoa(account.CustomerID) == customerID
}
