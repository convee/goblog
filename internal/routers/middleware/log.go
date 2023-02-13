package middleware

import (
	"fmt"
	"github.com/convee/artgo"
	"github.com/convee/goblog/pkg/logger"
	"time"
)

func Logger() artgo.HandlerFunc {
	return func(c *artgo.Context) {
		t := time.Now()
		c.Next()
		logger.Info(fmt.Sprintf("[%d] %s %v", c.StatusCode, c.Req.RequestURI, time.Since(t)))
	}
}
