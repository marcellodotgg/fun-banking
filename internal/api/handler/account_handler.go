package handler

import (
	"net/http"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type accountHandler struct {
	SignedIn       bool
	Account        domain.Account
	accountService service.AccountService
}

func NewAccountHandler() accountHandler {
	return accountHandler{
		SignedIn:       true,
		Account:        domain.Account{},
		accountService: service.NewAccountService(),
	}
}

func (ah accountHandler) Get(c *gin.Context) {
	accountID := c.Param("id")

	if err := ah.accountService.FindByID(accountID, &ah.Account); err != nil {
		// handle error
	}

	c.HTML(http.StatusOK, "account", ah)
}
