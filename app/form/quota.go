package form

import (
	"chainstack/utilities/uer"

	"github.com/gin-gonic/gin"
)

type UserQuota struct {
	Quota int `form:"quota" json:"quota""`
}

func (userQuota *UserQuota) FromCtx(c *gin.Context) error {
	if err := c.Bind(userQuota); err != nil {
		return uer.BadRequestError(err)
	}

	return nil
}
