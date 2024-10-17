package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// validationIf is a special validationRuleFunc that does a validation check for a given rule if condition is met.
func validationIf(cond bool, validationRule validation.Rule) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validationRule)
		}

		return nil
	}
}
