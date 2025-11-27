package utils

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	passMinLength int    = 8
	passMaxLength int    = 50
	specialChars  string = `!@#$%^&*(),.?":{}|<>`
)

var (
	rgxEmail        = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	rgxLowercase    = regexp.MustCompile(`[a-z]`)
	rgxNumber       = regexp.MustCompile(`[0-9]`)
	rgxSpecialChars = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	rgxUppercase    = regexp.MustCompile(`[A-Z]`)
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (u *Validator) Email(value *string) (bool, string) {
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

func (u *Validator) Password(value *string) (bool, string) {
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

func (u *Validator) Token(value *string) (bool, string) {
	if value == nil {
		return false, "Token is not provided"
	}

	val := strings.TrimSpace(*value)
	if val == "" {
		return false, "Token is empty"
	}

	return true, ""
}
