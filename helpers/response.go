package helpers

import (
	"github.com/labstack/echo/v4"
	"log"
)

func HandleError(message string, err interface{}) {
	log.Println("========== Start Error Message ==========")
	log.Println("Message => " + message + ".")
	if err != nil {
		log.Println("Error => ", err)
	}
	log.Println("========== End Of Error Message ==========")
	log.Println()
}

type JSONResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func HandleResponse(e echo.Context, statusCode int, message string, data interface{}) error {
	log.Println("-----------------------------------------------")
	log.Println("status code : ", statusCode)
	log.Println("message : ", message)
	//log.Println("data : ", data)

	dataResponse := &JSONResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}

	return e.JSON(statusCode, dataResponse)
}

type JSONResponseError struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Errors     []string    `json:"errors"`
	Data       interface{} `json:"data"`
}

func HandleResponseError(e echo.Context, statusCode int, message string, errors []string, data interface{}) error {
	log.Println("-----------------------------------------------")
	log.Println("status code : ", statusCode)
	log.Println("message : ", message)
	log.Println("data : ", data)

	dataResponse := &JSONResponseError{
		StatusCode: statusCode,
		Message:    message,
		Errors:     errors,
		Data:       data,
	}

	return e.JSON(statusCode, dataResponse)
}
