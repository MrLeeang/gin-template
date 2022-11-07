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

	r.Use(
		middleware.GinLogger(),
		gin.RecoveryWithWriter(logger.Logger.Out),
		middleware.Cors(),
		middleware.RateLimit(rate.NewLimiter(rate.Every(time.Millisecond*time.Duration(1000/config.Config.MaxRequest)), 50)),
	)

	// r.Static("/static", config.Config.UploadDir)
	r.StaticFS("/static", http.Dir(config.Config.UploadDir))

	appv1.MakeRouter(r)

	_ = r.Run(fmt.Sprintf(":%s", config.Config.Server.ServerPort))
}
