package main

import (
	"context"
	"flag"
	"gin-template/pkg/config"
	"gin-template/pkg/logger"
	"gin-template/service/handler"
	pb "gin-template/service/proto"

	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
)

var (
	service = "gin.template.service"
	version = "latest"
	debug   bool
)

// 定义日志中间件
func logMiddleware(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {

		md, ok := metadata.FromContext(ctx)
		if ok {
			traceId, ok := md.Get("trace_id")
			if ok {
				ctx = context.WithValue(ctx, "trace_id", traceId)
			}
		}
		return fn(ctx, req, rsp)
	}
}

func main() {

	flag.BoolVar(&debug, "debug", false, "Open debug mode (default false)")

	flag.Parse()

	config.InitializeConfig()

	config.Global.Debug = debug

	log := logger.InitializeZapLogger()

	defer log.Sync()

	// Create service
	srv := micro.NewService(
		micro.Name(service),
		micro.Version(version),
		micro.Registry(consul.NewRegistry(registry.Addrs(config.Global.Consul.Address))),
		micro.WrapHandler(logMiddleware), // 注册中间件,trace_id
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
