package entity

import (
	"chainstack/models"
	"chainstack/repo"
	"chainstack/utilities/uer"
)

type resourceEntity struct {
	resourceRepo repo.IResource
}

type Resource interface {
	GetByUserId(userId, nextId, limit int) ([]models.Resource, error)
}

func NewResource(resourceRepo repo.IResource) Resource {
	return &resourceEntity{
		resourceRepo: resourceRepo,
	}
}

func (r resourceEntity) GetByUserId(userId, nextId, limit int) (resources []models.Resource, err error) {
	resources, err = r.resourceRepo.GetByUserId(userId, nextId, limit)
	if err != nil {
		err = uer.InternalError(err)
		return
	}

	return
}
