package handler

import (
	"net/http"
	"strings"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type customerHandler struct {
	bankService     service.BankService
	customerService service.CustomerService
	userService     service.UserService
	ModalType       string
	Form            FormData
	Bank            domain.Bank
	Customer        domain.Customer
	SignedIn        bool
}

func NewCustomerHandler() customerHandler {
	return customerHandler{
		bankService:     service.NewBankService(),
		customerService: service.NewCustomerService(),
		userService:     service.NewUserService(),
		ModalType:       "create_customer_modal",
		Form:            NewFormData(),
		Bank:            domain.Bank{},
		SignedIn:        false,
	}
}

func (h customerHandler) OpenCreateModal(c *gin.Context) {
	h.ModalType = "create_customer_modal"
	h.Form = NewFormData()
	h.Form.Data["bank_id"] = c.Query("bank_id")
	c.HTML(http.StatusOK, "modal", h)
}

func (h customerHandler) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	h.SignedIn = c.GetString("user_id") != ""

	isCustomer := c.GetString("customer_id") == c.Param("id")
	if !isCustomer && !h.hasAccess(c.Param("id"), c.GetString("user_id")) {
		c.HTML(http.StatusForbidden, "forbidden", h)
		return
	}

	if err := h.customerService.FindByID(id, &h.Customer); err != nil {
		c.HTML(http.StatusNotFound, "not-found", h)
		return
	}

	c.HTML(http.StatusOK, "customer", h)
}

func (h customerHandler) OpenSettingsModal(c *gin.Context) {
	id := c.Param("id")
	h.ModalType = "customer_settings"
	h.Form = NewFormData()

	h.customerService.FindByID(id, &h.Customer)

	h.Form.Data["first_name"] = h.Customer.FirstName
	h.Form.Data["last_name"] = h.Customer.LastName
	h.Form.Data["pin"] = h.Customer.PIN

	c.HTML(http.StatusOK, "modal", h)
}

func (h customerHandler) Update(c *gin.Context) {
	h.Form = GetForm(c)

	if !h.hasAccess(c.Param("id"), c.GetString("user_id")) {
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

		h.Form.Errors["general"] = "Something went wrong updating your customer"
		c.HTML(http.StatusOK, "customer_settings_form", h)
		return
	}

	h.Form.Data["success"] = "Successfully updated the customer"
	c.HTML(http.StatusOK, "customer_settings_oob", h)
}

func (h customerHandler) hasAccess(customerID, userID string) bool {
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
