package validator

import "regexp"

type Validator struct {
	errors map[string]string
}

func New() *Validator {
	return &Validator{errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.errors[key]; !exists {
		v.errors[key] = message
	}
}

//nolint:revive,flag-parameter
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) Errors() map[string]string {
	return v.errors
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
