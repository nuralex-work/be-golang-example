package mdgeneral

import (
	"be-golang/helpers"
	"be-golang/models/conn"
	"be-golang/structs/general"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

func GetUser(u *general.GetUserParams) (res map[string]interface{}, err error) {
	var data []map[string]interface{}
	tx := conn.PostgresDB.Table("users usr")
	tx.Select("usr.*, roles.name role_name")
	tx.Joins("left join roles on roles.id = usr.roleid")

	var filter map[string]interface{}
	err = json.Unmarshal([]byte(u.Filter), &filter)
	if len(filter) > 0 {
		for index, element := range filter {
			value := fmt.Sprintf("%v", element)
			if index == "id" {
				tx.Where(index+"=?", value)
			} else {
				tx.Where(index+" LIKE ?", "%"+value+"%")
			}
		}
	}

	if u.Name != "" {
		tx.Where("usr.name LIKE ? OR LOWER(usr.name) LIKE ?", "%"+u.Name+"%", "%"+strings.ToLower(u.Name)+"%")
	}

	if u.Role != "" {
		u.Role = strings.ReplaceAll(u.Role, "_", "-")
		tx.Where("roles.name = ? OR LOWER(roles.name) = ?", u.Role, strings.ToLower(u.Role))
	}
	if u.Status != "" {
		tx.Where("usr.status = ? OR LOWER(usr.status) = ?", u.Status, strings.ToLower(u.Status))
	}
	var count int64
	tx.Count(&count)

	limit := 20
	if u.Limit != "" {
		limit, _ = strconv.Atoi(u.Limit)
		tx.Limit(limit)
	}

	offset := 0
	if u.Offset != "" {
		page, err := strconv.Atoi(u.Offset)
		if err != nil {
			page = 0
		}
		if page > 0 {
			offset = (page - 1) * limit
			tx.Offset(offset)
		}
	}

	tx.Order("createddate desc")
	tx.Find(&data)

	var response map[string]interface{}
	var tmp []map[string]interface{}

	if len(data) > 0 {
		for _, val := range data {
			val["password"] = ""
			tmp = append(tmp, val)
		}
	} else {
		tmp = []map[string]interface{}{}
	}

	response = map[string]interface{}{
		"list":        tmp,
		"total_items": count,
	}

	return response, nil
}
func GetUserById(u *general.GetRequestDefault) (res map[string]interface{}, err error) {
	var data map[string]interface{}
	tx := conn.PostgresDB.Table("users usr")
	tx.Select("usr.*, roles.name as role_name")
	tx.Joins("left join roles on roles.id = usr.roleid")
	tx.Where("usr.id", u.Id)

	tx.Find(&data)

	if len(data) > 0 {
		data["password"] = ""
	}

	return data, nil
}
func CreateUser(u *general.UserPayload) (res interface{}, err error) {
	if u.Password != "" {
		u.Password = helpers.StringToSHA1(u.Password)
	}
	u.Status = "active"
	tx := conn.PostgresDB.Table("users").Create(&u)
	return u, tx.Error
}

func CreateUserMulti(userId string, u *general.UserPayloadBulk) (res interface{}, err error) {
	for _, item := range u.UserPayload {
		if item.Password != "" {
			item.Password = helpers.StringToSHA1(item.Password)
		}
		item.Status = "active"
		item.CreatedId, _ = uuid.Parse(userId)
		item.CreatedDate = time.Now().Local()
	}

	tx := conn.PostgresDB.Table("users").Create(&u)
	return u, tx.Error
}

func UpdateUser(u *general.User) (res interface{}, err error) {
	var payload map[string]interface{}
	byteParam, _ := json.Marshal(u)
	err = json.Unmarshal(byteParam, &payload)
	if u.Password == "" {
		delete(payload, "password")
	} else {
		payload["password"] = helpers.StringToSHA1(u.Password)
	}
	if u.Status == "" {
		delete(payload, "status")
	}
	var data map[string]interface{}
	tx := conn.PostgresDB.Table("users").Where("id", u.ID).Updates(&payload).Find(&data)
	if len(data) == 0 {
		return nil, tx.Error
	}
	return data, tx.Error
}

func DeleteUser(Id uuid.UUID) (res interface{}, err error) {
	tx := conn.PostgresDB.Table("users").Where("id", Id).Updates(&map[string]interface{}{
		"status": "inactive",
	})
	return map[string]interface{}{
		"Id":     Id,
		"status": "inactive",
	}, tx.Error
}

func GetUserSales(u *general.GetUserParams) (res map[string]interface{}, err error) {
	var data []map[string]interface{}
	tx := conn.PostgresDB.Table("users usr")
	tx.Select("usr.id, usr.name, roles.name as role")
	tx.Joins("left join roles on roles.id = usr.roleid")
	tx.Where("roles.name IN (?)", []string{"sales-outdoor", "sales-indoor", "admin"})
	if u.Name != "" {
		tx.Where("usr.name LIKE ? OR LOWER(usr.name) LIKE ?", "%"+u.Name+"%", "%"+strings.ToLower(u.Name)+"%")
	}
	tx.Where("usr.status = ? OR LOWER(usr.status) = ?", "Active", "active")
	tx.Where("usr.isdefault = ?", false)

	var count int64
	tx.Count(&count)

	if u.Limit != "" {
		limit, _ := strconv.Atoi(u.Limit)
		tx.Limit(limit)
	}
	if u.Offset != "" {
		offset, _ := strconv.Atoi(u.Offset)
		tx.Offset(offset)
	}

	tx.Order("createddate desc")
	tx.Find(&data)

	var response map[string]interface{}
	var tmp []map[string]interface{}

	if len(data) > 0 {
		for _, val := range data {
			val["role"] = strings.ReplaceAll(val["role"].(string), "-", " ")
			tmp = append(tmp, val)
		}
	} else {
		tmp = []map[string]interface{}{}
	}

	response = map[string]interface{}{
		"list":        tmp,
		"total_items": count,
	}

	return response, nil
}
