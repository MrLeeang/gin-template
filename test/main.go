package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	logger, _ := zap.NewProduction(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.Fields(zap.String("appName", "test-gin")),
	)

	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		logger.Sugar().Infof("hello world")
		c.JSON(200, map[string]interface{}{"Hello": "World"})
	})

	r.Run(":8080")
}
