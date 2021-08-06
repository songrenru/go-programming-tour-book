package app

import (
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/pkg/errcode"
	"net/http"
)

type Response struct {
	Ctx *gin.Context
}

type Pager struct {
	Page int `json:"page"`
	PageSize int `json:"page_size"`
	TotalRows int `json:"total_rows"`
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{Ctx: ctx}
}

func (this *Response) ToResponse(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	this.Ctx.JSON(http.StatusOK, data)
}

func (this *Response) ToResponseList(list interface{}, totalRows int) {
	this.Ctx.JSON(http.StatusOK, gin.H{
		"list": list,
		"pager": Pager{
			Page: GetPage(this.Ctx),
			PageSize: GetPageSize(this.Ctx),
			TotalRows: totalRows,
		},
	})
}

func (this *Response) ToErrorResponse(err *errcode.Error) {
	resp := gin.H{"code": err.Code(), "msg": err.Msg()}
	details := err.Detail()
	if len(details) > 0 {
		resp["details"] = details
	}
	this.Ctx.JSON(err.StatusCode(), resp)
}
