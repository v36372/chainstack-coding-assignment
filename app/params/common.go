package params

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	paramUrlLimit    = "limit"
	ParamUrlNextId   = "next"
	paramUrlUserId   = "id"
	paramUrlUniqueId = "uid"

	maxLimit     = 100
	defaultLimit = 20
)

func parseUrlParamToInt(value string, defaultVal int) int {
	res, err := strconv.Atoi(value)
	if err != nil {
		return defaultVal
	}

	return res
}

func GetNextIdAndLimitParam(c *gin.Context) (nextId, limit int) {
	nextId = parseUrlParamToInt(c.Query(ParamUrlNextId), 0)

	limit = parseUrlParamToInt(c.Query(paramUrlLimit), defaultLimit)
	if limit <= 0 || limit > maxLimit {
		limit = defaultLimit
	}

	return
}
