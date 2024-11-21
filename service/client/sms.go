package client

import (
	"context"
	pb "gin-template/service/proto"
)

type ServiceSmsInterface interface {
	Call(ctx context.Context, code, phoneNum string) error
}

type ServiceSmsApi struct {
	c pb.SmsService
}

func (srv *ServiceSmsApi) Call(ctx context.Context, code, phoneNum string) error {

	c := context.Background()
	c = context.WithValue(c, "trace_id", ctx.Value("trace_id"))

	_, err := srv.c.Call(c, &pb.SmsCallRequest{Code: code, PhoneNum: phoneNum})

	return err
}
