package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type customerHandler struct {
	pageObject
	bankService            service.BankService
	customerService        service.CustomerService
	accountService         service.AccountService
	transactionService     service.TransactionService
	userService            service.UserService
	Bank                   domain.Bank
	Customer               domain.Customer
	MAX_TRANSACTION_AMOUNT int
}

func NewCustomerHandler() customerHandler {
	return customerHandler{
		bankService:            service.NewBankService(),
		customerService:        service.NewCustomerService(),
		userService:            service.NewUserService(),
		accountService:         service.NewAccountService(),
		transactionService:     service.NewTransactionService(),
		Bank:                   domain.Bank{},
		MAX_TRANSACTION_AMOUNT: domain.MAX_TRANSACTION_AMOUNT,
	}
}

func (h customerHandler) OpenCreateModal(c *gin.Context) {
	h.Reset(c)

	h.ModalType = "create_customer_modal"
	h.Form.Data["bank_id"] = c.Query("bank_id")

	c.HTML(http.StatusOK, "modal", h)
}

func (h customerHandler) GetCustomer(c *gin.Context) {
	h.Reset(c)

	isCustomer := c.GetString("customer_id") == c.Param("id")
	if !isCustomer && !h.isOwner(c.Param("id"), c.GetString("user_id")) {
		c.HTML(http.StatusForbidden, "forbidden", h)
		return
	}

	if err := h.customerService.FindByID(c.Param("id"), &h.Customer); err != nil {
		c.HTML(http.StatusNotFound, "not-found", h)
		return
	}

	c.HTML(http.StatusOK, "customer", h)
}

func (h customerHandler) OpenSettingsModal(c *gin.Context) {
	h.Reset(c)

	h.ModalType = "customer_settings"

	h.customerService.FindByID(c.Param("id"), &h.Customer)

	h.Form.Data["first_name"] = h.Customer.FirstName
	h.Form.Data["last_name"] = h.Customer.LastName
	h.Form.Data["pin"] = h.Customer.PIN

	c.HTML(http.StatusOK, "modal", h)
}

func (h customerHandler) Update(c *gin.Context) {
	h.Reset(c)

	if !h.isOwner(c.Param("id"), c.GetString("user_id")) {
		h.Form.Errors["general"] = "You do not have access to do that"
		c.HTML(http.StatusOK, "customer_settings_form", h)
		return
	}

	h.Customer = domain.Customer{
		FirstName: h.Form.Data["first_name"],
		LastName:  h.Form.Data["last_name"],
		PIN:       h.Form.Data["pin"],
	}

	if err := h.customerService.Update(c.Param("id"), &h.Customer); err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			h.Form.Errors["pin"] = "PIN is already in use by another customer"
			c.HTML(http.StatusOK, "customer_settings_form", h)
			return
		}
		if strings.Contains(err.Error(), "invalid PIN") {
			h.Form.Errors["pin"] = "PINs can only be 4 to 6 digits"
			c.HTML(http.StatusOK, "customer_settings_form", h)
			return
		}

		h.Form.Errors["general"] = "Something went wrong updating your customer"
		c.HTML(http.StatusOK, "customer_settings_form", h)
		return
	}

	h.Form.Data["success"] = "Successfully updated the customer"
	c.HTML(http.StatusOK, "customer_settings_oob", h)
}

func (h customerHandler) Delete(c *gin.Context) {
	if !h.isOwner(c.Param("id"), c.GetString("user_id")) {
		h.Form.Errors["general"] = "You do not have access to do that"
		c.HTML(http.StatusUnprocessableEntity, "customer_settings_form", h)
		return
	}

	if err := h.customerService.FindByID(c.Param("id"), &h.Customer); err != nil {
		h.Form.Errors["general"] = "Something went wrong deleting this customer"
		c.HTML(http.StatusUnprocessableEntity, "customer_settings_form", h)
		return
	}

	if err := h.customerService.Delete(c.Param("id")); err != nil {
		h.Form.Errors["general"] = "Something went wrong deleting this customer"
		c.HTML(http.StatusForbidden, "customer_settings_form", h)
		return
	}

	c.Header("HX-Trigger", "closeModal")
	c.Header("HX-Redirect", fmt.Sprintf("/banks/%d", h.Customer.BankID))
}

