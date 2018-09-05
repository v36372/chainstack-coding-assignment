package params

import "strconv"

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
