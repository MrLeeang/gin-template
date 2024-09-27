package views

import (
	"gin-template/db"
	"gin-template/models"
	"gin-template/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ActionRoleList(c *gin.Context) {
	// 页码
	page := c.Query("page")
	// 每页显示数量
	size := c.Query("size")
	// 查询字段
	keyword := c.Query("keyword")

	var pageInt int
	var sizeInt int

	if page != "" && size != "" {
		pageNum, err := strconv.Atoi(page)
		if err != nil {
			utils.ReturnResutl(c, utils.RetCode.ExceptionError, err.Error(), page)
			return
		}
		sizeNum, err := strconv.Atoi(size)
		if err != nil {
			utils.ReturnResutl(c, utils.RetCode.ExceptionError, err.Error(), size)
			return
		}
		pageInt = pageNum
		sizeInt = sizeNum
	}

	// 获取权限
	quryMap := map[string]string{}

	roleData, err := db.ListRoles(quryMap, keyword, pageInt, sizeInt)

	if err != nil {
		utils.ReturnResutl(c, utils.RetCode.ExceptionError, err.Error(), roleData)
		return
	}

	utils.ReturnResutl(c, utils.RetCode.Success, "", roleData)
}

func ActionRolePut(c *gin.Context) {
	var jsonData map[string]interface{}

	if err := c.ShouldBindJSON(&jsonData); err != nil {
		// 返回错误信息
		utils.ReturnResutl(c, utils.RetCode.ParamRequired, err.Error(), jsonData)
		return
	}

	uuid, ok := jsonData["uuid"]

	if !ok {
		utils.ReturnResutl(c, utils.RetCode.ParamRequired, "缺少参数uuid", jsonData)
		return
	}

	uuidString := uuid.(string)

	role, err := db.QueryRoleByUuid(uuidString)

	if err != nil {
		utils.ReturnResutl(c, utils.RetCode.NotFoundInfo, err.Error(), jsonData)
		return
	}

	err = db.UpdateRole(role.Uuid, jsonData)
	if err != nil {
		utils.ReturnResutl(c, utils.RetCode.NotFoundInfo, err.Error(), jsonData)
		return
	}

	role, _ = db.QueryRoleByUuid(role.Uuid)

	utils.ReturnResutl(c, utils.RetCode.Success, "", role)
}

func ActionRolePost(c *gin.Context) {
	var jsonData models.Role

	if err := c.ShouldBindJSON(&jsonData); err != nil {
		// 返回错误信息
		utils.ReturnResutl(c, utils.RetCode.ParamRequired, err.Error(), jsonData)
		return
	}

	if jsonData.Uuid == "" {
		jsonData.Uuid = utils.GetUuid()
	}

	err := db.Add(&jsonData)
	if err != nil {
		utils.ReturnResutl(c, utils.RetCode.ExceptionError, err.Error(), jsonData)
		return
	}

	utils.ReturnResutl(c, utils.RetCode.Success, "", jsonData)
}

func ActionRoleQuery(c *gin.Context) {
	uuid := c.Param("uuid")

	role, err := db.QueryRoleByUuid(uuid)

	if err != nil {
		utils.ReturnResutl(c, utils.RetCode.NotFoundInfo, err.Error(), role)
		return
	}

	utils.ReturnResutl(c, utils.RetCode.Success, "", role)
}

func ActionRoleDelete(c *gin.Context) {
	uuid := c.Param("uuid")

	role, err := db.QueryRoleByUuid(uuid)

	if err != nil {
		utils.ReturnResutl(c, utils.RetCode.NotFoundInfo, err.Error(), role)
		return
	}

	err = db.DeleteRoleByUuid(uuid)

	if err != nil {
		utils.ReturnResutl(c, utils.RetCode.ExceptionError, err.Error(), role)
		return
	}

	utils.ReturnResutl(c, utils.RetCode.Success, "", role)
}
