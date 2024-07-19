package handler

import (
	"net/http"
	"strconv"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/bytebury/fun-banking/internal/utils"
	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	SignedIn           bool
	Form               FormData
	transactionService service.TransactionService
	accountService     service.AccountService
	customerService    service.CustomerService
	Customer           domain.Customer
}

func NewTransactionHandler() transactionHandler {
	return transactionHandler{
		SignedIn:           true,
		Form:               NewFormData(),
		Customer:           domain.Customer{},
		accountService:     service.NewAccountService(),
		customerService:    service.NewCustomerService(),
		transactionService: service.NewTransactionService(),
	}
}

func (th transactionHandler) Create(c *gin.Context) {
	th.Form = GetForm(c)

	var account domain.Account
	if err := th.accountService.FindByID(th.Form.Data["account_id"], &account); err != nil {
		th.Form.Errors["general"] = "Something went wrong creating your transaction"
		c.HTML(http.StatusUnprocessableEntity, "transfer_money_form", th)
		return
	}

	th.Customer = account.Customer

	amount, err := strconv.ParseFloat(th.Form.Data["amount"], 64)

	if err != nil {
		th.Form.Errors["amount"] = "Amount is not a valid number"
		c.HTML(http.StatusUnprocessableEntity, "transfer_money_form", th.Customer)
		return
	}

	userID, _ := utils.ConvertToUintPointer(c.GetString("user_id"))

	transaction := domain.Transaction{
		AccountID:   account.ID,
		Amount:      th.getTransferAmount(amount, th.Form.Data["type"]),
		Description: th.Form.Data["description"],
		Status:      domain.TransactionPending,
		UserID:      userID,
	}

	if err := th.transactionService.Create(&transaction); err != nil {
		c.HTML(http.StatusUnprocessableEntity, "transfer_money_form", th.Customer)
		return
	}

	th.customerService.FindByID(strconv.Itoa(int(th.Customer.ID)), &th.Customer)

	c.HTML(http.StatusOK, "transfer_money_form_oob", th.Customer)
}

func (th transactionHandler) getTransferAmount(amount float64, transferType string) float64 {
	if transferType == "withdraw" {
		return amount * -1
	}
	return amount
}
