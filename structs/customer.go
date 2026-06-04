package structs

import (
	"be-golang/structs/general"
	"github.com/google/uuid"
)

type CustomerNs struct {
	EntityId string `gorm:"column:entityid;size:45" param:"entityid" form:"entityid" json:"entityid,omitempty"`
}
type Customer struct {
	CustomerPayload
	CustomerNs
	StatusParam
	general.CreateUserTime
}

type StatusParam struct {
	Status string `gorm:"column:status;default:'active'" param:"status" form:"status" json:"status" validate:"required"`
}
type CustomerPayload struct {
	Id           uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" param:"id" form:"id" json:"id,omitempty"`
	InternalId   string      `gorm:"column:internalid;size:45" param:"internalid" form:"internalid" json:"internalid,omitempty"`
	CustomerId   string      `gorm:"column:customerid;size:45" param:"customerid" form:"customerid" json:"customerid" validate:"required"`
	CompanyName  string      `gorm:"column:companyname;size:100" param:"companyname" form:"companyname" json:"companyname" validate:"required"`
	Email        string      `gorm:"column:email;size:100;index:idx_email;unique" param:"email" form:"email" json:"email"`
	Phone        string      `gorm:"column:phone;size:45;index:idx_phone;unique" param:"phone" form:"phone" json:"phone"`
	TelePhone    string      `gorm:"column:altphone;size:45" param:"altphone" form:"altphone" json:"altphone"`
	State        string      `gorm:"column:state;size:100" param:"state" form:"state" json:"state" validate:"required"`
	City         string      `gorm:"column:city;size:100" param:"city" form:"city" json:"city" validate:"required"`
	Addressee    string      `gorm:"column:addressee" param:"addressee" form:"addressee" json:"addressee"`
	Address1     string      `gorm:"column:addr1" param:"addr1" form:"addr1" json:"addr1" validate:"required"`
	Zip          string      `gorm:"column:zip;size:10" param:"zip" form:"zip" json:"zip"`
	Terms        string      `gorm:"column:terms" param:"terms" form:"terms" json:"terms"`
	CreditLimit  uint        `gorm:"column:creditlimit" param:"creditlimit" form:"creditlimit" json:"creditlimit"`
	ResaleNumber string      `gorm:"column:resalenumber;size:100" param:"resalenumber" form:"resalenumber" json:"resalenumber" validate:"required"`
	SalesId      interface{} `gorm:"type:uuid;column:salesid;" param:"salesid" form:"salesid" json:"salesid"`
	Category     string      `gorm:"type:uuid;column:category;" param:"category" form:"category" json:"category"`
}
type CustomerPayloadResp struct {
	CustomerPayload
	SalesRep string `gorm:"type:uuid;column:salesrep;" param:"salesrep" form:"salesrep" json:"salesrep"`
}
type CustomerPayloadCreate struct {
	CustomerPayload
	StatusParam
	general.CreateUserTime
}
type CustomerPayloadCreateMulti struct {
	CustomerPayloadCreate []*CustomerPayloadCreate `param:"customers" form:"customers" json:"customers"`
}
type GetRequestCustomer struct {
	general.GetRequestDefault
	CustomerId string `param:"customerid" form:"customerid" json:"customerid"`
	Customer   string `param:"customer" form:"customer" json:"customer"`
}
