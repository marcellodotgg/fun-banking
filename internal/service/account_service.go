package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/pagination"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"github.com/bytebury/fun-banking/internal/utils"
)

type AccountService interface {
	FindByID(id string, account *domain.Account) error
	Create(account *domain.Account) error
	Update(id string, account *domain.Account) error
	UpdateBalance(id string, account *domain.Account) error
	CashFlow(accountID string, cashflow *Cashflow) error
	Transactions(accountID string, pagingInfo *pagination.PagingInfo[domain.Transaction]) error
	TransactionsByPeriod(accountID string, period string, pagingInfo *pagination.PagingInfo[domain.Transaction]) error
}

type accountService struct{}

func NewAccountService() AccountService {
	return accountService{}
}

func (s accountService) FindByID(id string, account *domain.Account) error {
	return persistence.DB.
		Preload("Customer.Bank").
		Preload("Customer.Accounts").
		First(&account, "id = ?", id).Error
}

func (s accountService) Create(account *domain.Account) error {
	if err := persistence.DB.Create(&account).Error; err != nil {
		return err
	}
	return s.FindByID(strconv.Itoa(account.ID), account)
}

func (s accountService) Update(id string, account *domain.Account) error {
	if err := persistence.DB.Where("id = ?", id).Updates(&account).Error; err != nil {
		return err
	}
	return s.FindByID(id, account)
}

func (s accountService) UpdateBalance(id string, account *domain.Account) error {
	if err := persistence.DB.Where("id = ?", id).Select("balance").Updates(&account).Error; err != nil {
		return err
	}
	return s.FindByID(id, account)
}

func (s accountService) CashFlow(accountID string, cashflow *Cashflow) error {
	month := time.Now().Month()

	if err := s.sumDepositsByAccount(accountID, month, &cashflow.Deposits); err != nil {
		cashflow.Deposits = 0
	}

	if err := s.sumWithdrawalsByAccount(accountID, month, &cashflow.Withdrawals); err != nil {
		cashflow.Withdrawals = 0
	}

	return nil
}

func (s accountService) Transactions(accountID string, pagingInfo *pagination.PagingInfo[domain.Transaction]) error {
	if err := persistence.DB.Find(&domain.Transaction{}, "account_id = ?", accountID).Count(&pagingInfo.TotalItems).Error; err != nil {
		return err
	}
	return persistence.DB.
		Offset((pagingInfo.PageNumber-1)*pagingInfo.ItemsPerPage).
		Limit(pagingInfo.ItemsPerPage).
		Order("updated_at DESC").
		Find(&pagingInfo.Items, "account_id = ?", accountID).Error
}

func (s accountService) TransactionsByPeriod(accountID string, period string, pagingInfo *pagination.PagingInfo[domain.Transaction]) error {
	if err := persistence.DB.Find(&domain.Transaction{}, "account_id = ?", accountID).Count(&pagingInfo.TotalItems).Error; err != nil {
		return err
	}
	return persistence.DB.
		Offset((pagingInfo.PageNumber-1)*pagingInfo.ItemsPerPage).
		Limit(pagingInfo.ItemsPerPage).
		Order("updated_at DESC").
		Where("strftime('%Y-%m', updated_at) = ? AND account_id = ?", period, accountID).
		Find(&pagingInfo.Items, "account_id = ?", accountID).Error
}

func (s accountService) sumWithdrawalsByAccount(accountID string, month time.Month, sum *float64) error {
	period := fmt.Sprintf("%d-%s", time.Now().Year(), utils.ConvertMonthToNumeric(month))

	return persistence.DB.
		Model(&domain.Transaction{}).
		Where("strftime('%Y-%m', updated_at) = ? AND amount <= ? AND account_id = ?", period, 0, accountID).
		Where("status = ?", domain.TransactionApproved).
		Select("sum(amount)").
		Row().
		Scan(sum)
}

func (s accountService) sumDepositsByAccount(accountID string, month time.Month, sum *float64) error {
	period := fmt.Sprintf("%d-%s", time.Now().Year(), utils.ConvertMonthToNumeric(month))

	return persistence.DB.
		Model(&domain.Transaction{}).
		Where("strftime('%Y-%m', updated_at) = ? AND amount >= ? AND account_id = ?", period, 0, accountID).
		Where("status = ?", domain.TransactionApproved).
		Select("sum(amount)").
		Row().
		Scan(sum)
}
