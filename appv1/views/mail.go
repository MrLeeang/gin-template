package views

import (
	"gin-template/pkg/utils"
	"gin-template/service/client"

	"github.com/gin-gonic/gin"
)

type sendMailParams struct {
	ToAddress []string `json:"to_address" binding:"required"`
	Subject   string   `json:"subject" binding:"required"`
	Text      string   `json:"text" binding:"required"`
}

func ActionSendMail(c *gin.Context) {

	var params sendMailParams

	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ReturnResutl(c, utils.RetCode.ParamRequired, err.Error(), params)
		return
	}

	srv := client.NewService()
	if err := srv.Mail().Call(c, params.ToAddress, params.Subject, params.Text); err != nil {
		utils.ReturnResutl(c, utils.RetCode.ExceptionError, err.Error(), params)
		return
	}

	utils.ReturnResutl(c, utils.RetCode.Success, "", params)
}
