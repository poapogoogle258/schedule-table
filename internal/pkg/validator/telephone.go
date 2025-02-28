package validator

import (
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var telephoneRegex = regexp.MustCompile(`^[0-9]{10}$`)

func TelephoneFormat(fl validator.FieldLevel) bool {

	if v := fl.Field(); v.Kind() != reflect.String {
		return false
	}

	value := fl.Field().String()

	if len(value) == 10 && telephoneRegex.MatchString(value) {
		return true
	}

	return false

}
