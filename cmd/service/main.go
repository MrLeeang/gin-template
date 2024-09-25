package main

import (
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
)

func main() {
	// Create service
	srv := micro.NewService(
		micro.Name(service),
		micro.Version(version),
		micro.Address(config.Global.Service.Address),
		micro.Registry(consul.NewRegistry(registry.Addrs(config.Global.Consul.Address))),
	)
	srv.Init()

	// Register handler
	pb.RegisterServiceHandler(srv.Server(), new(handler.Service))
	pb.RegisterMailServiceHandler(srv.Server(), new(handler.MailService))
	pb.RegisterSmsServiceHandler(srv.Server(), new(handler.SmsService))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Panicf(err.Error())
	}
}
