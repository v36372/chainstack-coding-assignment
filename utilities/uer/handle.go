package uer

import (
	"github.com/gin-gonic/gin"
)

func HandleErrorGin(e error, c *gin.Context) {
	err := ToStatusError(e)
	c.JSON(err.Status, gin.H{
		"msg": err.Message,
	})
	c.Abort()
}

func HandleNotFound(c *gin.Context) {
	c.AbortWithStatus(404)
}

func HandleUnauthorized(c *gin.Context) {
	c.AbortWithStatus(401)
}

func HandlePermissionDenied(c *gin.Context) {
	c.AbortWithStatus(403)
}
