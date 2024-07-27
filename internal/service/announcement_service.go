package service

import (
	"strconv"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/pagination"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
)

type AnnouncementService interface {
	FindAll(pagingInfo *pagination.PagingInfo[domain.Announcement]) error
	FindByID(id string, announcement *domain.Announcement) error
	Recent(announcements *[]domain.Announcement) error
	Create(announcement *domain.Announcement) error
	Update(id string, announcement *domain.Announcement) error
	Delete(id string) error
}

type announcementService struct{}

func NewAnnoucementService() AnnouncementService {
	return announcementService{}
}

func (s announcementService) FindAll(pagingInfo *pagination.PagingInfo[domain.Announcement]) error {
	return persistence.DB.
		Model(&domain.Announcement{}).
		Count(&pagingInfo.TotalItems).
		Order("created_at desc").
		Offset((pagingInfo.PageNumber - 1) * pagingInfo.ItemsPerPage).
		Limit(pagingInfo.ItemsPerPage).
		Find(&pagingInfo.Items).Error
}

func (s announcementService) FindByID(id string, announcement *domain.Announcement) error {
	return persistence.DB.Find(&announcement, "id = ?", id).Error
}

func (s announcementService) Recent(announcements *[]domain.Announcement) error {
	return persistence.DB.Order("created_at desc").Limit(5).Find(&announcements).Error
}

func (s announcementService) Create(announcement *domain.Announcement) error {
	if err := persistence.DB.Create(&announcement).Error; err != nil {
		return err
	}
	return s.FindByID(strconv.Itoa(announcement.ID), announcement)
}

func (s announcementService) Update(id string, announcement *domain.Announcement) error {
	if err := persistence.DB.Where("id = ?", id).Updates(&announcement).Error; err != nil {
		return err
	}
	return s.FindByID(strconv.Itoa(announcement.ID), announcement)
}

func (s announcementService) Delete(id string) error {
	return persistence.DB.Delete(&domain.Announcement{}, "id = ?").Error
}
