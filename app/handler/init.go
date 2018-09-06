package handler

import (
	"chainstack/app/entity"
	"chainstack/config"
	"chainstack/middleware"
	"chainstack/repo"

	"github.com/gin-gonic/gin"
)

func InitEngine(conf *config.Config) *gin.Engine {
	if conf.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())

	if conf.App.Debug {
		r.Use(gin.Logger())
	}

	// ----------------------   INIT AUTHENTICATE MIDDLEWARE
	secCookie := middleware.NewSetCookie(conf.CookieToken.BlockKey, conf.CookieToken.HashKey)
	authMiddleware := middleware.NewAuthMiddleware(secCookie, middleware.Auth.GetLoggedInUser)
	middleware.InitAuth(authMiddleware.GetCurrentUser)

	// ----------------------   INIT HANDLER
	userHandler := userHandler{
		user:      entity.NewUser(repo.User),
		secCookie: secCookie,
	}

	userResourceHandler := userResourceHandler{
		resource: entity.NewResource(repo.Resource),
	}

	userQuotaHandler := userQuotaHandler{
		quota: entity.NewQuota(repo.UserQuota, repo.User),
	}

	// ----------------------   INIT ROUTE

	indexGroup := r.Group("/api/:version")
	{
		POST(indexGroup, "/login", userHandler.Login)
	}

	userGroup := indexGroup.Group("/users")
	userGroup.Use(authMiddleware.RequireLogin())
	userGroup.Use(authMiddleware.Interception())
	{
		GET(userGroup, "/:id/resources", userResourceHandler.ListResources)
	}

	adminGroup := userGroup
	adminGroup.Use(authMiddleware.RequireAdmin())
	{
		GET(adminGroup, "", userHandler.ListUsers)
		POST(adminGroup, "", userHandler.CreateUser)
		DELETE(adminGroup, "/:id", userHandler.DeleteUser)
		PUT(adminGroup, "/:id/quota", userQuotaHandler.UpdateQuota)
	}

	resourceGroup := indexGroup.Group("/resources")
	resourceGroup.Use(authMiddleware.RequireLogin())
	resourceGroup.Use(authMiddleware.Interception())
	{
		DELETE(resourceGroup, "/:uid", userResourceHandler.DeleteResource)
		POST(resourceGroup, "", userResourceHandler.CreateResource)
	}

	return r
}

func GET(group *gin.RouterGroup, relativePath string, f func(*gin.Context)) {
	group.GET(relativePath, f)
}

func POST(group *gin.RouterGroup, relativePath string, f func(*gin.Context)) {
	group.POST(relativePath, f)
}

func PUT(group *gin.RouterGroup, relativePath string, f func(*gin.Context)) {
	group.PUT(relativePath, f)
}

func DELETE(group *gin.RouterGroup, relativePath string, f func(*gin.Context)) {
	group.DELETE(relativePath, f)
}
