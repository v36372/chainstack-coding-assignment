package repo

import (
	"chainstack/infra"
	"chainstack/models"
	"chainstack/utilities/uer"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type resource struct {
	base
}

var Resource IResource

func init() {
	Resource = resource{}
}

type IResource interface {
	GetByUserId(userId, nextId, limit int) ([]models.Resource, error)
	GetById(id int) (*models.Resource, error)
	Create(*models.Resource) (*models.Resource, error)
	Delete(*models.Resource) error
}

func (r resource) GetById(id int) (*models.Resource, error) {
	var resource models.Resource
	err := infra.PostgreSql.Model(models.Resource{}).
		Where("id = ?", id).
		Find(&resource).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &resource, err
}

func (resource) createWithTx(resource *models.Resource, tx *gorm.DB) (*models.Resource, error) {
	err := tx.Create(resource).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return resource, nil
}

func (r resource) Create(resource *models.Resource) (*models.Resource, error) {
	var userQuota models.UserQuota
	tx := infra.PostgreSql.Begin()

	err := tx.Model(models.UserQuota{}).Where("user_id = ?", resource.CreatedBy).Find(&userQuota).Error
	if err == gorm.ErrRecordNotFound {
		return r.createWithTx(resource, tx)
	}
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if userQuota.CurrentQuotaLeft <= 0 {
		err = uer.BadRequestError(errors.New("Cant create new resource, you ran out of quota."))
		return nil, err
	}

	err = tx.Create(resource).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	userQuota.CurrentQuotaLeft -= 1
	userQuota.UpdatedAt = time.Now()
	userQuota.UpdatedBy = resource.CreatedBy
	err = tx.Save(&userQuota).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return resource, nil

}

func (r resource) Delete(resource *models.Resource) error {
	return r.delete(resource)
}

func (r resource) GetByUserId(userId, nextId, limit int) (resources []models.Resource, err error) {
	query := infra.PostgreSql.Model(models.Resource{}).
		Where("created_by = ?", userId)

	if nextId > 0 {
		query = query.Where("id < ?", nextId)
	}

	err = query.Order("id desc").Limit(limit).Find(&resources).Error
	return
}
