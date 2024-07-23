package service

import (
	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
)

type AnnouncementService interface {
	FindByID(id string, announcement *domain.Announcement) error
	Recent(announcements *[]domain.Announcement) error
	Create(announcement *domain.Announcement) error
	Update(announcement *domain.Announcement) error
	Delete(id string) error
}

type announcementService struct{}

func NewAnnoucementService() AnnouncementService {
	return announcementService{}
}

func (s announcementService) FindByID(id string, announcement *domain.Announcement) error {
	return persistence.DB.Find(&announcement, "id = ?", id).Error
}

func (s announcementService) Recent(announcements *[]domain.Announcement) error {
	return persistence.DB.Order("created_at desc").Limit(5).Find(&announcements).Error
}

func (s announcementService) Create(announcement *domain.Announcement) error {
	return persistence.DB.Create(&announcement).Error
}

func (s announcementService) Update(announcement *domain.Announcement) error {
	return persistence.DB.Updates(&announcement).Error
}

func (s announcementService) Delete(id string) error {
	return persistence.DB.Delete(&domain.Announcement{}, "id = ?").Error
}
