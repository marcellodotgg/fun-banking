package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/pagination"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/bytebury/fun-banking/internal/utils"
	"github.com/gin-gonic/gin"
)

type accountHandler struct {
	SignedIn           bool
	ModalType          string
	Form               FormData
	Account            domain.Account
	accountService     service.AccountService
	transactionService service.TransactionService
	customerService    service.CustomerService
}

func NewAccountHandler() accountHandler {
	return accountHandler{
		SignedIn:           true,
		ModalType:          "",
		Form:               NewFormData(),
		Account:            domain.Account{},
		accountService:     service.NewAccountService(),
		transactionService: service.NewTransactionService(),
		customerService:    service.NewCustomerService(),
	}
}

func (ah accountHandler) Get(c *gin.Context) {
	accountID := c.Param("id")

	if err := ah.accountService.FindByID(accountID, &ah.Account); err != nil {
		c.HTML(http.StatusNotFound, "not-found", ah)
		return
	}

	c.HTML(http.StatusOK, "account", ah)
}

func (ah accountHandler) OpenSettingsModal(c *gin.Context) {
	accountID := c.Param("id")
	ah.Form = NewFormData()
	ah.ModalType = "account_settings"
	ah.accountService.FindByID(accountID, &ah.Account)
	ah.Form.Data["name"] = ah.Account.Name

	c.HTML(http.StatusOK, "modal", ah)
}

func (ah accountHandler) OpenWithdrawOrDepositModal(c *gin.Context) {
	accountID := c.Param("id")
	ah.Form = NewFormData()
	ah.ModalType = "withdraw_or_deposit_modal"
	ah.accountService.FindByID(accountID, &ah.Account)

	c.HTML(http.StatusOK, "modal", ah)
}

func (ah accountHandler) WithdrawOrDeposit(c *gin.Context) {
	ah.Form = GetForm(c)
	accountID, _ := strconv.Atoi(c.Param("id"))
	amount, _ := strconv.ParseFloat(ah.Form.Data["amount"], 64)
	userID, _ := utils.ConvertToIntPointer(c.GetString("user_id"))

	if ah.Form.Data["type"] == "withdraw" {
		amount = amount * -1
	}

	transaction := domain.Transaction{
		AccountID:   accountID,
		Amount:      amount,
		Description: ah.Form.Data["description"],
		UserID:      userID,
	}

	if err := ah.transactionService.Create(&transaction); err != nil {
		if strings.Contains(err.Error(), "cannot be 0") {
			ah.Form.Errors["general"] = "Please fix the fields marked with errors"
			ah.Form.Errors["amount"] = "Amount cannot be 0"
			c.HTML(http.StatusUnprocessableEntity, "withdraw_or_deposit_form", ah)
			return
		}
		ah.Form.Errors["general"] = "Something happened trying to create that transaction"
		c.HTML(http.StatusUnprocessableEntity, "withdraw_or_deposit_form", ah)
		return
	}

	ah.accountService.FindByID(c.Param("id"), &ah.Account)

	c.Header("HX-Trigger", "closeModal")
	c.HTML(http.StatusOK, "account_oob", ah)
}

func (ah accountHandler) CashFlow(c *gin.Context) {
	var cashflow service.Cashflow

	if err := ah.transactionService.CashflowByAccount(c.Param("id"), &cashflow); err != nil {
		fmt.Println(err)
	}

	c.HTML(http.StatusOK, "chart_deposits_vs_withdrawals", cashflow)
}

func (ah accountHandler) Update(c *gin.Context) {
	ah.Form = GetForm(c)

	ah.Account.Name = ah.Form.Data["name"]
	if err := ah.accountService.Update(c.Param("id"), &ah.Account); err != nil {
		ah.Form.Errors["general"] = "Something happened trying to update your account"
		c.HTML(http.StatusUnprocessableEntity, "account_settings_form", ah)
		return
	}

	c.HTML(http.StatusOK, "account_settings_oob", ah)
}

func (ah accountHandler) GetTransactions(c *gin.Context) {
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

	if err := ah.transactionService.FindAllByAccount(accountID, &pagingInfo.Items, pagingInfo); err != nil {
		c.HTML(http.StatusNotFound, "not-found", ah)
		return
	}

	if err := ah.transactionService.CountAllByAccount(accountID, &pagingInfo.TotalItems); err != nil {
		c.HTML(http.StatusNotFound, "not-found", ah)
		return
	}

	c.HTML(http.StatusOK, "transactions", struct {
		PagingInfo pagination.PagingInfo[domain.Transaction]
		AccountID  string
	}{pagingInfo, accountID})
}

func (h accountHandler) OpenSendMoneyModal(c *gin.Context) {
	accountID := c.Param("id")
	h.ModalType = "send_money_modal"

	if err := h.accountService.FindByID(accountID, &h.Account); err != nil {
		c.HTML(http.StatusNotFound, "not-found", nil)
		return
	}

	c.HTML(http.StatusOK, "modal", h)
}

func (h accountHandler) SendMoney(c *gin.Context) {
	accountID := c.Param("id")
	h.Form = GetForm(c)
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

	amount, _ := strconv.ParseFloat(h.Form.Data["amount"], 64)

	transaction := domain.Transaction{
		Amount:      amount,
		Description: h.Form.Data["description"],
	}

	if err := h.transactionService.SendMoney(h.Account, recipient, &transaction); err != nil {
		// TODO: handl the error
	}

	c.Header("HX-Redirect", "/accounts/"+accountID)
}
