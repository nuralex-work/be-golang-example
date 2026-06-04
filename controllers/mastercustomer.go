package controllers

import (
	"be-golang/helpers"
	"be-golang/models"
	"be-golang/models/globfunc"
	"be-golang/structs"
	"be-golang/structs/general"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetMasterCustomer(e echo.Context) error {
	log.Println("Starting process get data")
	req := new(structs.GetRequestCustomer)
	byteParam, _ := json.Marshal(helpers.ReqGetParams(e))
	err := json.Unmarshal(byteParam, &req)
	data, err := models.GetMasterCustomer(req)
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	if data == nil {
		return helpers.HandleResponse(e, http.StatusNotFound, "Data Not Found", nil)
	}

	return helpers.HandleResponse(e, http.StatusOK, "Success", data)
}
func GetMasterCustomerById(e echo.Context) error {
	log.Println("Starting process get data by id")

	req := new(general.GetRequestDefault)
	Id, err := uuid.Parse(e.Param("id"))
	req.Id = Id
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	data, err := models.GetMasterCustomerById(req)
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	if data == nil {
		return helpers.HandleResponse(e, http.StatusNotFound, "Data Not Found", nil)
	}

	return helpers.HandleResponse(e, http.StatusOK, "Success", data)
}
func UpdateMasterCustomer(e echo.Context) error {
	log.Println("Starting process put data")

	ReqBody := new(structs.CustomerPayload)
	if err := e.Bind(ReqBody); err != nil {
		log.Println("error : ", err.Error())
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	if ReqBody.SalesId != "" && ReqBody.SalesId != nil {
		_, err := uuid.Parse(ReqBody.SalesId.(string))
		if err != nil {
			return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
		}
	} else {
		ReqBody.SalesId = nil
	}
	Id, err := uuid.Parse(e.Param("id"))
	ReqBody.Id = Id
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	if err := helpers.ValidationCustom(ReqBody, map[string]string{"companyname": "Customer Name"}); err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", err, nil)
	}
	ReqBody.Addressee = fmt.Sprintf("%s, %s, %s, %s", ReqBody.Address1, ReqBody.State, ReqBody.City, ReqBody.Zip)
	data, err := models.UpdateMasterCustomer(ReqBody)
	if err != nil {
		if strings.Contains(fmt.Sprintf("%v", err), "unique_customers_companyname") {
			return helpers.HandleResponseError(e, http.StatusConflict, "Failed", []string{helpers.ValidationStrDuplicate("Customer Name", "Name")}, nil)
		}
		if strings.Contains(fmt.Sprintf("%v", err), "unique_customers_customerid") {
			return helpers.HandleResponseError(e, http.StatusConflict, "Failed", []string{helpers.ValidationStrDuplicate("Customer ID", "ID")}, nil)
		}
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	if data == nil {
		return helpers.HandleResponse(e, http.StatusNotFound, "Data Not Found", nil)
	}

	return helpers.HandleResponse(e, http.StatusOK, "Success", data)
}
func ApprovalMasterCustomer(e echo.Context) error {
	log.Println("Starting process put data")

	ReqBody := new(structs.StatusParam)
	if err := e.Bind(ReqBody); err != nil {
		log.Println("error : ", err.Error())
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	//if err := helpers.Validation(ReqBody); err != nil {
	//	return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", err, nil)
	//}
	ReqBody.Status = "Active"
	Id, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	cust, err := globfunc.GetDataById("customers", "status", "id", Id)
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	if len(cust) == 0 {
		return helpers.HandleResponseError(e, http.StatusNotFound, "Failed", []string{"Customer Not Found"}, nil)
	}
	if cust["status"] != "" && cust["status"] != "Pending Approval" {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{"Customer Status Is Not Pending Approval"}, nil)
	}

	data, err := models.ApprovalMasterCustomer(Id, ReqBody)

	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	if data == nil {
		return helpers.HandleResponse(e, http.StatusNotFound, "Data Not Found", nil)
	}

	return helpers.HandleResponse(e, http.StatusOK, "Success", data)
}
func CreateMasterCustomer(e echo.Context) error {
	log.Println("Starting process create data")
	ReqBody := new(structs.CustomerPayloadCreate)
	if err := e.Bind(ReqBody); err != nil {
		log.Println("error : ", err.Error())
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	if ReqBody.SalesId != "" && ReqBody.SalesId != nil {
		_, err := uuid.Parse(ReqBody.SalesId.(string))
		if err != nil {
			return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
		}
	} else {
		ReqBody.SalesId = nil
	}

	Token := e.Request().Header.Get("Authorization")
	UserId, Roles, ok := helpers.GetDataJwt(Token)
	if ok {
		ReqBody.CreatedId, _ = uuid.Parse(UserId)
	}

	if Roles == "sales-outdoor" {
		ReqBody.Status = "Pending Approval"
	} else {
		ReqBody.Status = "Active"
	}

	ReqBody.CreatedDate = time.Now().Local()

	if err := helpers.ValidationCustom(ReqBody, map[string]string{"companyname": "Customer Name", "customerid": "Customer Id", "resalenumber": "NPWP/NIK"}); err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", err, nil)
	}

	lastCode, err := globfunc.GetLastCode("customers", "createddate", "internalid")
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	dateMonth := time.Now().Format("01")
	dateY := time.Now().Format("06")
	intMonth, _ := strconv.Atoi(dateMonth)
	ReqBody.InternalId = helpers.GenerateCodeSufix(lastCode, helpers.ConvertToRoman(intMonth)+"/"+dateY, 3, "/", 0)
	ReqBody.Addressee = fmt.Sprintf("%s, %s, %s, %s", ReqBody.Address1, ReqBody.City, ReqBody.State, ReqBody.Zip)

	//return helpers.HandleResponse(e, http.StatusOK, "Success", ReqBody)
	data, err := models.CreateMasterCustomer(ReqBody)
	if err != nil {
		if strings.Contains(fmt.Sprintf("%v", err), "unique_customers_companyname") {
			return helpers.HandleResponseError(e, http.StatusConflict, "Failed", []string{helpers.ValidationStrDuplicate("Customer Name", "Name")}, nil)
		}
		if strings.Contains(fmt.Sprintf("%v", err), "unique_customers_customerid") {
			return helpers.HandleResponseError(e, http.StatusConflict, "Failed", []string{helpers.ValidationStrDuplicate("Customer ID", "ID")}, nil)
		}
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}

	if data == nil {
		return helpers.HandleResponse(e, http.StatusNotFound, "Data Not Found", nil)
	}

	return helpers.HandleResponse(e, http.StatusCreated, "Success", data)
}

func CreateMasterCustomerMulti(e echo.Context) error {
	log.Println("Starting process create data")
	ReqBody := new(structs.CustomerPayloadCreateMulti)
	if err := e.Bind(ReqBody); err != nil {
		log.Println("error : ", err.Error())
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	Token := e.Request().Header.Get("Authorization")
	UserId, Roles, ok := helpers.GetDataJwt(Token)
	if !ok {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{"Invalid Token"}, nil)
	}

	if len(ReqBody.CustomerPayloadCreate) > 0 {
		for _, create := range ReqBody.CustomerPayloadCreate {
			if create.SalesId != "" && create.SalesId != nil {
				_, err := uuid.Parse(create.SalesId.(string))
				if err != nil {
					return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
				}
			} else {
				create.SalesId = nil
			}
			if Roles != "admin" {
				create.Status = "Pending Approval"
			} else {
				create.Status = "Active"
			}
			if err := helpers.ValidationCustom(create, map[string]string{"companyname": "Customer Name"}); err != nil {
				return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", err, nil)
			}
		}
	} else {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{helpers.ValidationStr("Customers")}, nil)
	}

	//return helpers.HandleResponse(e, http.StatusOK, "Success", ReqBody)
	data, err := models.CreateMasterCustomerMulti(UserId, ReqBody.CustomerPayloadCreate)
	if err != nil {
		if strings.Contains(fmt.Sprintf("%v", err), "unique_customers_companyname") {
			return helpers.HandleResponseError(e, http.StatusConflict, "Failed", []string{helpers.ValidationStrDuplicate("Customer Name", "Name")}, nil)
		}

		if strings.Contains(fmt.Sprintf("%v", err), "unique_customers_customerid") {
			return helpers.HandleResponseError(e, http.StatusConflict, "Failed", []string{helpers.ValidationStrDuplicate("Customer ID", "ID")}, nil)
		}
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}

	if data == nil {
		return helpers.HandleResponse(e, http.StatusNotFound, "Data Not Found", nil)
	}

	return helpers.HandleResponse(e, http.StatusCreated, "Success", data)
}

func DeleteMasterCustomer(e echo.Context) error {
	log.Println("Starting process delete data")

	req := new(structs.CustomerPayload)
	Id, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	req.Id = Id

	check, err := globfunc.GetDataById("customers", "*", "id", Id)
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}

	if len(check) == 0 {
		return helpers.HandleResponse(e, http.StatusNotFound, "Data Not Found", nil)
	}
	//return helpers.HandleResponse(e, http.StatusOK, "Success", req)

	data, err := models.DeleteMasterCustomer(req)
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	if data == nil {
		return helpers.HandleResponse(e, http.StatusNotFound, "Data Not Found", nil)
	}
	return helpers.HandleResponse(e, http.StatusNoContent, "Success", data)
}
