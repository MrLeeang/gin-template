package appv1

import (
	"gin-template/appv1/views"

	"github.com/gin-gonic/gin"
)

func MakeRouter(r *gin.Engine) {

	v1 := r.Group("/v1", PermissionHandler())

	// 登录
	r.POST("/v1/login", views.ActionLogin)
	// 退出
	v1.GET("/logout", views.ActionLogout)
	// 用户信息
	v1.GET("/user/info", views.ActionUserInfo)

	// 角色
	v1.GET("/role", views.ActionRoleList)
	v1.GET("/role/:uuid", views.ActionRoleQuery)
	v1.PUT("/role", views.ActionRolePut)
	v1.POST("/role", views.ActionRolePost)
	v1.DELETE("/role/:uuid", views.ActionRoleDelete)

	// 用户
	v1.GET("/user", views.ActionUserList)
	v1.GET("/user/:uuid", views.ActionUserQuery)
	v1.PUT("/user", views.ActionUserPut)
	v1.POST("/user", views.ActionUserPost)
	v1.DELETE("/user/:uuid", views.ActionUserDelete)
	// 操作日志
	v1.GET("/user/log", views.ActionUserLoginLog)

	// mail
	v1.POST("/sendmail", views.ActionSendMail)
	v1.POST("/sendsms", views.ActionSendSms)
}
