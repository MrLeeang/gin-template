package client

import (
	"context"
	pb "gin-template/service/proto"
)

type ServiceMailInterface interface {
	Call(ctx context.Context, toAddress []string, subject, text string) error
}

type ServiceMailApi struct {
	c pb.MailService
}

func (srv *ServiceMailApi) Call(ctx context.Context, toAddress []string, subject, text string) error {

	_, err := srv.c.Call(ctx, &pb.MailCallRequest{ToAddress: toAddress, Subject: subject, Text: text})

	return err
}
