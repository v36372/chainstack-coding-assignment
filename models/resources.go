package models

import "time"

type Resource struct {
	Id        int
	Content   string
	CreatedBy int
	CreatedAt time.Time
	UpdatedAt time.Time
}
