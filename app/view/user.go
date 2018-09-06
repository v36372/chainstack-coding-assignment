package view

import (
	"chainstack/models"
	"chainstack/repo"
	"chainstack/utilities/uer"
	"time"
)

type User struct {
	Id               int       `json:"id"`
	Email            string    `json:"email"`
	IsAdmin          bool      `json:"isAdmin"`
	Quota            *int      `json:"quota,omitempty"`
	QuotaStatus      string    `json:"quotaStatus,omitempty"`
	CurrentQuotaLeft *int      `json:"currentQuotaLeft,omitempty"`
	CreatedBy        int       `json:"createdBy"`
	CreatedAt        time.Time `json:"createdAt"`
}

func NewUserWithQuota(user *models.User, userQuota *models.UserQuota) User {
	userView := User{
		Id:        user.Id,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
		CreatedBy: user.CreatedBy,
		CreatedAt: user.CreatedAt,
	}

	if userQuota == nil {
		userView.Quota = nil
		userView.CurrentQuotaLeft = nil
		userView.QuotaStatus = "Infinite"
		return userView
	}

	var quota, currentQuotaLeft int
	quota = userQuota.Quota
	currentQuotaLeft = userQuota.CurrentQuotaLeft
	userView.Quota = &quota
	userView.CurrentQuotaLeft = &currentQuotaLeft
	return userView
}

func NewUser(user models.User, quotaMap map[int]models.UserQuota) User {
	userView := User{
		Id:        user.Id,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
		CreatedBy: user.CreatedBy,
		CreatedAt: user.CreatedAt,
	}

	userQuota, exist := quotaMap[user.Id]
	if !exist {
		userView.Quota = nil
		userView.CurrentQuotaLeft = nil
		userView.QuotaStatus = "Infinite"
		return userView
	}

	var quota, currentQuotaLeft int
	quota = userQuota.Quota
	currentQuotaLeft = userQuota.CurrentQuotaLeft
	userView.Quota = &quota
	userView.CurrentQuotaLeft = &currentQuotaLeft
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

	quotaMap := map[int]models.UserQuota{}
	for _, userQuota := range userQuotas {
		quotaMap[userQuota.UserId] = userQuota
	}

	userViews = make([]User, len(users))
	for i, user := range users {
		userViews[i] = NewUser(user, quotaMap)
	}

	return
}
