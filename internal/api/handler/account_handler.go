package handler

import (
	"net/http"
	"strconv"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/pagination"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type accountHandler struct {
	SignedIn           bool
	Account            domain.Account
	accountService     service.AccountService
	transactionService service.TransactionService
}

func NewAccountHandler() accountHandler {
	return accountHandler{
		SignedIn:           true,
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
