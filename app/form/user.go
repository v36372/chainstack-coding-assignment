package form

import (
	"chainstack/models"
	"chainstack/utilities/uer"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type UserCreateForm struct {
	Email    string `form:"email" json:"email" valid:"email,required"`
	Password string `form:"password" json:"password" valid:"required"`
	IsAdmin  bool   `form:"isAdmin" json:"isAdmin"`
}

type UserLogin struct {
	Email    string `form:"email" json:"email" valid:"email,required"`
	Password string `form:"password" json:"password" valid:"required"`
}

func (loginForm *UserLogin) FromCtx(c *gin.Context) error {
	if err := c.Bind(loginForm); err != nil {
		return uer.BadRequestError(err)
	}

	_, err := govalidator.ValidateStruct(loginForm)
	if err != nil {
		return uer.BadRequestError(err)
	}

	return nil
}

func (userForm *UserCreateForm) FromCtx(c *gin.Context) error {
	if err := c.Bind(userForm); err != nil {
		return uer.BadRequestError(err)
	}

	_, err := govalidator.ValidateStruct(userForm)
	if err != nil {
		return uer.BadRequestError(err)
	}

	return nil
}

func (userForm *UserCreateForm) ToDBModel() (user *models.User) {
	user = &models.User{
		Email:    userForm.Email,
		Password: userForm.Password,
		IsAdmin:  userForm.IsAdmin,
	}

	return
}
