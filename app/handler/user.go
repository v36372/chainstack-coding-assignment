package handler

import (
	"chainstack/app/entity"
	"chainstack/app/form"
	"chainstack/app/params"
	"chainstack/app/view"
	"chainstack/middleware"
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

	c.Status(200)
}

func (h userHandler) CreateUser(c *gin.Context) {
	currentUser := middleware.Auth.GetCurrentUser(c)
	if currentUser == nil {
		uer.HandleUnauthorized(c)
		return
	}

	if currentUser.IsAdmin == false {
		uer.HandlePermissionDenied(c)
		return
	}

	var userForm form.UserCreateForm
	err := userForm.FromCtx(c)
	if err != nil {
		uer.HandleErrorGin(err, c)
		return
	}

	user, userQuota := userForm.ToDBModel()
	if user == nil {
		err = uer.BadRequestError(errors.New("Invalid input"))
		uer.HandleErrorGin(err, c)
		return
	}

	err = h.user.Create(user, userQuota, currentUser.Id)
	if err != nil {
		uer.HandleErrorGin(err, c)
		return
	}

	c.Status(200)
}

func (h userHandler) ListUsers(c *gin.Context) {
	currentUser := middleware.Auth.GetCurrentUser(c)
	if currentUser == nil {
		uer.HandleUnauthorized(c)
		return
	}

	if currentUser.IsAdmin == false {
		uer.HandlePermissionDenied(c)
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
