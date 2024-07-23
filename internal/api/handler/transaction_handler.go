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
	SignedIn           bool
	Form               FormData
	ModalType          string
	transactionService service.TransactionService
	accountService     service.AccountService
	customerService    service.CustomerService
	userService        service.UserService
	Customer           domain.Customer
}

func NewTransactionHandler() transactionHandler {
	return transactionHandler{
		SignedIn:           true,
		Form:               NewFormData(),
		ModalType:          "",
		Customer:           domain.Customer{},
		accountService:     service.NewAccountService(),
		customerService:    service.NewCustomerService(),
		userService:        service.NewUserService(),
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

	userID, _ := utils.ConvertToIntPointer(c.GetString("user_id"))

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

func (h transactionHandler) Approve(c *gin.Context) {
	transactionID := c.Param("id")

	if err := h.transactionService.Update(transactionID, c.GetString("user_id"), domain.TransactionApproved); err != nil {
		c.HTML(http.StatusBadRequest, "", h)
		return
	}

	var transactions []domain.Transaction
	h.userService.FindPendingTransactions(c.GetString("user_id"), &transactions)
	c.HTML(http.StatusAccepted, "notifications_list_oob", transactions)

}

func (h transactionHandler) Decline(c *gin.Context) {
	transactionID := c.Param("id")

	if err := h.transactionService.Update(transactionID, c.GetString("user_id"), domain.TransactionDeclined); err != nil {
		c.HTML(http.StatusBadRequest, "", h)
		return
	}

	var transactions []domain.Transaction
	h.userService.FindPendingTransactions(c.GetString("user_id"), &transactions)
	c.HTML(http.StatusAccepted, "notifications_list_oob", transactions)
}

func (h transactionHandler) OpenBulkTransferModal(c *gin.Context) {
	h.ModalType = "bulk_transfer_modal"
	h.Form = GetForm(c)
	h.Form.Data["customer_ids"] = strings.Join(c.QueryArray("ids"), ",")
	c.HTML(http.StatusOK, "modal", h)
}

func (h transactionHandler) BulkTransfer(c *gin.Context) {
	h.Form = GetForm(c)
	customerIDs := strings.Split(h.Form.Data["customer_ids"], ",")

	amount, _ := strconv.ParseFloat(h.Form.Data["amount"], 64)
	userID, _ := utils.ConvertToIntPointer(c.GetString("user_id"))

	transaction := domain.Transaction{
		Amount:      h.getTransferAmount(amount, h.Form.Data["type"]),
		Description: h.Form.Data["description"],
		UserID:      userID,
	}

	if err := h.transactionService.BulkTransfer(customerIDs, &transaction); err != nil {
		c.HTML(http.StatusUnprocessableEntity, "bulk_transfer_form", h)
		return
	}

	c.Header("HX-Trigger", "closeModal")
	// TODO: OOB swap for this, too.
}

func (th transactionHandler) getTransferAmount(amount float64, transferType string) float64 {
	if transferType == "withdraw" {
		return amount * -1
	}
	return amount
}
