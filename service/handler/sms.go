package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-template/pkg/config"
	"gin-template/pkg/logger"
	pb "gin-template/service/proto"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
)

type SmsService struct{}

func (e *SmsService) Call(ctx context.Context, req *pb.SmsCallRequest, rsp *pb.SmsCallResponse) error {
	logger.Infof("Received Service.Call request: %v", req)

	code := req.Code
	phoneNum := req.PhoneNum

	accessKeyId := config.Global.Alibaba.AccessKeyId
	accessKeySecret := config.Global.Alibaba.AccessKeySecret

	client, _err := CreateClient(&accessKeyId, &accessKeySecret)
	if _err != nil {
		return _err
	}

	body := map[string]string{
		"code": code,
	}

	bodyByte, err := json.Marshal(body)

	if err != nil {
		return err
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String(config.Global.Alibaba.SignName),
		TemplateCode:  tea.String(config.Global.Alibaba.TemplateCode),
		PhoneNumbers:  tea.String(phoneNum),
		TemplateParam: tea.String(string(bodyByte)),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		res, _err := client.SendSmsWithOptions(sendSmsRequest, runtime)
		if _err != nil {
			return _err
		}

		if *res.Body.Code != "OK" {
			return fmt.Errorf(*res.Body.Message)
		}
		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 如有需要，请打印 error
		util.AssertAsString(error.Message)
	}
	return _err

}

func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}
