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

type userResourceHandler struct {
	resource  entity.Resource
	secCookie *middleware.SecCookie
}

func getPureRequestURI(URI string, prefix string) string {
	firstIndex := strings.IndexRune(URI, '?')
	if firstIndex != -1 {
		return URI[len(prefix):firstIndex]
	}
	return URI[len(prefix):]
}

func (h userResourceHandler) ListResources(c *gin.Context) {
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

func (h userResourceHandler) CreateResource(c *gin.Context) {
	currentUser := middleware.Auth.GetCurrentUser(c)
	if currentUser == nil {
		uer.HandleUnauthorized(c)
		return
	}

	var resourceForm form.Resource
	err := resourceForm.FromCtx(c)
	if err != nil {
		uer.HandleErrorGin(err, c)
		return
	}

	err = h.resource.Create(resourceForm.Content, currentUser.Id)
	if err != nil {
		uer.HandleErrorGin(err, c)
		return
	}

	c.Status(200)
}

func (h userResourceHandler) DeleteResource(c *gin.Context) {
	c.Status(200)
}
