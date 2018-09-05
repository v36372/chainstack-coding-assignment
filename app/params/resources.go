package params

import "github.com/gin-gonic/gin"

func GetUserIdUrlParam(c *gin.Context) int {
	return parseUrlParamToInt(c.Param(paramUrlUserId), 0)
}

func GetResourceUIDParam(c *gin.Context) int {
	return parseUrlParamToInt(c.Param(paramUrlUniqueId), 0)
}
