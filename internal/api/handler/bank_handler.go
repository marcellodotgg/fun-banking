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
	bankService service.BankService
	Banks       []domain.Bank
	ModalType   string
	Form        FormData
	Bank        domain.Bank
	SignedIn    bool
}

func NewBankHandler() bankHandler {
	return bankHandler{
		bankService: service.NewBankService(),
		Banks:       []domain.Bank{},
		ModalType:   "create_bank_modal",
		Form:        NewFormData(),
		Bank:        domain.Bank{},
		SignedIn:    true,
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
		UserID:      uint(userID),
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

	c.Header("HX-Redirect", fmt.Sprintf("/banks/%s/%s", bank.User.Username, bank.Slug))
}

func (h bankHandler) UpdateBank(c *gin.Context) {
	h.Form = GetForm(c)

	bankId := h.Form.Data["bank_id"]
	h.Bank.Name = h.Form.Data["name"]
	h.Bank.Description = h.Form.Data["description"]

	if err := h.bankService.Update(bankId, &h.Bank); err != nil {
		h.Form.Errors["general"] = "A bank with that name already exists"
		c.HTML(http.StatusUnprocessableEntity, "update_bank_form", h)
		return
	}

	c.Header("HX-Redirect", fmt.Sprintf("/banks/%s/%s", h.Bank.User.Username, h.Bank.Slug))
}

func (h bankHandler) CreateModal(c *gin.Context) {
	h.Form = NewFormData()
	h.ModalType = "create_bank_modal"
	c.HTML(http.StatusOK, "modal", h)
}

func (h bankHandler) ViewBank(c *gin.Context) {
	h.bankService.FindByUsernameAndSlug(c.Param("username"), c.Param("slug"), &h.Bank)
	c.HTML(http.StatusOK, "bank", h)
}

func (h bankHandler) Settings(c *gin.Context) {
	h.Form = NewFormData()
	h.ModalType = "update_bank_modal"
	bankId := c.Query("id")

	if err := h.bankService.FindByID(bankId, &h.Bank); err != nil {
		c.HTML(http.StatusNotFound, "modal", h)
		return
	}

	h.Form.Data["name"] = h.Bank.Name
	h.Form.Data["description"] = h.Bank.Description
	h.Form.Data["bank_id"] = bankId

	c.HTML(http.StatusOK, "modal", h)
}
