package view

import "github.com/gin-gonic/gin"

type Response struct {
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	NextUrl string `json:"nextUrl"`
}

func ResponseOK(c *gin.Context, data interface{}) {
	c.JSON(200, Response{Data: data})
}

func ResponseOKWithPagination(c *gin.Context, data interface{}, pagination *Pagination) {
	c.JSON(200, Response{
		Data:       data,
		Pagination: pagination,
	})
}

func NewPagination(nextUrl string) Pagination {
	return Pagination{
		NextUrl: nextUrl,
	}
}
