package views

import (
	"gin-template/db"
	"gin-template/models"
	"gin-template/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ActionLogin(c *gin.Context) {

	var jsonData models.User

	if err := c.ShouldBind(&jsonData); err != nil {
		utils.ReturnResutl(c, utils.RetCode.ParamError, err.Error(), map[string]interface{}{})
		return
	}

	if jsonData.Username == "" || jsonData.Password == "" {
		utils.ReturnResutl(c, utils.RetCode.ParamRequired, "用户名或密码不能为空", map[string]interface{}{})
		return
	}

	user, err := db.QueryUserByUsername(jsonData.Username)

	if err != nil {
		utils.ReturnResutl(c, utils.RetCode.NotFoundInfo, "用户不存在", map[string]interface{}{})
		return
	}

	if user.Password != utils.GetMd5Sum(jsonData.Password) {
		utils.ReturnResutl(c, utils.RetCode.LoginError, "密码错误", map[string]interface{}{})
		return
	}

	tokenStr, err := utils.GenerateToken(user.Uuid, c.ClientIP())

	if err != nil {
		utils.ReturnResutl(c, utils.RetCode.ExceptionError, err.Error(), map[string]interface{}{})
		return
	}

	retData := map[string]interface{}{
		"token": tokenStr,
	}

	db.Add(&models.UserLoginLog{
		UserUuid: user.Uuid,
		OptType:  "login",
	})

	utils.ReturnResutl(c, utils.RetCode.Success, "", retData)
}

func ActionLogout(c *gin.Context) {
	userUuid := c.Keys["Uid"].(string)

	user, err := db.QueryUserByUuid(userUuid)

	if err != nil {
		utils.ReturnResutl(c, utils.RetCode.Success, "", map[string]interface{}{})
		return
	}

	db.Add(&models.UserLoginLog{
		UserUuid: user.Uuid,
		OptType:  "logout",
	})

	utils.ReturnResutl(c, utils.RetCode.Success, "", map[string]interface{}{})
}

func ActionUserInfo(c *gin.Context) {
	userUuid := c.Keys["Uid"].(string)

	user, err := db.QueryUserByUuid(userUuid)

	if err != nil {
		utils.ReturnResutl(c, utils.RetCode.NotFoundInfo, err.Error(), map[string]interface{}{})
		return
	}

	roles := []string{}

	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	retData := map[string]interface{}{
		"name":         user.Name,
		"avatar":       user.Avatar,
		"introduction": user.Introduction,
		"roles":        roles,
	}

	utils.ReturnResutl(c, utils.RetCode.Success, "", retData)
}

func ActionUserList(c *gin.Context) {
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

	quryMap := map[string]string{}

	userData, err := db.ListUsers(quryMap, keyword, pageInt, sizeInt)

	if err != nil {
		utils.ReturnResutl(c, utils.RetCode.ExceptionError, err.Error(), userData)
		return
	}

	utils.ReturnResutl(c, utils.RetCode.Success, "", userData)
}

func ActionUserPut(c *gin.Context) {
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

	user, err := db.QueryUserByUuid(uuidString)

	if err != nil {
		utils.ReturnResutl(c, utils.RetCode.NotFoundInfo, err.Error(), jsonData)
		return
	}

	roles, ok := jsonData["roles"].([]interface{})

	if ok {
		db.Unscoped(&models.User2Role{}, "user_uuid=?", user.Uuid)
		for _, role := range roles {

			roleUuid, ok := role.(string)

			if !ok || roleUuid == "" {
				continue
			}
			db.Add(
				&models.User2Role{
					RoleUuid: roleUuid,
					UserUuid: user.Uuid,
				},
			)
		}
	}

	delete(jsonData, "roles")

	err = db.UpdateUser(user.Uuid, jsonData)
	if err != nil {
		utils.ReturnResutl(c, utils.RetCode.NotFoundInfo, err.Error(), jsonData)
		return
	}

	user, _ = db.QueryUserByUuid(user.Uuid)

	utils.ReturnResutl(c, utils.RetCode.Success, "", user)
}

func ActionUserPost(c *gin.Context) {
	var jsonData models.User

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

	for _, role := range jsonData.Roles {
		db.Add(models.User2Role{
			UserUuid: jsonData.Uuid,
			RoleUuid: role.Uuid,
		})
	}

	utils.ReturnResutl(c, utils.RetCode.Success, "", jsonData)
}

func ActionUserQuery(c *gin.Context) {
	uuid := c.Param("uuid")

	user, err := db.QueryUserByUuid(uuid)

	if err != nil {
		utils.ReturnResutl(c, utils.RetCode.NotFoundInfo, err.Error(), user)
		return
	}

	utils.ReturnResutl(c, utils.RetCode.Success, "", user)
}

func ActionUserDelete(c *gin.Context) {
	uuid := c.Param("uuid")

	user, err := db.QueryUserByUuid(uuid)

	if err != nil {
		utils.ReturnResutl(c, utils.RetCode.NotFoundInfo, err.Error(), user)
		return
	}

	err = db.DeleteUserByUuid(uuid)

	if err != nil {
		utils.ReturnResutl(c, utils.RetCode.ExceptionError, err.Error(), user)
		return
	}

	utils.ReturnResutl(c, utils.RetCode.Success, "", user)
}

func ActionUserLoginLog(c *gin.Context) {
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

	quryMap := map[string]string{}

	userData, err := db.ListUserLogs(quryMap, keyword, pageInt, sizeInt)

	if err != nil {
		utils.ReturnResutl(c, utils.RetCode.ExceptionError, err.Error(), userData)
		return
	}

	utils.ReturnResutl(c, utils.RetCode.Success, "", userData)
}
