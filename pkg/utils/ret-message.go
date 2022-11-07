package utils

import (
	"encoding/json"
	"gin-template/pkg/config"

	"github.com/gin-gonic/gin"
)

// RetCodeModel RetCodeModel
type RetCodeModel struct {
	Success        int
	ExceptionError int
	ParamRequired  int
	NotFoundInfo   int
	ParamError     int
	LoginError     int
	VerifyFailed   int
}

var (

	// RetCode RetCode
	RetCode = RetCodeModel{
		Success:        0,
		ExceptionError: 500,
		ParamRequired:  1001,
		NotFoundInfo:   1002,
		ParamError:     1003,
		LoginError:     1004,
		VerifyFailed:   1005,
	}
	// ErrorCodeMessage ErrorCodeMessage
	ErrorCodeMessage = map[int]string{
		0:    "SUCCESS",
		500:  "异常错误",
		1001: "缺少参数",
		1002: "资源未找到",
		1003: "参数错误",
		1004: "登录失败",
		1005: "权限验证失败",
	}
)

func ReturnResutl(c *gin.Context, retCode int, msg string, result interface{}) {
	if msg == "" {
		msg = ErrorCodeMessage[retCode]
	}

	if config.Config.Server.Encrypt {
		// 加密返回值
		byteStr, _ := json.Marshal(result)

		baseStr, _ := AesEncrypt(byteStr)

		result = baseStr
	}

	c.JSON(200, gin.H{
		"data":       result,
		"error_code": retCode,
		"msg":        msg,
	})
}
