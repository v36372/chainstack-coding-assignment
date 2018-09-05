package models

import "time"

type UserQuota struct {
	Id        int
	UserId    int
	Quota     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserQuota) TableName() string {
	return "user_quotas"
}
