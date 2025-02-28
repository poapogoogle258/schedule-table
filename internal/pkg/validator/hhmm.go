package validator

import (
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var hhmmTime = regexp.MustCompile(`^[0-9]{2}:[0-9]{2}$`)

func HHMMTimeFormat(fl validator.FieldLevel) bool {

	if v := fl.Field(); v.Kind() != reflect.String {
		return false
	}

	value := fl.Field().String()

	if len(value) == 5 && hhmmTime.MatchString(value) {
		return true
	}

	return false

}
