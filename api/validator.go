package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/muditshukla3/simplebank/util"
)

var validCurrency validator.Func = func(fieldlevel validator.FieldLevel) bool {
	if currency, ok := fieldlevel.Field().Interface().(string); ok {
		return util.IsCurrencySupported(currency)
	}

	return false
}
