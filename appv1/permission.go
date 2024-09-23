package appv1

import (
	"gin-template/db"
	"gin-template/pkg/utils"

	"github.com/gin-gonic/gin"
)

// PermissionHandler 登录验证
func PermissionHandler() gin.HandlerFunc {
	return func(context *gin.Context) {

		tokenStr := context.Request.Header.Get("X-Token")

		verify, err := utils.VerifyToken(tokenStr)
		if err != nil {
			utils.ReturnResutl(context, utils.RetCode.VerifyFailed, "", map[string]interface{}{})
			context.Abort()
			return
		}

		if verify.ClientIp != context.ClientIP() {
			utils.ReturnResutl(context, utils.RetCode.VerifyFailed, "错误的Token", map[string]interface{}{})
			context.Abort()
			return
		}

		user, err := db.QueryUserByUuid(verify.Uid)

		if err != nil { // 获取用户信息失败
			utils.ReturnResutl(context, utils.RetCode.VerifyFailed, "用户不存在", user)
			context.Abort()
			return
		}

		// gin上下文存储context.Keys，verifyMap保存到上下文中
		verifyMap := map[string]interface{}{}
		verifyMap["Uid"] = verify.Uid
		verifyMap["UserName"] = user.Username

		context.Keys = verifyMap

		context.Next()
	}
}
