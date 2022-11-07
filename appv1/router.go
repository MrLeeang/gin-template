package appv1

import (
	"gin-template/appv1/views"

	"github.com/gin-gonic/gin"
)

func MakeRouter(r *gin.Engine) {

	v1 := r.Group("/v1", PermissionHandler())

	v1.GET("/", views.ActionHelloWorld)

}
