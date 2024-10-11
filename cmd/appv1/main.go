package main

import (
	"context"
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

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

var (
	server       *http.Server
	changeConfig bool
)

func main() {

	config.InitializeConfig()

	// 初始化zaplogger
	logger.InitializeLogger()

	defer logger.Logger.Sync() // 确保在程序结束时 flush 日志

	db.InitializeDatabase()

	go runServer()

	go watchConfig()

	// 等待系统信号
	waitForShutdown()
}

func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if changeConfig {
			return
		}

		if err := viper.Unmarshal(&config.Global); err != nil {
			panic(err)
		}

		changeConfig = true

		restartServer()

		time.AfterFunc(2*time.Second, func() {
			changeConfig = false
		})
	})
}

func runServer() {

	if !config.Global.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(
		middleware.GinLogger(),
		middleware.GinRecovery(true),
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

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func restartServer() {

	logger.Infof("Restarting Server...")

	zap.L().Sync()

	logger.InitializeLogger()

	db.InitializeDatabase()

	// 优雅地关闭当前服务
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Infof("Server Shutdown Failed:%+v", err)
	}

	// 重新启动服务
	go runServer()
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
