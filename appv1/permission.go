package appv1

import (
	"gin-template/db"
	"gin-template/pkg/utils"

	"github.com/gin-gonic/gin"
)

// PermissionHandler 登录验证
func PermissionHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		tokenStr := ctx.Request.Header.Get("X-Token")

		verify, err := utils.VerifyToken(tokenStr)

		if err != nil {
			utils.ReturnResutl(ctx, utils.RetCode.VerifyFailed, "", map[string]interface{}{})
			ctx.Abort()
			return
		}

		if verify.ClientIp != ctx.ClientIP() {
			utils.ReturnResutl(ctx, utils.RetCode.VerifyFailed, "错误的Token", map[string]interface{}{})
			ctx.Abort()
			return
		}

		user, err := db.QueryUserByUuid(ctx, verify.Uid)

		if err != nil { // 获取用户信息失败
			utils.ReturnResutl(ctx, utils.RetCode.VerifyFailed, "用户不存在", user)
			ctx.Abort()
			return
		}

		// gin上下文存储ctx.Keys，verifyMap保存到上下文中
		if ctx.Keys == nil {
			ctx.Keys = map[string]any{}
		}

		ctx.Keys["Uid"] = verify.Uid
		ctx.Keys["UserName"] = user.Username

		ctx.Next()
	}
}
