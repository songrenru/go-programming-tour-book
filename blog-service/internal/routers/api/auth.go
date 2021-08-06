package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/internal/service"
	"github.com/go-programming-tour-book/blog-service/pkg/app"
	"github.com/go-programming-tour-book/blog-service/pkg/errcode"
)

func GetAuth(c *gin.Context) {
	param := service.AuthRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		response.ToErrorResponse(errcode.InvalidParams.WriteDetails(errs.Errors()...))
		global.Logger.Errorf("app.BindAndValid errs: %v", errs)
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.CheckAuth(&param)
	if err != nil {
		response.ToErrorResponse(errcode.UnauthorizedAuthNotExist)
		global.Logger.Errorf("svc.CheckAuth  err: %v", err)
		return
	}

	token, err := app.GenerateToken(param.AppKey, param.AppSecret)
	if err != nil {
		response.ToErrorResponse(errcode.UnauthorizedTokenGenerate)
		global.Logger.Errorf("app.GenerateToken  err: %v", err)
		return
	}

	response.ToResponse(gin.H{
		"token": token,
	})
}
