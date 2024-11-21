package views

import (
	"gin-template/pkg/utils"
	"gin-template/service/client"

	"github.com/gin-gonic/gin"
)

type sendSmsParams struct {
	Code     string `json:"code" binding:"required"`
	PhoneNum string `json:"phone_num" binding:"required"`
}

func ActionSendSms(c *gin.Context) {

	var params sendSmsParams

	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ReturnResutl(c, utils.RetCode.ParamRequired, err.Error(), params)
		return
	}

	srv := client.NewService()
	if err := srv.Sms().Call(c, params.Code, params.PhoneNum); err != nil {
		utils.ReturnResutl(c, utils.RetCode.ExceptionError, err.Error(), params)
		return
	}

	utils.ReturnResutl(c, utils.RetCode.Success, "", params)
}
