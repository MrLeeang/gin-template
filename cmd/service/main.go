package main

import (
	"flag"
	"gin-template/pkg/config"
	"gin-template/pkg/logger"
	"gin-template/service/handler"
	pb "gin-template/service/proto"

	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

var (
	service = "gin.template.service"
	version = "latest"
	debug   bool
)

func main() {

	flag.BoolVar(&debug, "debug", false, "Open debug mode (default false)")

	flag.Parse()

	config.InitializeConfig()

	config.Global.Debug = debug

	// 初始化zaplogger
	logger.InitializeLogger()

	defer logger.Logger.Sync() // 确保在程序结束时 flush 日志

	// Create service
	srv := micro.NewService(
		micro.Name(service),
		micro.Version(version),
		micro.Registry(consul.NewRegistry(registry.Addrs(config.Global.Consul.Address))),
	)

	if config.Global.Service.Address != "" {
		srv.Init(micro.Address(config.Global.Service.Address))
	}

	srv.Init()

	// Register handler
	pb.RegisterServiceHandler(srv.Server(), new(handler.Service))
	pb.RegisterMailServiceHandler(srv.Server(), new(handler.MailService))
	pb.RegisterSmsServiceHandler(srv.Server(), new(handler.SmsService))

	logger.Infof("run server success on %s !!!", config.Global.Service.Address)

	// Run service
	if err := srv.Run(); err != nil {
		panic(err.Error())
	}
}
