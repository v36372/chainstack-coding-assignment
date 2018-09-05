package form

import (
	"chainstack/utilities/uer"

	"github.com/gin-gonic/gin"
)

type UserLogin struct {
	Email    string `form:"email" json:"email" validator:"required"`
	Password string `form:"password" json:"password" validator:"required"`
}

func (inputForm *UserLogin) FromCtx(c *gin.Context) error {
	if err := c.Bind(inputForm); err != nil {
		return uer.BadRequestError(err)
	}

	return nil
}
