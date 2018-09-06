package entity

import (
	"chainstack/models"
	"chainstack/repo"
	"chainstack/utilities/uer"
	"errors"
	"time"
)

type quotaEntity struct {
	quotaRepo repo.IUserQuota
	userRepo  repo.IUser
}

type Quota interface {
	UpdateByUserId(userId, quota int, updatedBy int) (*models.User, *models.UserQuota, error)
}

func NewQuota(quotaRepo repo.IUserQuota) Quota {
	return &quotaEntity{
		quotaRepo: quotaRepo,
	}
}

func (q quotaEntity) UpdateByUserId(userId, quota int, updatedBy int) (user *models.User, userQuota *models.UserQuota, err error) {
	user, err = q.userRepo.GetById(userId)
	if err != nil {
		err = uer.InternalError(err)
		return
	}
	if user == nil {
		err = uer.NotFoundError(errors.New("User not found"))
		return
	}

	userQuota, err = q.quotaRepo.GetByUserId(userId)
	if err != nil {
		err = uer.InternalError(err)
		return
	}

	if userQuota != nil {
		userQuota.Quota = quota
		userQuota.UpdatedAt = time.Now()
		userQuota.UpdatedBy = updatedBy
		err = q.quotaRepo.Update(userQuota)
		if err != nil {
			err = uer.InternalError(err)
			return
		}

		return
	}

	if quota > 0 {
		userQuota = &models.UserQuota{
			UserId:    userId,
			Quota:     quota,
			CreatedBy: updatedBy,
		}

		userQuota, err = q.quotaRepo.Create(userQuota)
		if err != nil {
			err = uer.InternalError(err)
			return
		}
	}

	return
}
