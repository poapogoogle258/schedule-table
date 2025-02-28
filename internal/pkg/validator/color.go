package validator

import (
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var colorRegex = regexp.MustCompile(`^#[a-fA-F0-9]{6}$`)

func ColorFormat(fl validator.FieldLevel) bool {

	if v := fl.Field(); v.Kind() != reflect.String {
		return false
	}

	value := fl.Field().String()

	if len(value) == 7 && colorRegex.MatchString(value) {
		return true
	}

	return false

}
