package helpers

import (
	"fmt"
	"github.com/go-playground/validator"
	"gorm.io/gorm/utils"
	"regexp"
	"strings"
)

func Validation(req interface{}) (res []string) {
	validate := validator.New()
	// Validate the user
	err := validate.Struct(req)
	if err != nil {
		// Validation failed
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Printf("Please enter value(s) for: %s", err.Field())
			res = append(res, ValidationStr(err.Field()))
		}
		return
	} else {
		// Validation succeeded
		fmt.Println("Data  is valid!")
		return nil
	}
}
func ValidatePassword(password string) (bool, string) {
	var (
		minLength    = 8
		hasUppercase = regexp.MustCompile(`[A-Z]`)
		hasLowercase = regexp.MustCompile(`[a-z]`)
		hasNumber    = regexp.MustCompile(`[0-9]`)
		hasSpecial   = regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",<>\./?\\|]`)
	)
	if len(password) < minLength {
		return false, "minimum length 8"
	}
	// Check for at least one uppercase letter
	if !hasUppercase.MatchString(password) {
		return false, "at least one uppercase letter"
	}
	// Check for at least one lowercase letter
	if !hasLowercase.MatchString(password) {
		return false, "at least one lowercase letter"
	}
	// Check for at least one digit
	if !hasNumber.MatchString(password) {
		return false, "at least one digit number"
	}
	// Check for at least one special character
	if !hasSpecial.MatchString(password) {
		return false, "at least one special character"
	}

	return true, ""
}
func ValidationCustom(req interface{}, custom map[string]string) (res []string) {
	validate := validator.New()
	// Validate the user
	err := validate.Struct(req)
	if err != nil {
		// Validation failed
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Printf("Please enter value(s) for: %s", err.Param())
			res = append(res, ValidationStr(ValidationFormat(custom, err.Field())))
		}
		return
	} else {
		// Validation succeeded
		fmt.Println("Data  is valid!")
		return nil
	}
}
func ValidationFormat(data map[string]string, field string) string {
	str := field
	if len(data) > 0 {
		for i, val := range data {
			if i == strings.ToLower(field) {
				str = val
			}
		}
	}

	return str
}
func ValidationStr(str string) string {
	strSplit := SplitCamelCase(str)
	str = strings.Join(strSplit, " ")
	return fmt.Sprintf("Please enter value(s) for: %s", str)
}
func ValidationStrDelete(str string) string {
	return fmt.Sprintf("This %s cannot be deleted because the %s has associated transactions", str, str)
}
func ValidationStrStatus(str string) string {
	return fmt.Sprintf("This Data cannot be updated because the status already %s", str)
}
func ValidationStrStatusCreate(str string) string {
	return fmt.Sprintf("This Data cannot be created because the status is still %s", str)
}
func ValidationStrDuplicate(str, name string) string {
	return fmt.Sprintf("%s already exists. Please use a different: %s", str, name)
}
func ValidationStrDuplicateStok(str string) string {
	return fmt.Sprintf("This Data cannot be created because the data '%s' is duplicate", str)
}
func ValidationStrPass(str string) string {
	return fmt.Sprintf("Invalid Password Format: %s", str)
}
func ValidationMap(body map[string]interface{}, fieldRequired []string) []string {
	var err []string
	var strbody []string
	for i, val := range body {
		if utils.Contains(fieldRequired, strings.ToLower(i)) && val == "" {
			err = append(err, ValidationStr(i))
		}
		strbody = append(strbody, strings.ToLower(i))
	}

	for _, v := range fieldRequired {
		if !utils.Contains(strbody, strings.ToLower(v)) {
			err = append(err, ValidationStr(v))
		}
	}
	return err
}
