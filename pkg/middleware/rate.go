package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimit(limiter *rate.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		err := limiter.WaitN(ctx, 10)
		if err != nil {
			c.AbortWithStatus(http.StatusServiceUnavailable)
			return
		}
		// 处理请求
		c.Next()
	}
}
