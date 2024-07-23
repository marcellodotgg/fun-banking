package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type bankHandler struct {
	bankService     service.BankService
	customerService service.CustomerService
	Banks           []domain.Bank
	ModalType       string
	Form            FormData
	Bank            domain.Bank
	SignedIn        bool
}

func NewBankHandler() bankHandler {
	return bankHandler{
		bankService:     service.NewBankService(),
		customerService: service.NewCustomerService(),
		Banks:           []domain.Bank{},
		ModalType:       "create_bank_modal",
		Form:            NewFormData(),
		Bank:            domain.Bank{},
		SignedIn:        true,
	}
}

func (h bankHandler) MyBanks(c *gin.Context) {
	h.bankService.MyBanks(c.GetString("user_id"), &h.Banks)
	c.HTML(http.StatusOK, "my_banks", h)
}

func (h bankHandler) CreateBank(c *gin.Context) {
	h.Form = GetForm(c)

	bank := domain.Bank{
		UserID:      c.GetString("user_id"),
		Name:        h.Form.Data["name"],
		Description: h.Form.Data["description"],
	}

	if err := h.bankService.Create(&bank); err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			h.Form.Errors["general"] = "You already have a bank by that name"
			h.Form.Data["name"] = ""
		} else {
			h.Form.Errors["general"] = "Something went wrong creating your bank"
		}
		c.HTML(http.StatusUnprocessableEntity, "create_bank_form", h)
		return
	}

	c.Header("HX-Redirect", fmt.Sprintf("/banks/%s", bank.ID))
}

func (h bankHandler) UpdateBank(c *gin.Context) {
	h.Form = GetForm(c)

	h.Bank.Name = h.Form.Data["name"]
	h.Bank.Description = h.Form.Data["description"]

	if err := h.bankService.Update(c.Param("id"), &h.Bank); err != nil {
		h.Form.Errors["general"] = "A bank with that name already exists"
		c.HTML(http.StatusUnprocessableEntity, "update_bank_form", h)
		return
	}

	c.Header("HX-Redirect", fmt.Sprintf("/banks/%s", c.Param("id")))
}

func (h bankHandler) OpenCreateModal(c *gin.Context) {
	h.Form = NewFormData()
	h.ModalType = "create_bank_modal"
	c.HTML(http.StatusOK, "modal", h)
}

func (h bankHandler) OpenCreateCustomerModal(c *gin.Context) {
	h.ModalType = "create_customer_modal"
	h.Form = NewFormData()
	h.bankService.FindByID(c.Param("id"), &h.Bank)
	c.HTML(http.StatusOK, "modal", h)
}

func (h bankHandler) CreateCustomer(c *gin.Context) {
	h.Form = GetForm(c)

	customer := domain.Customer{
		BankID:    c.Param("id"),
		FirstName: h.Form.Data["first_name"],
		LastName:  h.Form.Data["last_name"],
		PIN:       h.Form.Data["pin"],
	}

	if err := h.bankService.FindByID(c.Param("id"), &h.Bank); err != nil {
		c.HTML(http.StatusUnprocessableEntity, "create_customer_form", h)
		return
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

	c.Header("HX-Trigger", "closeModal")
	c.HTML(http.StatusCreated, "customers_oob", h)
}

func (h bankHandler) ViewBank(c *gin.Context) {
	h.bankService.FindByID(c.Param("id"), &h.Bank)
	c.HTML(http.StatusOK, "bank", h)
}

func (h bankHandler) OpenSettingsModal(c *gin.Context) {
	h.Form = NewFormData()
	h.ModalType = "update_bank_modal"

	if err := h.bankService.FindByID(c.Param("id"), &h.Bank); err != nil {
		c.HTML(http.StatusNotFound, "modal", h)
		return
	}

	h.Form.Data["name"] = h.Bank.Name
	h.Form.Data["description"] = h.Bank.Description

	c.HTML(http.StatusOK, "modal", h)
}

func (h bankHandler) CustomerSearch(c *gin.Context) {
	var customers []domain.Customer
	h.customerService.FindAllByBankIDAndName(c.Param("id"), c.Query("name"), 5, &customers)

	c.HTML(http.StatusOK, "search_bank_customers", customers)
}

func (h bankHandler) FilterCustomers(c *gin.Context) {
	var customers []domain.Customer
	h.customerService.FindAllByBankIDAndName(c.Param("id"), c.Query("search"), 15, &customers)

	c.HTML(http.StatusOK, "customers_table", customers)
}
