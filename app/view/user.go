package view

import (
	"chainstack/models"
	"chainstack/repo"
	"chainstack/utilities/uer"
	"time"
)

type User struct {
	Id          int       `json:"id"`
	Email       string    `json:"email"`
	IsAdmin     bool      `json:"isAdmin"`
	Quota       int       `json:"quota"`
	QuotaStatus string    `json:"quotaStatus,omitempty"`
	CreatedBy   int       `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
}

func NewUser(user models.User, quotaMap map[int]int) User {
	userView := User{
		Id:        user.Id,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
		CreatedBy: user.CreatedBy,
		CreatedAt: user.CreatedAt,
	}

	userQuota, exist := quotaMap[user.Id]
	if !exist {
		userView.Quota = -1
		userView.QuotaStatus = "Infinite"
		return userView
	}

	userView.Quota = userQuota
	return userView
}

func NewUsers(users []models.User) (userViews []User, err error) {
	userIds := make([]int, len(users))
	for i, user := range users {
		userIds[i] = user.Id
	}

	userQuotas, err := repo.UserQuota.GetByUserIds(userIds)
	if err != nil {
		err = uer.InternalError(err)
		return
	}

	quotaMap := map[int]int{}
	for _, userQuota := range userQuotas {
		quotaMap[userQuota.UserId] = userQuota.Quota
	}

	userViews = make([]User, len(users))
	for i, user := range users {
		userViews[i] = NewUser(user, quotaMap)
	}

	return
}
