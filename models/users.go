package models

import "time"

type User struct {
	Id        int
	Email     string
	Password  string
	Salt      string
	IsAdmin   bool
	CreatedBy int
	CreatedAt time.Time
	UpdatedAt time.Time
}
