package middleware

import (
	"chainstack/models"
	"chainstack/repo"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ExpireCode       = "TOKEN_EXPIRED"
	BadCode          = "BAD_TOKEN"
	userGinKey       = "CurrentUser"
	IsLoggedInGinKey = "IsLoggedInGinKey"
)

var (
	ErrorPermissionDenied = errors.New("PermissionDenied")
)

type GetLoggedInUserFunc func(userIDStr string) (interface{}, error)
type GetCurrentUser func(c *gin.Context) (user interface{}, exists bool)

type authMiddlewareInterface interface {
	RequireLogin() gin.HandlerFunc
	Interception() gin.HandlerFunc
	GetCurrentUser(c *gin.Context) (user interface{}, exists bool)
}
type authMiddleware struct {
	secCookie       *SecCookie
	getLoggedInUser GetLoggedInUserFunc
}

func NewAuthMiddleware(secCookie *SecCookie, getLoggedInUser GetLoggedInUserFunc) authMiddlewareInterface {
	return &authMiddleware{
		secCookie:       secCookie,
		getLoggedInUser: getLoggedInUser,
	}
}

func (a *authMiddleware) Interception() gin.HandlerFunc {
	return func(c *gin.Context) {
		isLoggedIn := true
		userIdStr, err := a.secCookie.GetCurrentUserID(c.Request)
		if err != nil {
			if err != http.ErrNoCookie {
				if err.Error() == "securecookie: expired timestamp" {
					a.secCookie.ClearCookie(c.Writer, "auth", "/")
				} else if err.Error() == "securecookie: the value is not valid" {
					a.secCookie.ClearCookie(c.Writer, "auth", "/")
				}
				a.secCookie.ClearCookie(c.Writer, "auth", "/")
			}
			isLoggedIn = false
		}
		if isLoggedIn {
			user, err := a.getLoggedInUser(userIdStr)
			if err != nil {
				if err == ErrorPermissionDenied {
					a.secCookie.ClearCookie(c.Writer, "auth", "/")
					c.AbortWithStatus(403)
					return
				}
				a.secCookie.ClearCookie(c.Writer, "auth", "/")
				c.AbortWithStatus(401)
				return
			}

			c.Set(userGinKey, user)
		}
		c.Next()

	}
}

func (a *authMiddleware) RequireLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := a.secCookie.GetCurrentUserID(c.Request)

		if err != nil {
			c.JSON(401, gin.H{
				"msg": "You need to login to see this.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (a *authMiddleware) GetCurrentUser(c *gin.Context) (user interface{}, exists bool) {
	currentUser, exists := c.Get(userGinKey)
	return currentUser, exists
}

type auth struct {
	getCurrentUserFunc GetCurrentUser
}

var Auth auth

func InitAuth(f GetCurrentUser) {
	Auth.getCurrentUserFunc = f
}

func (a *auth) GetCurrentUser(c *gin.Context) *models.User {
	user, exists := a.getCurrentUserFunc(c)
	if exists {
		resp := user.(*models.User)
		return resp
	}
	return nil
}

func (auth) GetLoggedInUser(email string) (interface{}, error) {
	user, err := repo.User.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrorPermissionDenied
	}

	return user, nil
}
