package uer

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleErrorGin(e error, c *gin.Context) {
	err := ToStatusError(e)
	switch err.Status {
	case http.StatusNotFound, http.StatusUnauthorized, http.StatusForbidden:
	case http.StatusInternalServerError:
		c.AbortWithError(err.Status, e)
	case http.StatusBadRequest:
		c.AbortWithError(err.Status, err.Err)
		return
	}
	c.AbortWithStatus(err.Status)
}
