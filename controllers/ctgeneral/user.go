package ctgeneral

import (
	"be-golang/helpers"
	"be-golang/models/globfunc"
	"be-golang/models/mdgeneral"
	"be-golang/structs/general"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func GetUser(e echo.Context) error {
	log.Println("Starting process get data")
	req := new(general.GetUserParams)
	byteParam, _ := json.Marshal(helpers.ReqGetParams(e))
	err := json.Unmarshal(byteParam, &req)

	data, err := mdgeneral.GetUser(req)
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	if data == nil {
		return helpers.HandleResponse(e, http.StatusNotFound, "Data Not Found", nil)
	}

	return helpers.HandleResponse(e, http.StatusOK, "Success", data)
}
func GetUserById(e echo.Context) error {
	log.Println("Starting process get data by id")

	req := new(general.GetRequestDefault)
	Id, err := uuid.Parse(e.Param("id"))
	req.Id = Id
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	data, err := mdgeneral.GetUserById(req)
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	if data == nil {
		return helpers.HandleResponse(e, http.StatusNotFound, "Data Not Found", nil)
	}

	return helpers.HandleResponse(e, http.StatusOK, "Success", data)
}
func GetUserProfile(e echo.Context) error {
	log.Println("Starting process get data profil")
	req := new(general.GetRequestDefault)
	Token := e.Request().Header.Get("Authorization")
	UserId, _, ok := helpers.GetDataJwt(Token)
	var SalesId uuid.NullUUID
	if ok {
		usrId, _ := uuid.Parse(UserId)
		SalesId = uuid.NullUUID{UUID: usrId, Valid: true}
	}

	if SalesId.Valid {
		req.Id = SalesId.UUID
	}
	data, err := mdgeneral.GetUserById(req)
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	if data == nil {
		return helpers.HandleResponse(e, http.StatusNotFound, "Data Not Found", nil)
	}

	return helpers.HandleResponse(e, http.StatusOK, "Success", data)
}
func UpdateUser(e echo.Context) error {
	log.Println("Starting process put data")

	ReqBody := new(general.User)
	if err := e.Bind(ReqBody); err != nil {
		log.Println("error : ", err.Error())
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	Id, err := uuid.Parse(e.Param("id"))
	ReqBody.ID = Id
	//return helpers.HandleResponse(e, http.StatusOK, "Success", ReqBody)
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}

	if err := helpers.ValidationCustom(ReqBody, map[string]string{"roleid": "Role ID"}); err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", err, nil)
	}
	if ok, str := helpers.ValidatePassword(ReqBody.Password); !ok && ReqBody.Password != "" {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{helpers.ValidationStrPass(str)}, nil)
	}
	data, err := mdgeneral.UpdateUser(ReqBody)
	if err != nil {
		if strings.Contains(fmt.Sprintf("%v", err), "unique_users_username") {
			return helpers.HandleResponseError(e, http.StatusConflict, "Failed", []string{helpers.ValidationStrDuplicate("Username", "Username")}, nil)
		}
		if strings.Contains(fmt.Sprintf("%v", err), "unique_users_name") {
			return helpers.HandleResponseError(e, http.StatusConflict, "Failed", []string{helpers.ValidationStrDuplicate("Name", "Name")}, nil)
		}
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	if data == nil {
		return helpers.HandleResponse(e, http.StatusNotFound, "Data Not Found", nil)
	}

	return helpers.HandleResponse(e, http.StatusOK, "Success", data)
}
func CreateUser(e echo.Context) error {
	log.Println("Starting process create data")
	ReqBody := new(general.UserPayload)
	if err := e.Bind(ReqBody); err != nil {
		log.Println("error : ", err.Error())
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}

	Token := e.Request().Header.Get("Authorization")
	UserId, ok := helpers.GetUserId(Token)
	if ok {
		ReqBody.CreatedId, _ = uuid.Parse(UserId)
	}
	ReqBody.CreatedDate = time.Now().Local()

	if err := helpers.ValidationCustom(ReqBody, map[string]string{"roleid": "Role ID"}); err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", err, nil)
	}
	if ok, str := helpers.ValidatePassword(ReqBody.Password); !ok {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{helpers.ValidationStrPass(str)}, nil)
	}
	//return helpers.HandleResponse(e, http.StatusCreated, "Success", ReqBody)
	data, err := mdgeneral.CreateUser(ReqBody)
	if err != nil {
		if strings.Contains(fmt.Sprintf("%v", err), "unique_users_username") {
			return helpers.HandleResponseError(e, http.StatusConflict, "Failed", []string{helpers.ValidationStrDuplicate("Username", "Username")}, nil)
		}
		if strings.Contains(fmt.Sprintf("%v", err), "unique_users_name") {
			return helpers.HandleResponseError(e, http.StatusConflict, "Failed", []string{helpers.ValidationStrDuplicate("Name", "Name")}, nil)
		}
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}

	if data == nil {
		return helpers.HandleResponse(e, http.StatusNotFound, "Data Not Found", nil)
	}

	return helpers.HandleResponse(e, http.StatusCreated, "Success", data)
}
func CreateUserMulti(e echo.Context) error {
	log.Println("Starting process create data")
	ReqBody := new(general.UserPayloadBulk)
	if err := e.Bind(ReqBody); err != nil {
		log.Println("error : ", err.Error())
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	//return helpers.HandleResponse(e, http.StatusCreated, "Success", ReqBody)
	Token := e.Request().Header.Get("Authorization")
	UserId, ok := helpers.GetUserId(Token)
	if !ok {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{"Invalid Token"}, nil)
	}
	if len(ReqBody.UserPayload) > 0 {
		for _, create := range ReqBody.UserPayload {
			if err := helpers.ValidationCustom(create, map[string]string{"roleid": "Role ID"}); err != nil {
				return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", err, nil)
			}
		}
	} else {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{helpers.ValidationStr("Items")}, nil)
	}

	data, err := mdgeneral.CreateUserMulti(UserId, ReqBody)
	if err != nil {
		if strings.Contains(fmt.Sprintf("%v", err), "unique_users_username") {
			return helpers.HandleResponseError(e, http.StatusConflict, "Failed", []string{helpers.ValidationStrDuplicate("Username", "Username")}, nil)
		}
		if strings.Contains(fmt.Sprintf("%v", err), "unique_users_name") {
			return helpers.HandleResponseError(e, http.StatusConflict, "Failed", []string{helpers.ValidationStrDuplicate("Name", "Name")}, nil)
		}
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}

	if data == nil {
		return helpers.HandleResponse(e, http.StatusNotFound, "Data Not Found", nil)
	}

	return helpers.HandleResponse(e, http.StatusCreated, "Success", data)
}
func DeleteUser(e echo.Context) error {
	log.Println("Starting process delete data")
	Id, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	checkdefault, err := globfunc.GetDataAnyCondition("users", map[string]interface{}{
		"id":        Id,
		"isdefault": true,
	})

	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	if len(checkdefault) > 0 {
		return helpers.HandleResponse(e, http.StatusBadRequest, "Can't Delete Default User", nil)
	}

	check, err := globfunc.GetDataById("users", "id", "id", Id)
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}

	if len(check) == 0 {
		return helpers.HandleResponse(e, http.StatusNotFound, "Data Not Found", nil)
	}

	data, err := mdgeneral.DeleteUser(Id)
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	if data == nil {
		return helpers.HandleResponse(e, http.StatusNotFound, "Data Not Found", nil)
	}
	return helpers.HandleResponse(e, http.StatusNoContent, "Success", data)
}
func GetUserSales(e echo.Context) error {
	log.Println("Starting process get data")
	req := new(general.GetUserParams)
	byteParam, _ := json.Marshal(helpers.ReqGetParams(e))
	err := json.Unmarshal(byteParam, &req)

	data, err := mdgeneral.GetUserSales(req)
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	if data == nil {
		return helpers.HandleResponse(e, http.StatusNotFound, "Data Not Found", nil)
	}

	return helpers.HandleResponse(e, http.StatusOK, "Success", data)
}
