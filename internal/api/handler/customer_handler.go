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
	customerService service.CustomerService
	ModalType       string
	Form            FormData
}

func NewCustomerHandler() customerHandler {
	return customerHandler{
		customerService: service.NewCustomerService(),
		ModalType:       "create_customer_modal",
		Form:            NewFormData(),
	}
}

func (h customerHandler) OpenCreateModal(c *gin.Context) {
	h.ModalType = "create_customer_modal"
	h.Form = NewFormData()
	h.Form.Data["bank_id"] = c.Query("bank_id")
	c.HTML(http.StatusOK, "modal", h)
}

func (h customerHandler) CreateCustomer(c *gin.Context) {
	c.Request.ParseForm()

	h.Form = NewFormData()
	for key, values := range c.Request.PostForm {
		if len(values) > 0 {
			h.Form.Data[key] = values[0]
		}
	}

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

	// need to do that better.
	c.HTML(http.StatusOK, "bank", h)
}
