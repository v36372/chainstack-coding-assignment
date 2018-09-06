package handler

import (
	"chainstack/app/entity"
	"chainstack/app/form"
	"chainstack/app/params"
	"chainstack/app/view"
	"chainstack/middleware"
	"chainstack/models"
	"chainstack/utilities/uer"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	user      entity.User
	secCookie *middleware.SecCookie
}

func (h userHandler) Login(c *gin.Context) {
	var loginForm form.UserLogin
	err := loginForm.FromCtx(c)
	if err != nil {
		uer.HandleErrorGin(err, c)
		return
	}

	user, err := h.user.Login(loginForm.Email, loginForm.Password)
	if err != nil {
		uer.HandleErrorGin(err, c)
		return
	}

	_, err = h.secCookie.SetAuthorizationToken("auth", user.Email, "/", c.Writer)
	if err != nil {
		err = uer.InternalError(err)
		uer.HandleErrorGin(err, c)
		return
	}

	userViews, err := view.NewUsers([]models.User{*user})
	if err != nil {
		uer.HandleErrorGin(err, c)
		return
	}

	view.ResponseOK(c, userViews[0])
}

func (h userHandler) DeleteUser(c *gin.Context) {
	currentUser := middleware.Auth.GetCurrentUser(c)
	if currentUser == nil {
		uer.HandleUnauthorized(c)
		return
	}

	userId := params.GetUserIdUrlParam(c)
	if userId == 0 {
		err := uer.NotFoundError(errors.New("User not found"))
		uer.HandleErrorGin(err, c)
		return
	}

	if userId == currentUser.Id {
		err := uer.BadRequestError(errors.New("You can not delete yourself"))
		uer.HandleErrorGin(err, c)
		return
	}

	err := h.user.Delete(userId, currentUser.Id)
	if err != nil {
		uer.HandleErrorGin(err, c)
		return
	}

	c.JSON(200, gin.H{
		"msg": "User and its resources are deleted.",
	})
}

func (h userHandler) CreateUser(c *gin.Context) {
	currentUser := middleware.Auth.GetCurrentUser(c)
	if currentUser == nil {
		uer.HandleUnauthorized(c)
		return
	}

	var userForm form.UserCreateForm
	err := userForm.FromCtx(c)
	if err != nil {
		uer.HandleErrorGin(err, c)
		return
	}

	user := userForm.ToDBModel()
	if user == nil {
		err = uer.BadRequestError(errors.New("Invalid input"))
		uer.HandleErrorGin(err, c)
		return
	}

	user, err = h.user.Create(user, currentUser.Id)
	if err != nil {
		uer.HandleErrorGin(err, c)
		return
	}

	userViews, err := view.NewUsers([]models.User{*user})
	if err != nil {
		uer.HandleErrorGin(err, c)
		return
	}

	view.ResponseOK(c, userViews[0])
}

func (h userHandler) ListUsers(c *gin.Context) {
	currentUser := middleware.Auth.GetCurrentUser(c)
	if currentUser == nil {
		uer.HandleUnauthorized(c)
		return
	}

	nextId, limit := params.GetNextIdAndLimitParam(c)
	users, err := h.user.ListUsers(nextId, limit)
	if err != nil {
		uer.HandleErrorGin(err, c)
		return
	}

	userViews, err := view.NewUsers(users)
	if err != nil {
		uer.HandleErrorGin(err, c)
		return
	}

	var pagination *view.Pagination
	if len(users) == limit {
		url := c.Request.RequestURI
		if strings.Index(url, "/") == 0 {
			url = url[1:]
		}
		version := c.Param("version")
		// NextUrl
		nextUrl := fmt.Sprintf("%s?%s=%d",
			getPureRequestURI(url, version),
			params.ParamUrlNextId,
			users[len(users)-1].Id)
		pagination = &view.Pagination{
			NextUrl: nextUrl,
		}
	}
	view.ResponseOKWithPagination(c, userViews, pagination)

}