func (h customerHandler) OpenAccountModal(c *gin.Context) {
	h.Reset(c)
	h.ModalType = "create_account"
	h.customerService.FindByID(c.Param("id"), &h.Customer)

	c.HTML(http.StatusOK, "modal", h)
}

func (h customerHandler) OpenAccount(c *gin.Context) {
	h.Reset(c)

	if err := h.customerService.FindByID(c.Param("id"), &h.Customer); err != nil {
		h.Form.Errors["general"] = "Something went wrong creating the account"
		c.HTML(http.StatusOK, "account/create_form", h)
		return
	}

	account := domain.Account{
		CustomerID: h.Customer.ID,
		Name:       h.Form.Data["name"],
	}

	if err := h.accountService.Create(&account); err != nil {
		if strings.Contains(err.Error(), "maximum") {
			h.Form.Errors["general"] = "You have already met the maximum amount of accounts"
		} else {
			h.Form.Errors["general"] = "Something went wrong creating the account"
		}

		c.HTML(http.StatusOK, "account/create_form", h)
		return
	}

	c.Header("HX-Redirect", fmt.Sprintf("/accounts/%d", account.ID))
}

func (h customerHandler) OpenTransferMoneyModal(c *gin.Context) {
	h.Reset(c)
	h.ModalType = "transfer_money_modal"

	if err := h.customerService.FindByID(c.Param("id"), &h.Customer); err != nil {
		c.HTML(http.StatusNotFound, "not-found", nil)
		return
	}

	c.HTML(http.StatusOK, "modal", h)
}

func (h customerHandler) TransferMoney(c *gin.Context) {
	h.Reset(c)

	if err := h.customerService.FindByID(c.Param("id"), &h.Customer); err != nil {
		h.Form.Errors["general"] = "Could not find customer"
		c.HTML(http.StatusNotFound, "account/transfer-money-form", h)
		return
	}

	var fromAccount domain.Account
	var toAccount domain.Account

	if err := persistence.DB.First(&fromAccount, "id = ?", h.Form.Data["from_account"]).Error; err != nil {
		h.Form.Errors["general"] = "Account does not exist"
		c.HTML(http.StatusNotFound, "account/transfer-money-form", h)
		return
	}

	if err := persistence.DB.First(&toAccount, "id = ?", h.Form.Data["to_account"]).Error; err != nil {
		h.Form.Errors["general"] = "Account does not exist"
		c.HTML(http.StatusNotFound, "account/transfer_money_form", h)
		return
	}

	amount, err := strconv.ParseFloat(h.Form.Data["amount"], 64)

	if err != nil {
		h.Form.Errors["amount"] = "Invalid currency value"
		c.HTML(http.StatusUnprocessableEntity, "account/transfer_money_form", h)
	}

	if err := h.transactionService.TransferMoney(fromAccount, toAccount, amount); err != nil {
		switch err := err.Error(); err {
		case "amount must be greater than 0":
			h.Form.Errors["amount"] = "Amount must be greater than 0"
			c.HTML(http.StatusUnprocessableEntity, "account/transfer_money_form", h)
			return
		case "not enough money":
			h.Form.Errors["from_account"] = "You do not have enough money in this account"
			c.HTML(http.StatusUnprocessableEntity, "account/transfer_money_form", h)
			return
		case "cannot transfer to same account":
			h.Form.Errors["to_account"] = "You cannot transfer money to the same account"
			c.HTML(http.StatusUnprocessableEntity, "account/transfer_money_form", h)
			return
		case "cannot transfer to other customers accounts":
			h.Form.Errors["general"] = "You do not have enough money in this account"
			c.HTML(http.StatusUnprocessableEntity, "account/transfer_money_form", h)
			return
		}
	}

	if err := h.customerService.FindByID(c.Param("id"), &h.Customer); err != nil {
		h.Form.Errors["general"] = "Could not find customer"
		c.HTML(http.StatusNotFound, "account/transfer-money-form", h)
		return
	}

	c.Header("HX-Trigger", "closeModal")
	c.HTML(http.StatusOK, "account/transfer_money_oob", h)
}

func (h customerHandler) isOwner(customerID, userID string) bool {
	var customer domain.Customer
	if err := h.customerService.FindByID(customerID, &customer); err != nil {
		return false
	}

	var user domain.User
	if err := h.userService.FindByID(userID, &user); err != nil {
		return false
	}

	return customer.Bank.UserID == user.ID || user.IsAdmin()
}
