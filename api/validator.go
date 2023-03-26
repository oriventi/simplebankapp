package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/oriventi/simplebank/util"
)

var validatorCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if value, ok := fieldLevel.Field().Interface().(string); ok {
		if util.IsSupportedCurrency(value) {
			return true
		}
	}
	return false
}
