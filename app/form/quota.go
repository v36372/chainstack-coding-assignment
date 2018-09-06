package form

import (
	"chainstack/utilities/uer"
	"errors"

	"github.com/gin-gonic/gin"
)

type UserQuota struct {
	Quota int `form:"quota" json:"quota"`
}

func (userQuota *UserQuota) FromCtx(c *gin.Context) error {
	if err := c.Bind(userQuota); err != nil {
		return uer.BadRequestError(err)
	}

	if userQuota.Quota < 0 {
		err := uer.BadRequestError(errors.New("Quota should not be negative"))
		return err
	}
	return nil
}
