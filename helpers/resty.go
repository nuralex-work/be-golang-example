package helpers

import (
	"be-golang/structs/general"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
)

type RestyClient struct {
	baseUrl string
	client  *resty.Client
}

func NewRestyClient(baseUrl string) RestyClient {
	return RestyClient{
		baseUrl: baseUrl,
		client:  resty.New(),
	}
}

func NewPath(path string, params map[string]string) string {
	newpath := "?"
	if params != nil {
		c := 0
		for i, v := range params {
			if c == 0 {
				newpath += i + "=" + v
			} else {
				newpath += "&" + i + "=" + v
			}
			c++
		}
	} else {
		newpath = ""
	}

	return path + newpath
}

func (resty *RestyClient) Get(path string, header map[string]string, params map[string]string) (string, error) {
	//path = NewPath(path, params)
	//fmt.Println(resty.baseUrl + path, "===== tes")
	resp, err := resty.prepareRequest(header, params, "").Get(resty.baseUrl + path)

	//log trace
	logTraceHttp(resp, err, "GET")

	return string(resp.Body()), err
}

func (resty *RestyClient) Post(path string, header map[string]string, params map[string]string, body interface{}) (string, error) {
	//path = NewPath(path, params)
	var bodyJson map[string]interface{}
	if body != nil {
		bodyByte, _ := json.Marshal(body)
		err := json.Unmarshal(bodyByte, &bodyJson)
		if err != nil {
			return "", nil
		}
	}
	fmt.Println("urls: ", resty.baseUrl+path)
	fmt.Println("params: ", params)
	resp, err := resty.prepareRequest(header, params, bodyJson).Post(resty.baseUrl + path)

	logTraceHttp(resp, err, "POST")

	return string(resp.Body()), err
}

func (resty *RestyClient) PostV2(path string, header map[string]string, params map[string]string, body interface{}, isFormData bool) (map[string]interface{}, error) {
	if isFormData {
		switch val := body.(type) {
		case map[string]string:
			resp, err := resty.prepareRequestV2(header, params, val).Post(resty.baseUrl + path)

			logTraceHttp(resp, err, "POST V2")

			if resp.StatusCode() != 200 {
				return nil, fmt.Errorf("Error on request")
			}

			var result map[string]interface{}
			err = json.Unmarshal(resp.Body(), &result)

			if err != nil {
				fmt.Println("Error:", err)
				return nil, fmt.Errorf("Result is Empty")
			}
			return result, err
		default:
			return nil, fmt.Errorf("Body must be in map[string]string format")
		}
	} else {
		var bodyJson string
		if body != nil {
			bodyByte, _ := json.Marshal(body)
			bodyJson = string(bodyByte)
		}
		resp, err := resty.prepareRequestV2(header, params, bodyJson).Post(resty.baseUrl + path)

		logTraceHttp(resp, err, "POST V2")

		if resp.StatusCode() < 200 && resp.StatusCode() > 201 {
			return nil, fmt.Errorf("Error on request")
		}

		var result map[string]interface{}
		err = json.Unmarshal(resp.Body(), &result)

		if err != nil {
			fmt.Println("Error:", err)
			return nil, fmt.Errorf("Result is Empty")
		}
		return result, err
	}
}
func (resty *RestyClient) PostFile(endpoint string, header map[string]string, params map[string]string, formdata map[string]string, filename string, base64file []byte) (map[string]interface{}, error) {

	decUp := make(map[string]interface{})
	var errResponse error
	for i := 0; i <= 10; i++ {
		log.Println("Looping Upload Image Count : ", i)
		resp, err := resty.prepareRequestFile(header, params, formdata, filename, base64file).Post(resty.baseUrl + endpoint)
		logTraceHttp(resp, err, "POST FILE")
		er := json.Unmarshal(resp.Body(), &decUp)
		if er != nil {
			errResponse = er
			break
		}
		if resp.StatusCode() == 200 {
			if decUp["Code"].(float64) == 200 {
				errResponse = err
				break
			}
		}
		errResponse = err
	}

	return decUp["Data"].(map[string]interface{}), errResponse
}
func (resty *RestyClient) PostOauth2(path string, header map[string]string, params map[string]string, body interface{}) (general.OauthResponse, error) {
	var result general.OauthResponse
	switch val := body.(type) {
	case map[string]string:
		resp, err := resty.prepareRequestV2(header, params, val).Post(resty.baseUrl + path)
		logTraceHttp(resp, err, "POST OAUTH2")

		if resp.StatusCode() != 200 {
			return result, fmt.Errorf("Error on request")
		}

		err = json.Unmarshal(resp.Body(), &result)

		if err != nil {
			fmt.Println("Error:", err)
			return result, fmt.Errorf("Result is Empty")
		}
		return result, err
	default:
		return result, fmt.Errorf("Body must be in map[string]string format")
	}
}

