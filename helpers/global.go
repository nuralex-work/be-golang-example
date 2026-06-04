package helpers

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func FormatIDR(amount int) string {
	var result string
	amountStr := strconv.FormatInt(int64(amount), 10)
	amountLen := len(amountStr)

	for i := 0; i < amountLen; i++ {
		if i%3 == 0 && i != 0 {
			result = "." + result
		}
		result = string(amountStr[amountLen-i-1]) + result
	}

	return result
}
func SplitCamelCase(s string) []string {
	// Regular expression to match camel case
	re := regexp.MustCompile("([a-z])([A-Z])")
	// Replace camel case with a space in between
	s = re.ReplaceAllString(s, "${1} ${2}")
	// Split the string by spaces and return the result
	return strings.Fields(s)
}
func TimeStmpToString(param int64) (result string) {
	timestamp := param
	timeObj := time.Unix(timestamp, 0)
	result = timeObj.Format("2006-01-02 15:04:05")
	return
}

func getLastCode(lastcode, prefix, sufix, sparator string) (code int, jml int) {
	lastcode = strings.ReplaceAll(lastcode, prefix, "")
	lastcode = strings.ReplaceAll(lastcode, sufix, "")
	lastcode = strings.ReplaceAll(lastcode, sparator, "")

	jml = len(lastcode)
	code, _ = strconv.Atoi(lastcode)
	return
}
func GenerateCodePrefix(lastcode string, prefix string, lenchar int, sparator string) (rsp string) {
	code, lastnumber := getLastCode(lastcode, prefix, "", sparator)
	if lenchar > 0 {
		lastnumber = lenchar
	}
	code++
	rsp = prefix + sparator + strings.Repeat("0", lastnumber-len(strconv.Itoa(code))) + strconv.Itoa(code)
	return rsp
}
func GenerateCodeSufix(lastcode string, sufix string, lenchar int, sparator string, increment int) (rsp string) {
	code, lastnumber := getLastCode(lastcode, "", sufix, sparator)
	if lenchar > 0 {
		lastnumber = lenchar
	}
	code++

	code += increment
	rsp = strings.Repeat("0", lastnumber-len(strconv.Itoa(code))) + strconv.Itoa(code) + sparator + sufix
	return rsp
}
func GenerateCodePrefSufixMulti(lastcode string, prefix string, sufix string, lenchar int, sparator string, increment int) (rsp string) {
	code, lastnumber := getLastCode(lastcode, prefix, sufix, sparator)
	if lenchar > 0 {
		lastnumber = lenchar
	}
	code++
	code += increment
	rsp = prefix + sparator + strings.Repeat("0", lastnumber-len(strconv.Itoa(code))) + strconv.Itoa(code) + sparator + sufix
	return rsp
}
func GenerateCodePrefSufix(lastcode string, prefix string, sufix string, lenchar int, sparator string) (rsp string) {
	code, lastnumber := getLastCode(lastcode, prefix, sufix, sparator)
	if lenchar > 0 {
		lastnumber = lenchar
	}
	code++
	rsp = prefix + sparator + strings.Repeat("0", lastnumber-len(strconv.Itoa(code))) + strconv.Itoa(code) + sparator + sufix
	return rsp
}
func RequestGetParams(e echo.Context) (res map[string]interface{}) {
	e.MultipartForm()
	fromParam := make(map[string]interface{})
	fromJson := make(map[string]interface{})

	r := e.Request()
	if err := r.ParseForm(); err != nil {
	}
	jsn, _ := json.Marshal(r.Form)
	if err := json.Unmarshal(jsn, &fromParam); err != nil {
	}

	if err := json.NewDecoder(e.Request().Body).Decode(&fromJson); err != nil {
	}

	if len(fromJson) > 0 {
		res = fromJson
	} else if len(fromParam) > 0 {
		var data = make(map[string]interface{})
		for k, val := range fromParam {
			var afterDec []interface{}
			v, _ := json.Marshal(val)
			if err := json.Unmarshal(v, &afterDec); err != nil {
			}
			data[k] = afterDec[0]
		}
		res = data
	}

	if res == nil {
		res = map[string]interface{}{}
	}

	return res
}
func ReqGetParams(e echo.Context) (res map[string]interface{}) {
	fromParam := make(map[string]interface{})
	r := e.Request()
	if err := r.ParseForm(); err != nil {
	}

	jsn, _ := json.Marshal(r.Form)
	if err := json.Unmarshal(jsn, &fromParam); err != nil {
	}

	if len(fromParam) > 0 {
		var data = make(map[string]interface{})
		for k, val := range fromParam {
			var afterDec []interface{}
			v, _ := json.Marshal(val)
			if err := json.Unmarshal(v, &afterDec); err != nil {
			}
			data[k] = afterDec[0]
		}
		res = data
	}

	if res == nil {
		res = map[string]interface{}{}
	}

	return res
}

func Base64Decode(plaintext string) string {
	text, err := base64.StdEncoding.DecodeString(plaintext)

	if err != nil {
		return ""
	}
	return string(text)
}
func StructToURLValues(data interface{}) url.Values {
	values := url.Values{}
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		tag := field.Tag.Get("url")
		if tag != "" {
			values.Add(tag, fmt.Sprintf("%v", value.Interface()))
		}
	}
	return values
}

func StringToSHA1(text string) string {
	var sha = sha1.New()
	sha.Write([]byte(text))
	var encrypted = sha.Sum(nil)
	var encryptedString = fmt.Sprintf("%x", encrypted)

	return encryptedString
}

func InterfaceToSliceMap(data interface{}) (result []map[string]interface{}) {
	m, _ := json.Marshal(data)
	err := json.Unmarshal(m, &result)
	if err != nil {
		fmt.Println("error conv: ", err)
		return
	}
	return
}
func GetKeyStruct(data interface{}, exclude string) (res []string) {
	var result map[string]interface{}
	jsn, _ := json.Marshal(data)
	json.Unmarshal(jsn, &result)

	for key, _ := range result {
		if !strings.Contains(exclude, key) {
			res = append(res, strings.ToLower(key))
		}

	}
	return
}
func ConvertToRoman(num int) string {
	if num < 1 || num > 3999 {
		return "Input must be between 1 and 3999"
	}

	// Define the mapping of integers to Roman numerals
	romanNumerals := []struct {
		value   int
		numeral string
	}{
		{1000, "M"},
		{900, "CM"},
		{500, "D"},
		{400, "CD"},
		{100, "C"},
		{90, "XC"},
		{50, "L"},
		{40, "XL"},
		{10, "X"},
		{9, "IX"},
		{5, "V"},
		{4, "IV"},
		{1, "I"},
	}

	result := ""
	for _, rn := range romanNumerals {
		// While the number is greater than or equal to the value
		for num >= rn.value {
			result += rn.numeral // Append the Roman numeral
			num -= rn.value      // Decrease the number
		}
	}
	return result
}
func InArray(value string, array []string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}
