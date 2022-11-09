package client

import (
	"context"
	pb "gin-template/service/proto"
)

type ServiceMailInterface interface {
	Call(toAddress []string, subject, text string) error
}

type ServiceMailApi struct {
	c pb.MailService
}

func (srv *ServiceMailApi) Call(toAddress []string, subject, text string) error {

	_, err := srv.c.Call(context.Background(), &pb.MailCallRequest{ToAddress: toAddress, Subject: subject, Text: text})

	return err
}
