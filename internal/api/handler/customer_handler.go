package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type customerHandler struct {
	bankService     service.BankService
	customerService service.CustomerService
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
		ModalType:       "create_customer_modal",
		Form:            NewFormData(),
		Bank:            domain.Bank{},
		SignedIn:        true,
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

	if err := h.customerService.FindByID(id, &h.Customer); err != nil {
		c.HTML(http.StatusNotFound, "not-found", h)
		return
	}

	c.HTML(http.StatusOK, "customer", h)
}

func (h customerHandler) OpenSettings(c *gin.Context) {
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
	id := c.Param("id")
	h.Form = GetForm(c)

	h.customerService.FindByID(id, &h.Customer)

	h.Customer.FirstName = h.Form.Data["first_name"]
	h.Customer.LastName = h.Form.Data["last_name"]
	h.Customer.PIN = h.Form.Data["pin"]

	if err := h.customerService.Update(&h.Customer); err != nil {
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

func (h customerHandler) CreateCustomer(c *gin.Context) {
	h.Form = GetForm(c)

	bankID, err := strconv.Atoi(h.Form.Data["bank_id"])

	if err != nil {
		c.HTML(http.StatusBadRequest, "create_customer_form", h)
		return
	}

	customer := domain.Customer{
		BankID:    uint(bankID),
		FirstName: h.Form.Data["first_name"],
		LastName:  h.Form.Data["last_name"],
		PIN:       h.Form.Data["pin"],
	}

	if err := h.customerService.Create(&customer); err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			h.Form.Errors["general"] = "A customer with that PIN already exists in this bank"
			h.Form.Errors["pin"] = "A customer with that PIN already exists in this bank"
		} else {
			h.Form.Errors["general"] = "Something went wrong creating your customer"
		}

		c.HTML(http.StatusUnprocessableEntity, "create_customer_form", h)
		return
	}

	if err := h.bankService.FindByID(h.Form.Data["bank_id"], &h.Bank); err != nil {
		c.HTML(http.StatusUnprocessableEntity, "create_customer_form", h)
		return
	}

	c.Header("HX-Trigger", "closeModal")
	c.HTML(http.StatusCreated, "customers_oob", h)
}
