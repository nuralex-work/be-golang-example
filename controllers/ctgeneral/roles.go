package ctgeneral

import (
	"be-golang/helpers"
	"be-golang/models/mdgeneral"
	"be-golang/structs/general"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func GetRoles(e echo.Context) error {
	log.Println("Starting process get data")
	req := new(general.GetRequestDefault)
	byteParam, _ := json.Marshal(helpers.ReqGetParams(e))
	err := json.Unmarshal(byteParam, &req)

	data, err := mdgeneral.GetRoles(req)
	if err != nil {
		return helpers.HandleResponseError(e, http.StatusBadRequest, "Failed", []string{fmt.Sprintf("%v", err)}, nil)
	}
	if data == nil {
		return helpers.HandleResponse(e, http.StatusNotFound, "Data Not Found", nil)
	}

	return helpers.HandleResponse(e, http.StatusOK, "Success", data)
}
