package handler

import (
	"context"
	"gin-template/pkg/config"
	"gin-template/pkg/logger"
	pb "gin-template/service/proto"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type MailService struct{}

func (e *MailService) Call(ctx context.Context, req *pb.MailCallRequest, rsp *pb.MailCallResponse) error {
	logger.Infof("Received Service.Call request: %v", req)

	em := email.NewEmail()
	// 设置 sender 发送方 的邮箱 ， 此处可以填写自己的邮箱
	em.From = config.Global.Mail.From

	// 设置 receiver 接收方 的邮箱  此处也可以填写自己的邮箱， 就是自己发邮件给自己
	em.To = req.ToAddress

	// 设置主题
	em.Subject = req.Subject

	// 简单设置文件发送的内容，暂时设置成纯文本
	em.HTML = []byte(req.Text)

	//设置服务器相关的配置
	return em.Send(config.Global.Mail.Address, smtp.PlainAuth("", config.Global.Mail.Username, config.Global.Mail.Password, config.Global.Mail.Host))

}
