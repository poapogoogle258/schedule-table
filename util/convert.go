package util

import (
	"fmt"

	"github.com/jinzhu/copier"
)

func Convert[T any](formValue interface{}) *T {
	value := new(T)
	if err := copier.Copy(&value, formValue); err != nil {
		panic(err)
	}

	fmt.Println(value, formValue)
	return value
}
