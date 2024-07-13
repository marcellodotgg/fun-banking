package handler

import (
	"net/http"
	"strconv"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type customerHandler struct {
	customerService service.CustomerService
	ModalType       string
	Form            FormData
	BankID          string
}

func NewCustomerHandler() customerHandler {
	return customerHandler{
		customerService: service.NewCustomerService(),
		ModalType:       "create_customer_modal",
		Form:            NewFormData(),
		BankID:          "",
	}
}

func (h customerHandler) OpenCreateModal(c *gin.Context) {
	h.ModalType = "create_customer_modal"
	h.BankID = c.Query("bank_id")
	h.Form = NewFormData()
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
		// TODO
	}

	customer := domain.Customer{
		BankID:    uint(bankID),
		FirstName: h.Form.Data["first_name"],
		LastName:  h.Form.Data["last_name"],
		PIN:       h.Form.Data["pin"],
	}

	if err := h.customerService.Create(&customer); err != nil {
		// handle errors
	}

	// need to do that better.
	c.Header("HX-Redirect", "/")
}
