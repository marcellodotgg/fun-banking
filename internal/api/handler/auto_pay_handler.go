package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"github.com/gin-gonic/gin"
)

func (h accountHandler) OpenAutoPayModal(c *gin.Context) {
	h.Reset(c)

	yesterday := time.Now().Format("2006-01-02")

	h.ModalType = "create_auto_pay"
	h.Form.Data["min_date"] = yesterday

	if err := h.accountService.FindByID(c.Param("id"), &h.Account); err != nil {
		c.HTML(http.StatusNotFound, "not-found", nil)
		return
	}

	c.HTML(http.StatusOK, "modal", h)
}

func (h accountHandler) AutoPay(c *gin.Context) {
	h.Reset(c)

	if err := h.accountService.FindByID(c.Param("id"), &h.Account); err != nil {
		c.HTML(http.StatusNotFound, "not-found", nil)
		return
	}

	persistence.DB.Find(&h.AutoPays, "account_id = ?", c.Param("id"))

	c.HTML(http.StatusOK, "auto_pay.html", h)
}

func (h accountHandler) CreateAutoPay(c *gin.Context) {
	h.Reset(c)
	h.ModalType = "create_auto_pay"

	startDate, startDateErr := time.Parse("2006-01-02", h.Form.Data["start_date"])
	amount, amountErr := strconv.ParseFloat(h.Form.Data["amount"], 64)
	accountId, _ := strconv.Atoi(c.Param("id"))

	if startDateErr != nil {
		h.Form.Errors["start_date"] = "You provided an invalid start date"
		c.HTML(http.StatusUnprocessableEntity, "create_auto_pay_form.html", h)
		return
	}

	if amountErr != nil {
		h.Form.Errors["amount"] = "Invalid value for amount"
		c.HTML(http.StatusUnprocessableEntity, "create_auto_pay_form.html", h)
		return
	}

	if h.Form.Data["type"] == "withdraw" {
		amount *= -1
	}

	autoPay := domain.AutoPay{
		Cadence:     h.Form.Data["cadence"],
		StartDate:   startDate,
		Amount:      amount,
		Description: h.Form.Data["description"],
		AccountID:   accountId,
		Active:      true,
	}

	// TODO: Move to service
	if err := persistence.DB.Create(&autoPay).Error; err != nil {
		h.Form.Errors["general"] = "Something went wrong creating your Auto Pay Rule"
		c.HTML(http.StatusUnprocessableEntity, "create_auto_pay_form.html", h)
		return
	}

	persistence.DB.Find(&h.AutoPays, "account_id = ?", c.Param("id"))
	h.accountService.FindByID(c.Param("id"), &h.Account)

	c.Header("HX-Trigger", "closeModal")
	c.HTML(http.StatusOK, "auto_pay_oob.html", h)
}

func (h accountHandler) UpdateAutoPay(c *gin.Context) {
	h.Reset(c)

	var autoPay domain.AutoPay
	if err := persistence.DB.First(&autoPay, "id = ?", c.Param("auto_pay_id")).Error; err != nil {
		c.HTML(http.StatusNotFound, "not-found", h)
		return
	}

	persistence.DB.Model(&autoPay).Select("Active").Updates(domain.AutoPay{Active: h.Form.Data["checked"] == "on"})
	persistence.DB.Find(&h.AutoPays, "account_id = ?", c.Param("id"))
	h.accountService.FindByID(c.Param("id"), &h.Account)

	c.HTML(http.StatusOK, "auto_pay_oob.html", h)
}
