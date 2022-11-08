package db

import (
	"gin-template/pkg/models"
)

func QueryUserByUuid(uuid string) (models.User, error) {
	var user models.User
	if err := DB.First(&user, "uuid=?", uuid).Error; err != nil {
		return user, err
	}

	err := DB.Model(new(models.Role)).Select("role.uuid,role.name,role.display_name").Joins("left JOIN user_2_role on user_2_role.role_uuid = role.uuid").Where("user_2_role.user_uuid=?", uuid).Scan(&user.Roles).Error

	return user, err
}

func QueryUserByUsername(username string) (models.User, error) {
	var user models.User
	err := DB.First(&user, "username=?", username).Error
	return user, err
}

func AddUserLoginLog(log models.UserLoginLog) error {
	return DB.Create(&log).Error
}

func DeleteUserByUuid(uuid string) error {
	return DB.Delete(&models.User{}, "uuid=?", uuid).Error
}

func AddUser(user models.User) error {
	return DB.Create(&user).Error
}

func UpdateUser(uuid string, jsonData map[string]interface{}) error {
	return DB.Model(&models.User{}).Where("uuid=?", uuid).Updates(jsonData).Error
}

type userData struct {
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
	Users []*models.User `json:"users"`
}

func ListUsers(params map[string]string, keyword string, page int, size int) (userData, error) {

	data := userData{
		Page: page,
		Size: size,
	}

	db := DB.Where("id>?", 0)

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

	return data, db.Error
}
