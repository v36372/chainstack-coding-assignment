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

	r.Use(authMiddleware.Interception())

	// ----------------------   INIT HANDLER
	userHandler := userHandler{
		user:      entity.NewUser(repo.User),
		secCookie: secCookie,
	}

	// ----------------------   INIT ROUTE

	indexGroup := r.Group("/")
	{
		POST(indexGroup, "/login", userHandler.Login)
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