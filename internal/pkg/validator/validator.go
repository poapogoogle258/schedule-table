package validator

import (
	v "github.com/go-playground/validator/v10"
)

var validate *v.Validate

func init() {
	validate = v.New(v.WithRequiredStructEnabled())
	validate.RegisterValidation("telephone", TelephoneFormat)
	validate.RegisterValidation("color", ColorFormat)
	validate.RegisterValidation("hhmm", HHMMTimeFormat)
	validate.RegisterValidation("byweekday", ValidateWeekday)
	validate.RegisterValidation("bymonth", ValidateMonth)
}

func Validate(s any) error {
	return validate.Struct(s)
}
