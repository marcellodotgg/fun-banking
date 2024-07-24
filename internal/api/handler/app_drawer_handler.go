package handler

import (
	"net/http"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type appDrawerHandler struct {
	SignedIn         bool
	CustomerSignedIn bool
	CurrentUser      domain.User
	CurrentCustomer  domain.Customer
	userService      service.UserService
	customerService  service.CustomerService
}

func NewAppDrawerHandler() appDrawerHandler {
	return appDrawerHandler{
		SignedIn:         false,
		CustomerSignedIn: false,
		CurrentUser:      domain.User{},
		CurrentCustomer:  domain.Customer{},
		userService:      service.NewUserService(),
		customerService:  service.NewCustomerService(),
	}
}

func (h appDrawerHandler) Open(c *gin.Context) {
	userID := c.GetString("user_id")
	customerID := c.GetString("customer_id")

	h.SignedIn = userID != ""
	h.CustomerSignedIn = customerID != ""

	if h.SignedIn {
		h.userService.FindByID(userID, &h.CurrentUser)
	}

	if h.CustomerSignedIn {
		h.customerService.FindByID(customerID, &h.CurrentCustomer)
	}

	c.HTML(http.StatusOK, "app_drawer", h)
}
