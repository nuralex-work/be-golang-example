package globfunc

import (
	"be-golang/models/conn"
	"time"

	"github.com/google/uuid"
)

func GetLastCode(table, order, selected string) (res string, err error) {
	tgl := time.Now()

	tx := conn.PostgresDB.Table(table)
	tx.Select(selected).Order(order + " desc")
	tx.Where("extract(month from createddate) = ? and extract(year from createddate) = ?", tgl.Format("01"), tgl.Format("2006"))
	tx.Limit(1)
	tx.Find(&res)

	if tx.Error != nil {
		return "", tx.Error
	}
	return
}

func GetDataById(table, selected, where string, Id uuid.UUID) (res map[string]interface{}, err error) {
	tx := conn.PostgresDB.Table(table)
	tx.Select(selected)
	tx.Where(where, Id)
	tx.Find(&res)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return
}

func GetDataAnyCondition(table string, where map[string]interface{}) (res map[string]interface{}, err error) {
	tx := conn.PostgresDB.Table(table)
	tx.Select("*")

	for i, m := range where {
		tx.Where(i, m)
	}
	tx.Find(&res)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return
}
