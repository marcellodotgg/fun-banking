package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/bytebury/fun-banking/internal/utils"
	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	pageObject
	transactionService service.TransactionService
	accountService     service.AccountService
	customerService    service.CustomerService
	userService        service.UserService
	bankService        service.BankService
	Customer           domain.Customer
	Bank               domain.Bank
}

func NewTransactionHandler() transactionHandler {
	return transactionHandler{
		Customer:           domain.Customer{},
		Bank:               domain.Bank{},
		accountService:     service.NewAccountService(),
		customerService:    service.NewCustomerService(),
		bankService:        service.NewBankService(),
		userService:        service.NewUserService(),
		transactionService: service.NewTransactionService(),
	}
}

func (h transactionHandler) Create(c *gin.Context) {
	h.Reset(c)

	var account domain.Account
	if err := h.accountService.FindByID(h.Form.Data["account_id"], &account); err != nil {
		h.Form.Errors["general"] = "Something went wrong creating your transaction"
		c.HTML(http.StatusUnprocessableEntity, "transfer_money_form", h)
		return
	}

	h.Customer = account.Customer

	amount, err := utils.GetDollarAmount(h.Form.Data["amount"])

	if err != nil {
		h.Form.Errors["amount"] = "Amount is not a valid number"
		c.HTML(http.StatusUnprocessableEntity, "transfer_money_form", h.Customer)
		return
	}

	userID, _ := utils.ConvertToIntPointer(c.GetString("user_id"))

	transaction := domain.Transaction{
		AccountID:   account.ID,
		Amount:      h.getTransferAmount(amount, h.Form.Data["type"]),
		Description: h.Form.Data["description"],
		Status:      domain.TransactionPending,
		UserID:      userID,
	}

	if err := h.transactionService.Create(&transaction); err != nil {
		c.HTML(http.StatusUnprocessableEntity, "transfer_money_form", h.Customer)
		return
	}

	h.customerService.FindByID(strconv.Itoa(int(h.Customer.ID)), &h.Customer)

	c.HTML(http.StatusOK, "transfer_money_form_oob", h.Customer)
}

func (h transactionHandler) Approve(c *gin.Context) {
	if err := h.transactionService.Update(c.Param("id"), c.GetString("user_id"), domain.TransactionApproved); err != nil {
		c.HTML(http.StatusBadRequest, "", h)
		return
	}

	var transactions []domain.Transaction
	h.userService.FindPendingTransactions(c.GetString("user_id"), &transactions)
	c.HTML(http.StatusAccepted, "notifications_list_oob", transactions)

}

func (h transactionHandler) Decline(c *gin.Context) {
	if err := h.transactionService.Update(c.Param("id"), c.GetString("user_id"), domain.TransactionDeclined); err != nil {
		c.HTML(http.StatusBadRequest, "", h)
		return
	}

	var transactions []domain.Transaction
	h.userService.FindPendingTransactions(c.GetString("user_id"), &transactions)
	c.HTML(http.StatusAccepted, "notifications_list_oob", transactions)
}

func (h transactionHandler) OpenBulkTransferModal(c *gin.Context) {
	h.Reset(c)

	h.ModalType = "bulk_transfer_modal"
	h.Form.Data["customer_ids"] = strings.Join(c.QueryArray("ids"), ",")

	c.HTML(http.StatusOK, "modal", h)
}

func (h transactionHandler) BulkTransfer(c *gin.Context) {
	h.Reset(c)

	customerIDs := strings.Split(h.Form.Data["customer_ids"], ",")

	amount, _ := utils.GetDollarAmount(h.Form.Data["amount"])
	userID, _ := utils.ConvertToIntPointer(c.GetString("user_id"))

	transaction := domain.Transaction{
		Amount:      h.getTransferAmount(amount, h.Form.Data["type"]),
		Description: h.Form.Data["description"],
		UserID:      userID,
	}

	if len(customerIDs) <= 0 {
		c.HTML(http.StatusUnprocessableEntity, "bulk_transfer_form", h)
		return
	}

	// TODO - this should really all be transactional
	if err := h.transactionService.BulkTransfer(customerIDs, &transaction); err != nil {
		c.HTML(http.StatusUnprocessableEntity, "bulk_transfer_form", h)
		return
	}

	if err := h.customerService.FindByID(customerIDs[0], &h.Customer); err != nil {
		h.Form.Errors["general"] = "Something went wrong finding your bank"
		c.HTML(http.StatusUnprocessableEntity, "bulk_transfer_form", h)
		return
	}

	if err := h.bankService.FindByID(strconv.Itoa(h.Customer.BankID), &h.Bank); err != nil {
		c.HTML(http.StatusUnprocessableEntity, "bulk_transfer_form", h)
		return
	}

	c.Header("HX-Trigger", "closeModal")
	c.HTML(http.StatusAccepted, "customers_oob", h)
}

func (h transactionHandler) getTransferAmount(amount float64, transferType string) float64 {
	if transferType == "withdraw" {
		return amount * -1
	}
	return amount
}
