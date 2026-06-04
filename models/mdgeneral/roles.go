package mdgeneral

import (
	"be-golang/models/conn"
	"be-golang/structs/general"
	"encoding/json"
	"fmt"
	"strconv"
)

func GetRoles(u *general.GetRequestDefault) (res map[string]interface{}, err error) {
	var data []map[string]interface{}
	tx := conn.PostgresDB.Table("roles")
	tx.Select("*")

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

	var count int64
	tx.Count(&count)

	tx.Order("name desc")

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

	tx.Find(&data)

	var response map[string]interface{}
	var tmp []map[string]interface{}

	if len(data) > 0 {
		for _, val := range data {
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
