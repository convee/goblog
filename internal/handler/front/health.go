package front

import (
	"context"
	"net/http"
	"time"

	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/daos"
)

func Health(c *artgo.Context) {
	c.Writer.WriteHeader(http.StatusNoContent)
}

func Ready(c *artgo.Context) {
	ctx, cancel := context.WithTimeout(c.Req.Context(), time.Second)
	defer cancel()
	if err := daos.GetDB().PingContext(ctx); err != nil {
		c.Error(http.StatusServiceUnavailable, "database unavailable")
		return
	}
	c.Writer.WriteHeader(http.StatusNoContent)
}
