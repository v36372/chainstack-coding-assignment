package params

import "github.com/gin-gonic/gin"

func GetUserIdUrlParam(c *gin.Context) int {
	return parseUrlParamToInt(c.Param(paramUrlUserId), 0)
}

func GetResourceUIDParam(c *gin.Context) int {
	return parseUrlParamToInt(c.Param(paramUrlUniqueId), 0)
}

func GetNextIdGetResourcesParam(c *gin.Context) (nextId, limit int) {
	nextId = parseUrlParamToInt(c.Query(ParamUrlNextId), 0)

	limit = parseUrlParamToInt(c.Query(paramUrlLimit), defaultLimit)
	if limit <= 0 || limit > maxLimit {
		limit = defaultLimit
	}

	return
}
