package validator

import (
	"fmt"
	"strings"
)

type Validator struct{}

func Init() *Validator {
	return &Validator{}
}

func (v *Validator) Email(value *string) (bool, string) {
	if value == nil {
		return false, "Email is not provided"
	}

	val := strings.TrimSpace(*value)
	if val == "" {
		return false, "Email is empty"
	}
	if !rgxEmail.MatchString(val) {
		return false, "Email is not valid"
	}

	return true, ""
}

func (v *Validator) Password(value *string) (bool, string) {
	if value == nil {
		return false, "Password is not provided"
	}
	if *value == "" {
		return false, "Password is empty"
	}
	if len(*value) < passMinLength {
		return false, fmt.Sprintf("Password is less than %d characters", passMinLength)
	}
	if len(*value) > passMaxLength {
		return false, fmt.Sprintf("Password is greater than %d characters", passMaxLength)
	}
	if !rgxLowercase.MatchString(*value) {
		return false, "Password must include lowercase letters"
	}
	if !rgxUppercase.MatchString(*value) {
		return false, "Password must include uppercase letters"
	}
	if !rgxNumber.MatchString(*value) {
		return false, "Password must include numbers"
	}
	if !rgxSpecialChars.MatchString(*value) {
		return false, fmt.Sprintf("Password must include at least one special character: %s", specialChars)
	}
	return true, ""
}
