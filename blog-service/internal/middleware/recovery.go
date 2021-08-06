package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/pkg/app"
	"github.com/go-programming-tour-book/blog-service/pkg/email"
	"github.com/go-programming-tour-book/blog-service/pkg/errcode"
	"time"
)

func Recovery() gin.HandlerFunc {
	defaultMailer := email.NewEamil(&email.SMTPInfo{
		Host: global.EmailSetting.Host,
		Port: global.EmailSetting.Port,
		IsSSL: global.EmailSetting.IsSSL,
		UserName: global.EmailSetting.UserName,
		Password: global.EmailSetting.Password,
		From: global.EmailSetting.From,
	})
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				app.NewResponse(c).ToErrorResponse(errcode.ServerError)
				global.Logger.WithCallersFrames().Errorf("panic recover err: %v", err)
				err := defaultMailer.SendMail(
					global.EmailSetting.To,
					fmt.Sprintf("异常时间： %d", time.Now().Unix()),
					fmt.Sprintf("异常内容： %v", err),
				)
				if err != nil {
					global.Logger.Panicf("mail.SendMail err: %v", err)
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}
