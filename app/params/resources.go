package params

import "github.com/gin-gonic/gin"

func GetResourceUIDParam(c *gin.Context) int {
	return parseUrlParamToInt(c.Param(paramUrlUniqueId), 0)
}
