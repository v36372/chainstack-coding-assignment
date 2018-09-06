package models

import "time"

type UserQuota struct {
	Id               int
	UserId           int
	Quota            int
	CurrentQuotaLeft int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	UpdatedBy        int
	CreatedBy        int
}

func (UserQuota) TableName() string {
	return "user_quotas"
}
