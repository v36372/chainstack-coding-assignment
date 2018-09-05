package handler

import (
	"chainstack/app/entity"
	"chainstack/app/form"
	"chainstack/utilities/uer"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	user entity.User
}

func (h userHandler) Login(c *gin.Context) {
	var loginForm form.UserLogin
	err := loginForm.FromCtx(c)
	if err != nil {
		uer.HandleErrorGin(err, c)
		return
	}

	_, err = h.user.Login(loginForm.Email, loginForm.Password)
	if err != nil {
		uer.HandleErrorGin(err, c)
		return
	}

	c.Status(200)
}
