package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

func validationIf(cond bool, validationRule validation.Rule) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validationRule)
		}

		return nil
	}
}
