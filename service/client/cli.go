package client

import (
	"gin-template/pkg/config"

	pb "gin-template/service/proto"

	"go-micro.dev/v4/registry"

	"github.com/go-micro/plugins/v4/registry/consul"

	"go-micro.dev/v4"
)

var (
	service = "gin.template.service"
	version = "latest"
)

type Service struct {
	Srv micro.Service
}

func NewService() *Service {

	// Create service
	srv := micro.NewService(
		micro.Name(service),
		micro.Version(version),
		micro.Registry(consul.NewRegistry(registry.Addrs(config.Config.Consul.Address))),
	)

	return &Service{
		Srv: srv,
	}
}

func (s *Service) Mail() ServiceMailInterface {
	return &ServiceMailApi{
		c: pb.NewMailService(service, s.Srv.Client()),
	}
}
func (s *Service) Sms() ServiceSmsInterface {
	return &ServiceSmsApi{
		c: pb.NewSmsService(service, s.Srv.Client()),
	}
}
