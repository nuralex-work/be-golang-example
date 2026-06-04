package helpers

import (
	"fmt"
	"reflect"
	"strings"
)

func StructToMapNilSafety(data interface{}, custom map[string]string) map[string]interface{} {

	v := reflect.ValueOf(data)
	// Make sure we have a struct
	if v.Kind() == reflect.Ptr {
		v = v.Elem() // get the value the pointer points to
	}
	t := v.Type()
	if v.Kind() != reflect.Struct {
		fmt.Println("Provided value is not a struct ", v.Kind())
		return nil
	}

	var res = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i) // reflect.StructField, contains field name
		value := v.Field(i) // reflect.Value, contains field value

		if value.Interface() != nil && value.Interface() != "" && field.Name != "Id" {
			typeVal := reflect.TypeOf(value.Interface()).String()
			if typeVal == "uint" && value.Interface().(uint) == 0 {
				continue
			}

			s := fmt.Sprintf("%v", value.Interface())
			f := ValidationFormat(custom, strings.ToLower(field.Name))

			//arr := SplitCamelCase(ValidationFormat(custom, field.Name))
			//f := strings.Join(arr, " ")

			res[f] = s
		}
	}
	return res
}