func (resty *RestyClient) Put(path string, header map[string]string, params map[string]string, body interface{}) (string, error) {
	//path = NewPath(path, params)
	fmt.Println("staarting resty put")
	var bodyJson map[string]interface{}
	if body != nil {
		bodyByte, _ := json.Marshal(body)
		err := json.Unmarshal(bodyByte, &bodyJson)
		if err != nil {
			return "", nil
		}
	}
	fmt.Println(bodyJson, "== body")

	resp, err := resty.prepareRequest(header, params, bodyJson).Put(resty.baseUrl + path)
	//
	logTraceHttp(resp, err, "PUT")

	return string(resp.Body()), err
}

func (resty *RestyClient) Delete(path string, header map[string]string, params map[string]string, body interface{}) (string, error) {
	path = NewPath(path, params)
	var bodyJson string
	if body != nil {
		bodyByte, _ := json.Marshal(body)
		bodyJson = string(bodyByte)
	}
	resp, err := resty.prepareRequest(header, params, bodyJson).Delete(resty.baseUrl + path)

	logTraceHttp(resp, err, "DELETE")

	return string(resp.Body()), err
}

func (resty *RestyClient) prepareRequest(header map[string]string, params map[string]string, body interface{}) *resty.Request {
	request := resty.client.R().EnableTrace()
	request = request.SetHeaders(header)
	request = request.SetQueryParams(params)

	if body != "" {
		request = request.SetBody(body)
	}

	//log body request
	log.Println("==================== Request =====================")
	log.Println("Body    : ", body)
	log.Println("Header  : ", header)
	log.Println("Params  : ", params)
	log.Println("==================== Request =====================")

	return request
}
func (resty *RestyClient) prepareRequestV2(header map[string]string, params map[string]string, body interface{}) *resty.Request {
	request := resty.client.R().EnableTrace()
	request = request.SetHeaders(header)
	request = request.SetQueryParams(params)

	switch val := body.(type) {
	case map[string]string:
		request = request.SetFormData(val)
	case string:
		request = request.SetBody(body)
	}

	//log body request
	//log.Println("==================== Request =====================")
	//log.Println("Body    : ", body)
	//log.Println("Header  : ", header)
	//log.Println("Params  : ", params)
	//log.Println("==================== Request =====================")

	return request
}

func (resty *RestyClient) prepareRequestFile(header map[string]string, params map[string]string, formdata map[string]string, filename string, base64file []byte) *resty.Request {
	request := resty.client.R().EnableTrace()
	request = request.SetHeaders(header)
	request = request.SetQueryParams(params)

	if formdata != nil {
		//request = request.SetBody(body)
		request = request.SetFileReader("object", filename, strings.NewReader(string(base64file))).
			SetFormData(formdata)
	}

	//log body request
	log.Println("==================== Request =====================")
	log.Println("Body    	 : ", formdata)
	log.Println("Filename    : ", filename)
	log.Println("Header  	 : ", header)
	log.Println("Params  	 : ", params)
	log.Println("==================== Request =====================")

	return request
}

func logTraceHttp(resp *resty.Response, err error, method string) {
	if err != nil {
		log.Println("=============== Error Http Client ==============")
		log.Println("error : ", err.Error())
		log.Println("=============== Error Http Client ==============")
		return
	}
	//ti := resp.Request.TraceInfo()
	//log.Println("=============== Trace Http Client ==============")
	//log.Println("DNS Lookup    : ", ti.DNSLookup)
	//log.Println("Conn Time     : ", ti.ConnTime)
	//log.Println("Server Time   : ", ti.ServerTime)
	//log.Println("Response Time : ", ti.ResponseTime)
	//log.Println("Total Time	   : ", ti.TotalTime)
	//log.Println("Remote Addr   : ", ti.RemoteAddr)
	//log.Println("=============== End Trace Http Client ==============")

	log.Println("================== Response Http Client ====================")
	log.Println("Status Code : ", resp.StatusCode())
	log.Println("Type : ", method)
	//log.Println("Receive at	: ", resp.ReceivedAt())
	//log.Println("Time        : ", resp.Time())
	log.Println("Error       : ", err)
	//log.Println("Proto       : ", resp.Proto())
	log.Println("Body        : ", string(resp.Body()))
	log.Println("================== End Response Http Client ====================")
}
