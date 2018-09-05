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
	GetById(userId int) (*models.User, error)
	List(nextId, limit int) ([]models.User, error)
	CreateUserAndQuota(*models.User, *models.UserQuota) error
	Create(*models.User) error
	Update(*models.User) error
	Delete(*models.User) error
	DeleteUserAndQuota(userId int) error
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

func (u user) DeleteUserAndQuota(userId int) error {
	var quota models.UserQuota
	tx := infra.PostgreSql.Begin()

	err := tx.Model(models.UserQuota{}).Where("user_id = ?", userId).Find(&quota).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return err
	}

	err = tx.Where("id = ?", userId).Delete(models.User{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if quota.Id == 0 {
		err = tx.Commit().Error
		if err != nil {
			tx.Rollback()
			return err
		}

		return nil
	}

	err = tx.Delete(&quota).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
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

func (user) GetById(userId int) (*models.User, error) {
	var user models.User

	err := infra.PostgreSql.Model(models.User{}).
		Where("id = ?", userId).
		Limit(1).
		Find(&user).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &user, err
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
