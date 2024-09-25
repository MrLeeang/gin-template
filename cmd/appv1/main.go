package main

import (
	"fmt"
	"gin-template/appv1"
	"gin-template/pkg/config"
	"gin-template/pkg/logger"
	"gin-template/pkg/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	defer logger.Logger.Sync() // 确保在程序结束时 flush 日志

	r.Use(
		middleware.GinLogger(logger.Logger),
		middleware.GinRecovery(logger.Logger, true),
		middleware.Cors(),
		middleware.RateLimit(rate.NewLimiter(rate.Every(time.Millisecond*time.Duration(1000/config.Global.Server.MaxRequest)), 50)),
	)

	r.StaticFS("/static", http.Dir(config.Global.Server.UploadDir))

	appv1.MakeRouter(r)

	logger.Infof("run server success on %s !!!", config.Global.Server.ServerPort)

	_ = r.Run(fmt.Sprintf(":%s", config.Global.Server.ServerPort))

}
