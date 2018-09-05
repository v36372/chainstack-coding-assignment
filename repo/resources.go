package repo

import (
	"chainstack/infra"
	"chainstack/models"
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
	Create(*models.Resource) error
	Delete(*models.Resource) error
}

func (r resource) Create(resource *models.Resource) error {
	return r.create(resource)
}

func (r resource) Delete(resource *models.Resource) error {
	return r.delete(resource)
}

func (r resource) GetByUserId(userId, nextId, limit int) (resources []models.Resource, err error) {
	query := infra.PostgreSql.Model(models.Resource{}).
		Where("created_by= ?", userId)

	if nextId > 0 {
		query = query.Where("id < ?", nextId)
	}

	err = query.Order("id desc").Limit(limit).Find(&resources).Error
	return
}
