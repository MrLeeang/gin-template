package client

import (
	"context"
	pb "gin-template/service/proto"
)

type ServiceSmsInterface interface {
	Call(code, phoneNum string) error
}

type ServiceSmsApi struct {
	c pb.SmsService
}

func (srv *ServiceSmsApi) Call(code, phoneNum string) error {

	_, err := srv.c.Call(context.Background(), &pb.SmsCallRequest{Code: code, PhoneNum: phoneNum})

	return err
}
