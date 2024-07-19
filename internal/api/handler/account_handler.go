package handler

import (
	"fmt"
	"net/http"
	"strconv"

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
}

func NewAccountHandler() accountHandler {
	return accountHandler{
		SignedIn:           true,
		ModalType:          "",
		Form:               NewFormData(),
		Account:            domain.Account{},
		accountService:     service.NewAccountService(),
		transactionService: service.NewTransactionService(),
	}
}

func (ah accountHandler) Get(c *gin.Context) {
	accountID := c.Param("id")

	if err := ah.accountService.FindByID(accountID, &ah.Account); err != nil {
		// handle error
	}

	c.HTML(http.StatusOK, "account", ah)
}

func (ah accountHandler) OpenSettings(c *gin.Context) {
	accountID := c.Param("id")
	ah.Form = NewFormData()
	ah.ModalType = "account_settings"

	if err := ah.accountService.FindByID(accountID, &ah.Account); err != nil {
		// TODO: handle error
	}

	ah.Form.Data["name"] = ah.Account.Name

	c.HTML(http.StatusOK, "modal", ah)
}

func (ah accountHandler) OpenWithdrawOrDeposit(c *gin.Context) {
	accountID := c.Param("id")
	ah.Form = NewFormData()
	ah.ModalType = "withdraw_or_deposit_modal"

	if err := ah.accountService.FindByID(accountID, &ah.Account); err != nil {
		// TODO: handle error
	}

	c.HTML(http.StatusOK, "modal", ah)
}

func (ah accountHandler) CashFlow(c *gin.Context) {
	var cashflow service.Cashflow

	if err := ah.transactionService.CashflowByAccount(c.Param("id"), &cashflow); err != nil {
		fmt.Println(err)
	}

	c.HTML(http.StatusOK, "chart_deposits_vs_withdrawals", cashflow)
}

func (ah accountHandler) WithdrawOrDeposit(c *gin.Context) {
	ah.Form = GetForm(c)
	accountID, _ := strconv.Atoi(c.Param("id"))
	amount, _ := strconv.ParseFloat(ah.Form.Data["amount"], 64)
	userID, _ := utils.ConvertToUintPointer(c.GetString("user_id"))

	if ah.Form.Data["type"] == "withdraw" {
		amount = amount * -1
	}

	transaction := domain.Transaction{
		AccountID:   uint(accountID),
		Amount:      amount,
		Description: ah.Form.Data["description"],
		UserID:      userID,
	}

	if err := ah.transactionService.Create(&transaction); err != nil {
		// TODO: handle the error
	}

	c.Header("HX-Trigger", "closeModal")
	c.HTML(http.StatusOK, "account", ah)
}

func (ah accountHandler) Update(c *gin.Context) {
	accountID := c.Param("id")
	ah.Form = GetForm(c)

	if err := ah.accountService.FindByID(accountID, &ah.Account); err != nil {
		// TODO: handle error
	}

	ah.Account.Name = ah.Form.Data["name"]
	if err := ah.accountService.Update(&ah.Account); err != nil {
		// TODO: handle error
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
		// handle error
	}

	if err := ah.transactionService.CountAllByAccount(accountID, &pagingInfo.TotalItems); err != nil {
		// handle error
	}

	c.HTML(http.StatusOK, "transactions", struct {
		PagingInfo pagination.PagingInfo[domain.Transaction]
		AccountID  string
	}{pagingInfo, accountID})
}
