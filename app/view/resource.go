package view

import (
	"chainstack/models"
	"time"
)

type Resource struct {
	Id        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedBy int       `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewResource(resource models.Resource) Resource {
	return Resource{
		Id:        resource.Id,
		Content:   resource.Content,
		CreatedBy: resource.CreatedBy,
		CreatedAt: resource.CreatedAt,
	}
}

func NewResources(resources []models.Resource) (resourceViews []Resource) {
	resourceViews = make([]Resource, len(resources))
	for i, resource := range resources {
		resourceViews[i] = NewResource(resource)
	}

	return
}
