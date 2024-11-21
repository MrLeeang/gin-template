package client

import (
	"gin-template/pkg/config"

	pb "gin-template/service/proto"

	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/registry"
	"golang.org/x/net/context"

	"github.com/go-micro/plugins/v4/registry/consul"

	"go-micro.dev/v4"
)

var (
	service = "gin.template.service"
	version = "latest"
)

type Service struct {
	Srv micro.Service
	ctx context.Context
}

func NewService(ctx context.Context) *Service {

	// Create service
	srv := micro.NewService(
		micro.Name(service),
		micro.Version(version),
		micro.Registry(consul.NewRegistry(registry.Addrs(config.Global.Consul.Address))),
	)

	if traceId := ctx.Value("trace_id"); traceId != nil {
		ctx = metadata.NewContext(ctx, map[string]string{
			"trace_id": traceId.(string),
		})
	}

	return &Service{
		Srv: srv,
		ctx: ctx,
	}
}

func (s *Service) Mail() ServiceMailInterface {
	return &ServiceMailApi{
		c:   pb.NewMailService(service, s.Srv.Client()),
		ctx: s.ctx,
	}
}
func (s *Service) Sms() ServiceSmsInterface {
	return &ServiceSmsApi{
		c:   pb.NewSmsService(service, s.Srv.Client()),
		ctx: s.ctx,
	}
}
