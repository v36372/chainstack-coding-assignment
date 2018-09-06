package entity

import (
	"chainstack/models"
	"chainstack/repo"
	"chainstack/utilities/uer"
	"errors"
)

type resourceEntity struct {
	resourceRepo repo.IResource
}

type Resource interface {
	GetByUserId(userId, nextId, limit int) ([]models.Resource, error)
	Create(content string, createdBy int) (*models.Resource, error)
	Delete(id, deleteBy int, isDeletorAdmin bool) error
}

func NewResource(resourceRepo repo.IResource) Resource {
	return &resourceEntity{
		resourceRepo: resourceRepo,
	}
}

func (r resourceEntity) Create(content string, createdBy int) (*models.Resource, error) {
	resource := &models.Resource{
		Content:   content,
		CreatedBy: createdBy,
	}
	resource, err := r.resourceRepo.Create(resource)
	if err != nil {
		err = uer.InternalError(err)
		return nil, err
	}

	return resource, err
}

func (r resourceEntity) Delete(id, deletedBy int, isDeletorAdmin bool) error {
	resource, err := r.resourceRepo.GetById(id)
	if err != nil {
		err = uer.InternalError(err)
		return err
	}

	if resource == nil {
		err = uer.NotFoundError(errors.New("Resource not found"))
		return err
	}

	if resource.CreatedBy != deletedBy && isDeletorAdmin == false {
		err = uer.ForbiddenError(errors.New("You can not touch this"))
		return err
	}

	err = r.resourceRepo.Delete(resource)
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
