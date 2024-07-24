package handler

import (
	"fmt"
	"net/http"
	"strconv"
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

	userID, err := strconv.Atoi(c.GetString("user_id"))
	if err != nil {
		h.Form.Errors["general"] = "Bad user. Are you signed in?"
		c.HTML(http.StatusUnprocessableEntity, "create_bank_form", h)
		return
	}

	bank := domain.Bank{
		UserID:      userID,
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

	c.Header("HX-Redirect", fmt.Sprintf("/banks/%d", bank.ID))
}

func (h bankHandler) UpdateBank(c *gin.Context) {
	h.Form = GetForm(c)

	bankID := c.Param("id")
	h.Bank.Name = h.Form.Data["name"]
	h.Bank.Description = h.Form.Data["description"]

	if err := h.bankService.Update(bankID, &h.Bank); err != nil {
		h.Form.Errors["general"] = "A bank with that name already exists"
		c.HTML(http.StatusUnprocessableEntity, "update_bank_form", h)
		return
	}

	c.Header("HX-Redirect", fmt.Sprintf("/banks/%s", bankID))
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
	bankID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.HTML(http.StatusBadRequest, "create_customer_form", h)
		return
	}

	customer := domain.Customer{
		BankID:    bankID,
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

	if err := h.bankService.FindByID(c.Param("id"), &h.Bank); err != nil {
		h.Form.Errors["general"] = "Something went wrong fetching bank information, please refresh."
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
	bankID := c.Param("id")

	if err := h.bankService.FindByID(bankID, &h.Bank); err != nil {
		c.HTML(http.StatusNotFound, "modal", h)
		return
	}

	h.Form.Data["name"] = h.Bank.Name
	h.Form.Data["description"] = h.Bank.Description

	c.HTML(http.StatusOK, "modal", h)
}

func (h bankHandler) CustomerSearch(c *gin.Context) {
	bankID := c.Param("id")
	searchStr := c.Query("name")

	var customers []domain.Customer
	h.customerService.FindAllByBankIDAndName(bankID, searchStr, 5, &customers)

	c.HTML(http.StatusOK, "search_bank_customers", customers)
}

func (h bankHandler) FilterCustomers(c *gin.Context) {
	bankID := c.Param("id")
	searchStr := c.Query("search")

	var customers []domain.Customer
	h.customerService.FindAllByBankIDAndName(bankID, searchStr, 15, &customers)

	c.HTML(http.StatusOK, "customers_table", customers)
}
