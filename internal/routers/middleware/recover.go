package middleware

import (
	"fmt"
	"github.com/convee/artgo"
	"github.com/convee/goblog/conf"
	"github.com/convee/goblog/pkg/ding"
	"github.com/convee/goblog/pkg/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
)

func RecoverWrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				logger.Error("http_router_panic", zap.Any("err", r), zap.Stack(string(debug.Stack())))
				if !conf.Conf.App.DisableDingDing {
					_, _ = ding.SendAlert(fmt.Sprintf("http_router_panic:err:%v;stack:%s", r, string(debug.Stack())), false)
				}
				http.Error(writer, err.Error(), http.StatusInternalServerError)

			}
		}()
		h.ServeHTTP(writer, request)
	})
}

func Recover() artgo.HandlerFunc {
	return func(c *artgo.Context) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				logger.Error("http_router_panic", zap.Any("err", r), zap.Stack(string(debug.Stack())))
				if !conf.Conf.App.DisableDingDing {
					_, _ = ding.SendAlert(fmt.Sprintf("http_router_panic:err:%v;stack:%s", r, string(debug.Stack())), false)
				}
				c.Error(http.StatusInternalServerError, err.Error())

			}
		}()
		c.Next()
	}
}
