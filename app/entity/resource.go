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
	Create(content string, createdBy int) error
}

func NewResource(resourceRepo repo.IResource) Resource {
	return &resourceEntity{
		resourceRepo: resourceRepo,
	}
}

func (r resourceEntity) Create(content string, createdBy int) error {
	resource := &models.Resource{
		Content:   content,
		CreatedBy: createdBy,
	}
	err := r.resourceRepo.Create(resource)
	if err != nil {
		err = uer.InternalError(err)
		return err
	}

	return err
}

func (r resourceEntity) GetByUserId(userId, nextId, limit int) (resources []models.Resource, err error) {
	resources, err = r.resourceRepo.GetByUserId(userId, nextId, limit)
	if err != nil {
		err = uer.InternalError(err)
		return
	}

	return
}
