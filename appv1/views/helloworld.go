package views

import (
	"gin-template/pkg/utils"

	"github.com/gin-gonic/gin"
)

func ActionHelloWorld(c *gin.Context) {
	utils.ReturnResutl(c, utils.RetCode.Success, "", map[string]interface{}{"text": "hello world"})
}
