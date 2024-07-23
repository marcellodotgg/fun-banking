package handler

import (
	"net/http"
	"strconv"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/gin-gonic/gin"
)

type announcementHandler struct {
	SignedIn            bool
	Form                FormData
	Announcement        domain.Announcement
	announcementService service.AnnouncementService
}

func NewAnnouncementHandler() announcementHandler {
	return announcementHandler{
		SignedIn:            true,
		Form:                NewFormData(),
		Announcement:        domain.Announcement{},
		announcementService: service.NewAnnoucementService(),
	}
}

func (h announcementHandler) Dashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "announcements_dashboard", h)
}

func (h announcementHandler) Announcements(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", h)
}

func (h announcementHandler) RecentAnnouncements(c *gin.Context) {
	var announcements []domain.Announcement
	h.announcementService.Recent(&announcements)
	c.HTML(http.StatusOK, "recent_announcements", announcements)
}

func (h announcementHandler) Create(c *gin.Context) {
	userID, _ := strconv.Atoi(c.GetString("user_id"))
	h.Form = GetForm(c)

	announcement := domain.Announcement{
		Title:       h.Form.Data["title"],
		Description: h.Form.Data["description"],
		UserID:      uint(userID),
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
	c.HTML(http.StatusOK, "index.html", h)
}

func (h announcementHandler) FindByID(c *gin.Context) {
	h.announcementService.FindByID(c.Param("id"), &h.Announcement)
	c.HTML(http.StatusOK, "announcement", h)
}

func (h announcementHandler) Destroy(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", h)
}
