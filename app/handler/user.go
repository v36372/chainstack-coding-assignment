package handler

import (
	"github.com/gin-gonic/gin"
)

type userHandler struct {
}

func (h userHandler) Login(c *gin.Context) {
	c.Status(200)
}
