package db

import "gin-template/models"

func QueryRoleByUuid(uuid string) (models.Role, error) {
	var role models.Role
	err := Session.First(&role, "uuid=?", uuid).Error

	return role, err
}

func DeleteRoleByUuid(uuid string) error {
	return Session.Delete(&models.Role{}, "uuid=?", uuid).Error
}

func AddRole(role models.Role) error {
	return Session.Create(&role).Error
}

func UpdateRole(uuid string, jsonData map[string]interface{}) error {
	return Session.Model(&models.Role{}).Where("uuid=?", uuid).Updates(jsonData).Error
}

type roleData struct {
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
	Roles []*models.Role `json:"roles"`
}

func ListRoles(params map[string]string, keyword string, page int, size int) (roleData, error) {

	data := roleData{
		Page: page,
		Size: size,
	}

	db := Session.Where("id>?", 0)

	// 多字段查询
	for key, value := range params {
		db = db.Where(key+"=?", value)
	}
	// 模糊查询
	if keyword != "" {
		db = db.Where("name like ? ", "%"+keyword+"%")
	}

	db = db.Order("id desc")

	db.Model(new(models.Role)).Count(&data.Total)

	// 分页
	if page != 0 && size != 0 {
		db = db.Offset((page - 1) * size).Limit(size)
	}

	db = db.Find(&data.Roles)

	return data, db.Error
}
