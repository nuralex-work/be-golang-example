package models

import (
	"be-golang/helpers"
	"be-golang/models/conn"
	"be-golang/models/globfunc"
	"be-golang/structs"
	"be-golang/structs/general"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GetMasterCustomer(u *structs.GetRequestCustomer) (res map[string]interface{}, err error) {
	var data []map[string]interface{}
	tx := conn.PostgresDB.Table("customers cust")
	tx.Joins("left join users usr on usr.id = cust.salesid")
	tx.Select("cust.*, usr.name salesrep")

	if u.Status != "" && u.Status != "all" {
		tx.Where("(cust.status = ? OR LOWER(cust.status) = ?)", u.Status, u.Status)
	}
	if u.CustomerId != "" {
		tx.Where("LOWER(cust.customerid) LIKE ?", "%"+strings.ToLower(u.CustomerId)+"%")
	}
	if u.Customer != "" {
		tx.Where("LOWER(cust.companyname) LIKE ?", "%"+strings.ToLower(u.Customer)+"%")
	}

	var count int64
	tx.Count(&count)

	tx.Order("cust.customerid")

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
func GetMasterCustomerById(u *general.GetRequestDefault) (res map[string]interface{}, err error) {
	var data map[string]interface{}
	tx := conn.PostgresDB.Table("customers cust")
	tx.Joins("left join users usr on usr.id = cust.salesid")
	tx.Select("cust.*, usr.name salesrep")
	tx.Where("cust.id", u.Id)

	tx.Find(&data)
	return data, nil
}
func CreateMasterCustomer(u *structs.CustomerPayloadCreate) (res interface{}, err error) {
	tx := conn.PostgresDB.Table("customers").Create(&u)
	return u, tx.Error
}
func CreateMasterCustomerMulti(userId string, u []*structs.CustomerPayloadCreate) (res interface{}, err error) {
	counter := 0
	for _, cust := range u {
		cust.Addressee = fmt.Sprintf("%s, %s, %s, %s", cust.Address1, cust.City, cust.State, cust.Zip)
		cust.CreatedId, _ = uuid.Parse(userId)
		cust.CreatedDate = time.Now().Local()
		lastCode, err := globfunc.GetLastCode("customers", "createddate", "internalid")
		if err != nil {
			return nil, err
		}
		dateMonth := time.Now().Format("01")
		dateY := time.Now().Format("06")
		intMonth, _ := strconv.Atoi(dateMonth)
		cust.InternalId = helpers.GenerateCodeSufix(lastCode, helpers.ConvertToRoman(intMonth)+"/"+dateY, 3, "/", counter)
		counter++
	}
	tx := conn.PostgresDB.Table("customers").Create(&u)
	return u, tx.Error
}
func UpdateMasterCustomer(u *structs.CustomerPayload) (res interface{}, err error) {
	var data map[string]interface{}
	tx := conn.PostgresDB.Table("customers").Where("id", u.Id).Updates(&u).Find(&data)
	if len(data) == 0 {
		return nil, tx.Error
	}
	return u, tx.Error
}
func ApprovalMasterCustomer(Id uuid.UUID, req *structs.StatusParam) (res interface{}, err error) {
	var data map[string]interface{}
	tx := conn.PostgresDB.Table("customers").Where("id", Id).Updates(&req).Find(&data)
	if len(data) == 0 {
		return nil, tx.Error
	}
	return data, tx.Error
}
func DeleteMasterCustomer(u *structs.CustomerPayload) (res interface{}, err error) {

	tx := conn.PostgresDB.Table("customers").Where("id", u.Id).Delete(&u)
	return map[string]interface{}{
		"Id":     u.Id,
		"status": true,
	}, tx.Error
}
