package form

import (
	"chainstack/utilities/uer"

	"github.com/gin-gonic/gin"
)

type Resource struct {
	Content string `form:"content" json:"content""`
}

func (resourceForm *Resource) FromCtx(c *gin.Context) error {
	if err := c.Bind(resourceForm); err != nil {
		return uer.BadRequestError(err)
	}

	return nil
}
