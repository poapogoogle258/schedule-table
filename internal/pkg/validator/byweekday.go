package validator

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

func ValidateWeekday(fl validator.FieldLevel) bool {

	if v := fl.Field(); v.Kind() != reflect.Slice {
		return false
	}

	value, ok := fl.Field().Interface().([]int)
	if !ok {
		return false
	}

	for _, v := range value {
		if v < 0 || v > 6 {
			return false
		}
	}

	return true

}

func ValidateMonth(fl validator.FieldLevel) bool {

	if v := fl.Field(); v.Kind() != reflect.Slice {
		return false
	}

	value, ok := fl.Field().Interface().([]int)
	if !ok {
		return false
	}

	for _, v := range value {
		if v < 0 || v > 11 {
			return false
		}
	}

	return true

}
