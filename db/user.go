package db

import (
	"context"
	"gin-template/models"
)

func QueryUserByUuid(ctx context.Context, uuid string) (models.User, error) {
	var user models.User
	if err := Session.WithContext(ctx).First(&user, "uuid=?", uuid).Error; err != nil {
		return user, err
	}

	err := Session.WithContext(ctx).Model(new(models.Role)).Select("role.uuid,role.name,role.display_name").Joins("left JOIN user_2_role on user_2_role.role_uuid = role.uuid").Where("user_2_role.user_uuid=?", uuid).Scan(&user.Roles).Error

	return user, err
}

func QueryUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User
	err := Session.WithContext(ctx).First(&user, "username=?", username).Error
	return user, err
}

func DeleteUserByUuid(ctx context.Context, uuid string) error {

	if err := Session.WithContext(ctx).Delete(&models.User2Role{}, "user_uuid=?", uuid).Error; err != nil {
		return err
	}
	if err := Session.WithContext(ctx).Delete(&models.UserLoginLog{}, "user_uuid=?", uuid).Error; err != nil {
		return err
	}

	return Session.WithContext(ctx).Delete(&models.User{}, "uuid=?", uuid).Error
}

func UpdateUser(ctx context.Context, uuid string, jsonData map[string]interface{}) error {
	return Session.WithContext(ctx).Model(&models.User{}).Where("uuid=?", uuid).Updates(jsonData).Error
}

type userData struct {
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
	Users []*models.User `json:"users"`
}

func ListUsers(ctx context.Context, params map[string]string, keyword string, page int, size int) (userData, error) {

	data := userData{
		Page: page,
		Size: size,
	}

	db := Session.WithContext(ctx).Where("id>?", 0)

	// 多字段查询
	for key, value := range params {
		db = db.Where(key+"=?", value)
	}
	// 模糊查询
	if keyword != "" {
		db = db.Where("name like ? ", "%"+keyword+"%")
	}

	db = db.Order("id desc")

	db.Model(new(models.User)).Count(&data.Total)

	// 分页
	if page != 0 && size != 0 {
		db = db.Offset((page - 1) * size).Limit(size)
	}

	db = db.Find(&data.Users)

	for _, user := range data.Users {
		Session.WithContext(ctx).Model(new(models.Role)).Select("role.uuid,role.name,role.display_name").Joins("left JOIN user_2_role on user_2_role.role_uuid = role.uuid").Where("user_2_role.user_uuid=?", user.Uuid).Scan(&user.Roles)
	}

	return data, db.Error
}

type userLogData struct {
	Total int64                  `json:"total"`
	Page  int                    `json:"page"`
	Size  int                    `json:"size"`
	Logs  []*models.UserLoginLog `json:"logs"`
}

func ListUserLogs(ctx context.Context, params map[string]string, keyword string, page int, size int) (userLogData, error) {

	data := userLogData{
		Page: page,
		Size: size,
	}

	db := Session.WithContext(ctx).Where("id>?", 0)

	// 多字段查询
	for key, value := range params {
		db = db.Where(key+"=?", value)
	}
	// 模糊查询
	if keyword != "" {
		db = db.Where("name like ? ", "%"+keyword+"%")
	}

	db = db.Order("id desc")

	db.Model(new(models.UserLoginLog)).Count(&data.Total)

	// 分页
	if page != 0 && size != 0 {
		db = db.Offset((page - 1) * size).Limit(size)
	}

	db = db.Find(&data.Logs)

	return data, db.Error
}
