package repo

import (
	"chainstack/infra"
	"chainstack/models"

	"github.com/jinzhu/gorm"
)

type userQuota struct {
	base
}

var UserQuota IUserQuota

func init() {
	UserQuota = userQuota{}
}

type IUserQuota interface {
	GetByUserId(userId int) (*models.UserQuota, error)
	GetByUserIds(userIds []int) ([]models.UserQuota, error)
	Create(*models.UserQuota) (*models.UserQuota, error)
	Delete(*models.UserQuota) error
	Update(*models.UserQuota) error
}

func (u userQuota) Create(quota *models.UserQuota) (*models.UserQuota, error) {
	value, err := u.create(quota)
	return value.(*models.UserQuota), err
}

func (u userQuota) Delete(quota *models.UserQuota) error {
	return u.delete(quota)
}

func (u userQuota) Update(quota *models.UserQuota) error {
	return u.save(quota)
}

func (u userQuota) GetByUserId(userId int) (*models.UserQuota, error) {
	var userQuota models.UserQuota

	err := infra.PostgreSql.Model(models.UserQuota{}).
		Where("user_id = ?", userId).
		Limit(1).
		Find(&userQuota).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &userQuota, err
}

func (u userQuota) GetByUserIds(userIds []int) (userQuotas []models.UserQuota, err error) {

	err = infra.PostgreSql.Model(models.UserQuota{}).
		Where("user_id IN (?)", userIds).
		Find(&userQuotas).
		Error

	return
}
