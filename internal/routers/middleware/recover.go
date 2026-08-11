package middleware

import (
	"fmt"
	"github.com/convee/artgo"
	"github.com/convee/goblog/conf"
	"github.com/convee/goblog/pkg/ding"
	"github.com/convee/goblog/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
)

func Recover() artgo.HandlerFunc {
	return func(c *artgo.Context) {
		defer func() {
			r := recover()
			if r != nil {
				logger.Error("http_router_panic", zap.Any("err", r), zap.ByteString("stack", debug.Stack()))
				if !conf.Conf.App.DisableDingDing {
					_, _ = ding.SendAlert(fmt.Sprintf("http_router_panic: %v", r), false)
				}
				c.Error(http.StatusInternalServerError, "服务器内部错误")

			}
		}()
		c.Next()
	}
}
