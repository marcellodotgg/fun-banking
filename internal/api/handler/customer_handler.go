package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type customerHandler struct {
	pageObject
	bankService     service.BankService
	customerService service.CustomerService
	accountService  service.AccountService
	userService     service.UserService
	Bank            domain.Bank
	Customer        domain.Customer
}

func NewCustomerHandler() customerHandler {
	return customerHandler{
		bankService:     service.NewBankService(),
		customerService: service.NewCustomerService(),
		userService:     service.NewUserService(),
		accountService:  service.NewAccountService(),
		Bank:            domain.Bank{},
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
