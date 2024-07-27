package handler

import (
	"net/http"
	"strconv"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/pagination"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type announcementHandler struct {
	SignedIn                bool
	CurrentUser             domain.User
	Form                    FormData
	Announcement            domain.Announcement
	AnnouncementsPagination pagination.PagingInfo[domain.Announcement]
	userService             service.UserService
	announcementService     service.AnnouncementService
}

func NewAnnouncementHandler() announcementHandler {
	return announcementHandler{
		SignedIn:                true,
		CurrentUser:             domain.User{},
		Form:                    NewFormData(),
		Announcement:            domain.Announcement{},
		AnnouncementsPagination: pagination.PagingInfo[domain.Announcement]{},
		userService:             service.NewUserService(),
		announcementService:     service.NewAnnoucementService(),
	}
}

func (h announcementHandler) Dashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "announcements_dashboard", h)
}

func (h announcementHandler) RecentAnnouncements(c *gin.Context) {
	var announcements []domain.Announcement
	h.announcementService.Recent(&announcements)
	c.HTML(http.StatusOK, "recent_announcements", announcements)
}

func (h announcementHandler) FindAll(c *gin.Context) {
	pageNumber, err := strconv.Atoi(c.Query("page"))

	if err != nil {
		pageNumber = 1
	}

	h.AnnouncementsPagination.PageNumber = pageNumber
	h.AnnouncementsPagination.ItemsPerPage = 8

	if err := h.announcementService.FindAll(&h.AnnouncementsPagination); err != nil {
		c.HTML(http.StatusNotFound, "not-found", h)
		return
	}

	c.HTML(http.StatusOK, "announcements", h)
}

func (h announcementHandler) Create(c *gin.Context) {
	userID, _ := strconv.Atoi(c.GetString("user_id"))
	h.Form = GetForm(c)

	announcement := domain.Announcement{
		Title:       h.Form.Data["title"],
		Description: h.Form.Data["description"],
		UserID:      userID,
	}

	if err := h.announcementService.Create(&announcement); err != nil {
		h.Form.Errors["general"] = "Something went wrong posting your announcement."
		c.HTML(http.StatusUnprocessableEntity, "announcement_form", h.Form)
		return
	}

	h.Form = NewFormData()
	h.Form.Data["success"] = "Successfully posted the announcement"
	c.HTML(http.StatusOK, "announcement_form", h.Form)
}

func (h announcementHandler) Update(c *gin.Context) {
	announcementID := c.Param("id")
	userID, _ := strconv.Atoi(c.GetString("user_id"))
	h.Form = GetForm(c)

	announcement := domain.Announcement{
		Title:       h.Form.Data["title"],
		Description: h.Form.Data["description"],
		UserID:      userID,
	}

	if err := h.announcementService.Update(announcementID, &announcement); err != nil {
		h.Form.Errors["general"] = "Something went wrong posting your announcement."
		c.HTML(http.StatusUnprocessableEntity, "edit_announcement_form", h.Form)
		return
	}

	h.Form.Data["success"] = "Successfully updated the announcement"
	c.HTML(http.StatusOK, "edit_announcement_form", h.Form)
}

func (h announcementHandler) Edit(c *gin.Context) {
	announcementID := c.Param("id")
	h.announcementService.FindByID(announcementID, &h.Announcement)

	h.Form = NewFormData()
	h.Form.Data["id"] = strconv.Itoa(h.Announcement.ID)
	h.Form.Data["title"] = h.Announcement.Title
	h.Form.Data["description"] = h.Announcement.Description

	c.HTML(http.StatusOK, "edit_announcement", h)
}

func (h announcementHandler) FindByID(c *gin.Context) {
	userID := c.GetString("user_id")
	h.userService.FindByID(userID, &h.CurrentUser)
	h.announcementService.FindByID(c.Param("id"), &h.Announcement)
	c.HTML(http.StatusOK, "announcement", h)
}

func (h announcementHandler) Destroy(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", h)
}
