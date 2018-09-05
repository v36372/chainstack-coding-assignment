package repo

import (
	"chainstack/infra"
	"chainstack/models"

	"github.com/jinzhu/gorm"
)

type user struct {
	base
}

var User IUser

func init() {
	User = user{}
}

type IUser interface {
	GetByEmail(email string) (*models.User, error)
	List(nextId, limit int) ([]models.User, error)
	CreateUserAndQuota(*models.User, *models.UserQuota) error
	Create(*models.User) error
	Update(*models.User) error
	Delete(*models.User) error
}

func (u user) Create(user *models.User) error {
	return u.create(user)
}

func (u user) Update(user *models.User) error {
	return u.save(user)
}

func (u user) Delete(user *models.User) error {
	return u.delete(user)
}

func (u user) CreateUserAndQuota(user *models.User, quota *models.UserQuota) (err error) {
	if quota == nil {
		return u.create(user)
	}
	tx := infra.PostgreSql.Begin()
	err = tx.Create(user).Error
	if err != nil {
		tx.Rollback()
		return
	}

	quota.UserId = user.Id
	err = tx.Create(quota).Error
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return
	}

	return nil
}

func (user) List(nextId, limit int) (users []models.User, err error) {
	query := infra.PostgreSql.Model(models.User{})

	if nextId > 0 {
		query = query.Where("id < ?", nextId)
	}

	err = query.Order("id desc").Limit(limit).Find(&users).Error
	return
}

func (user) GetByEmail(email string) (*models.User, error) {
	var user models.User

	err := infra.PostgreSql.Model(models.User{}).
		Where("email = ?", email).
		Limit(1).
		Find(&user).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &user, err
}
