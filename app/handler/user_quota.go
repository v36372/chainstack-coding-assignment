package handler

import (
	"chainstack/app/entity"
	"chainstack/app/form"
	"chainstack/app/params"
	"chainstack/app/view"
	"chainstack/middleware"
	"chainstack/utilities/uer"
	"errors"

	"github.com/gin-gonic/gin"
)

type userQuotaHandler struct {
	quota entity.Quota
}

func (h userQuotaHandler) UpdateQuota(c *gin.Context) {
	currentUser := middleware.Auth.GetCurrentUser(c)
	if currentUser == nil {
		uer.HandleUnauthorized(c)
		return
	}

	var updateQuotaForm form.UserQuota
	err := updateQuotaForm.FromCtx(c)
	if err != nil {
		uer.HandleErrorGin(err, c)
		return
	}

	userId := params.GetUserIdUrlParam(c)
	if userId == 0 {
		err := uer.NotFoundError(errors.New("User not found"))
		uer.HandleErrorGin(err, c)
		return
	}

	user, userQuota, err := h.quota.UpdateByUserId(userId, updateQuotaForm.Quota, currentUser.Id)
	if err != nil {
		uer.HandleErrorGin(err, c)
		return
	}

	userView := view.NewUserWithQuota(user, userQuota)
	view.ResponseOK(c, userView)
}
