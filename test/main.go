package main

import (
	"context"
	"fmt"
	"gin-template/pkg/config"

	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"

	pb "gin-template/service/proto"
)

var (
	service = "gin.template.service"
	version = "latest"
)

func main() {
	srv := micro.NewService(
		micro.Name(service),
		micro.Version(version),
		micro.Registry(consul.NewRegistry(registry.Addrs(config.Global.Consul.Address))),
	)

	resp, err := pb.NewService(service, srv.Client()).Call(context.Background(), &pb.CallRequest{Name: "lihongwei"})

	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Msg)
}
