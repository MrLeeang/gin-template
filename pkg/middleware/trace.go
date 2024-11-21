package middleware

import (
	"gin-template/pkg/utils"

	"github.com/gin-gonic/gin"
)

func Trace() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		traceId := ctx.Request.Header.Get("trace_id")

		if traceId == "" {
			traceId = utils.GetUuid()
		}

		ctx.Keys = map[string]any{}

		ctx.Keys["trace_id"] = traceId

		// 设置响应头
		ctx.Header("trace_id", traceId)

		// 处理请求
		ctx.Next()
	}
}
