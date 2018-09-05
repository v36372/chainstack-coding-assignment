package handler

import (
	"chainstack/app/entity"
	"chainstack/app/form"
	"chainstack/app/params"
	"chainstack/app/view"
	"chainstack/middleware"
	"chainstack/utilities/uer"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	user      entity.User
	resource  entity.Resource
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

func getPureRequestURI(URI string, prefix string) string {
	firstIndex := strings.IndexRune(URI, '?')
	if firstIndex != -1 {
		return URI[len(prefix):firstIndex]
	}
	return URI[len(prefix):]
}

func (h userHandler) ListResources(c *gin.Context) {
	currentUser := middleware.Auth.GetCurrentUser(c)
	if currentUser == nil {
		uer.HandleUnauthorized(c)
		return
	}

	userId := params.GetUserIdUrlParam(c)
	if userId == 0 {
		userId = currentUser.Id
	}

	if userId != currentUser.Id && currentUser.IsAdmin == false {
		uer.HandlePermissionDenied(c)
		return
	}

	nextId, limit := params.GetNextIdGetResourcesParam(c)

	resources, err := h.resource.GetByUserId(userId, nextId, limit)
	if err != nil {
		uer.HandleErrorGin(err, c)
		return
	}

	var pagination *view.Pagination
	resourceViews := view.NewResources(resources)
	if len(resources) == limit {
		url := c.Request.RequestURI
		if strings.Index(url, "/") == 0 {
			url = url[1:]
		}
		version := c.Param("version")
		// NextUrl
		nextUrl := fmt.Sprintf("%s?%s=%d",
			getPureRequestURI(url, version),
			params.ParamUrlNextId,
			resources[len(resources)-1].Id)
		pagination = &view.Pagination{
			NextUrl: nextUrl,
		}
	}
	view.ResponseOKWithPagination(c, resourceViews, pagination)
}

func (h userHandler) CreateResource(c *gin.Context) {
	c.Status(200)
}

func (h userHandler) DeleteResource(c *gin.Context) {
	c.Status(200)
}
