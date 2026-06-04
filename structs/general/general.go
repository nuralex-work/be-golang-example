package general

import (
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"time"
)

type GetRequestDefault struct {
	Id     uuid.UUID `query:"id" form:"id" json:"id" param:"id"`
	Search string    `query:"search" form:"search" json:"search" param:"search"`
	Limit  string    `query:"limit" form:"limit" json:"limit" param:"limit"`
	Offset string    `query:"offset" form:"offset" json:"offset" param:"offset"`
	Sort   string    `query:"sort" form:"sort" json:"sort" param:"sort"`
	Filter string    `query:"filter" form:"filter" json:"filter" param:"filter"`
	Status string    `query:"status" form:"status" json:"status" param:"status"`
}
type GetRequestDefaultMigrate struct {
	Id     uuid.UUID `query:"id" form:"id" json:"id" param:"id"`
	Search string    `query:"search" form:"search" json:"search" param:"search"`
	Limit  int       `query:"limit" form:"limit" json:"limit" param:"limit"`
	Offset int       `query:"offset" form:"offset" json:"offset" param:"offset"`
	Sort   string    `query:"sort" form:"sort" json:"sort" param:"sort"`
	Filter string    `query:"filter" form:"filter" json:"filter" param:"filter"`
	Status string    `query:"status" form:"status" json:"status" param:"status"`
}
type GetRequestCustDate struct {
	GetRequestDefault
	Customer   string `param:"customer" form:"customer" json:"customer"`
	CustomerId string `param:"customerid" form:"customerid" json:"customerid"`
	Entity     string `param:"entity" form:"entity" json:"entity"`
	StartDate  string `param:"startdate" form:"startdate" json:"startdate"`
	EndDate    string `param:"enddate" form:"enddate" json:"enddate"`
	TranId     string `param:"tranid" form:"tranid" json:"tranid"`
}
type GetRequestMigrate struct {
	GetRequestDefaultMigrate
	Customer   string `param:"customer" form:"customer" json:"customer"`
	CustomerId string `param:"customerid" form:"customerid" json:"customerid"`
	Entity     string `param:"entity" form:"entity" json:"entity"`
	StartDate  string `param:"startdate" form:"startdate" json:"startdate"`
	EndDate    string `param:"enddate" form:"enddate" json:"enddate"`
	TranId     string `param:"tranid" form:"tranid" json:"tranid"`
}
type GetRequestItemDate struct {
	GetRequestDefault
	DisplayName string `param:"displayname" form:"displayname" json:"displayname"`
	ItemId      string `param:"itemid" form:"itemid" json:"itemid"`
	StartDate   string `param:"startdate" form:"startdate" json:"startdate"`
	EndDate     string `param:"enddate" form:"enddate" json:"enddate"`
}
type GetRequestAgreementDate struct {
	GetRequestDefault
	AgreementName string `param:"agreementname" form:"agreementname" json:"agreementname"`
	AgreementCode string `param:"agreementcode" form:"agreementcode" json:"agreementcode"`
	StartDate     string `param:"startdate" form:"startdate" json:"startdate"`
	EndDate       string `param:"enddate" form:"enddate" json:"enddate"`
}
type GetRequestDefCust struct {
	GetRequestDefault
	Customer string `query:"customerid" form:"customerid" json:"customerid" param:"customerid"`
	Type     string `query:"type" form:"type" json:"type" param:"type"`
}
type CreateUserTime struct {
	CreatedId   uuid.UUID `gorm:"column:createdby;type:uuid" param:"createdby" form:"createdby" json:"createdby"`
	CreatedDate time.Time `gorm:"column:createddate" param:"createddate" form:"createddate" json:"createddate"`
}
type ApprovalUserTime struct {
	ApprovedBy   uuid.UUID `gorm:"column:approvedby;type:uuid" param:"approvedby" form:"approvedby" json:"approvedby"`
	ApprovedDate time.Time `gorm:"column:approveddate" param:"approveddate" form:"approveddate" json:"approveddate"`
}
type GetSelectParams struct {
	Id   uuid.UUID `query:"id" form:"id" json:"id"`
	Name string    `query:"name" form:"name" json:"name"`
}

type OauthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	CreatedAt   string `json:"created_at"`
	Exp         string `json:"exp"`
	TokenType   string `json:"token_type"`
	Duration    string `json:"duration"`
	Roles       string `json:"roles"`
}

type CustomValidator struct {
	Validator *validator.Validate
}

type GetRequestDefaultV2 struct {
	Id        int    `query:"id" form:"id" json:"id" param:"id"`
	Search    string `query:"search" form:"search" json:"search" param:"search"`
	Limit     string `query:"limit" form:"limit" json:"limit" param:"limit"`
	Offset    string `query:"offset" form:"offset" json:"offset" param:"offset"`
	Sort      string `query:"sort" form:"sort" json:"sort" param:"sort"`
	Filter    string `query:"filter" form:"filter" json:"filter" param:"filter"`
	Status    string `query:"status" form:"status" json:"status" param:"status"`
	StartDate string `param:"startdate" form:"startdate" json:"startdate"`
	EndDate   string `param:"enddate" form:"enddate" json:"enddate"`
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}
