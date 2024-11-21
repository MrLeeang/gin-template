package main

import (
	"context"
	"flag"
	"fmt"
	"gin-template/appv1"
	"gin-template/db"
	"gin-template/pkg/config"
	"gin-template/pkg/logger"
	"gin-template/pkg/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	server *http.Server
	debug  bool
)

func main() {

	flag.BoolVar(&debug, "debug", false, "Open debug mode (default false)")

	flag.Parse()

	config.InitializeConfig()

	config.Global.Debug = debug

	log := logger.InitializeZapLogger()

	defer log.Sync()

	db.InitializeDatabase()

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(
		middleware.Trace(),
		middleware.GinLogger(log),
		middleware.GinRecovery(log, true),
		middleware.Cors(),
		middleware.RateLimit(rate.NewLimiter(rate.Every(time.Millisecond*time.Duration(1000/config.Global.Server.MaxRequest)), 50)),
	)

	r.StaticFS("/static", http.Dir(config.Global.Server.UploadDir))

	// 健康监测
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, map[string]interface{}{"status": "ok"})
	})

	appv1.MakeRouter(r)

	// 启动 HTTP 服务
	server = &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Global.Server.ServerPort),
		Handler: r,
	}

	logger.Infof("run server success on %s !!!", config.Global.Server.ServerPort)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// 等待系统信号
	waitForShutdown()
}

func waitForShutdown() {
	// 捕获系统信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c // 等待信号

	logger.Infof("Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Infof("Server Shutdown Failed:%+v", err)
	}
	logger.Infof("Server exited")
}
