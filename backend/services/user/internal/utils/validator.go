package utils

import (
	"fmt"
	"regexp"
	"time"
)

const (
	nameMinLength int = 3
	nameMaxLength int = 100
	bioMaxLength  int = 200
)

var rgxPhone = regexp.MustCompile(`^(?:0|\+62|62)8\d{8,11}$`)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (u *Validator) Name(value *string, optional bool) (bool, string) {
	if optional && value == nil {
		return true, ""
	}
	if !optional && value == nil {
		return false, "Name is not provided"
	}
	if *value == "" {
		return false, "Name is invalid"
	}
	if len(*value) < nameMinLength {
		return false, fmt.Sprintf("Name is less than %d characters", nameMinLength)
	}
	if len(*value) > nameMaxLength {
		return false, fmt.Sprintf("Name is greater than %d characters", nameMaxLength)
	}
	return true, ""
}

func (u *Validator) Bio(value *string) (bool, string) {
	if value == nil {
		return true, ""
	}
	if *value == "" {
		return true, ""
	}
	if len(*value) > bioMaxLength {
		return false, fmt.Sprintf("Bio is greater than %d characters", bioMaxLength)
	}
	return true, ""
}

func (u *Validator) Sex(value *string) (bool, string) {
	if value == nil {
		return true, ""
	}
	if *value == "" {
		return true, ""
	}
	if *value != "male" && *value != "female" {
		return false, fmt.Sprintf("Sex is invalid: %s", *value)
	}
	return true, ""
}

func (u *Validator) Birthdate(value *time.Time) (bool, string) {
	if value == nil {
		return true, ""
	}
	if value.After(time.Now().UTC()) {
		return false, fmt.Sprintf(
			"Birthdate is invalid: %d-%d-%d",
			value.Year(), value.Month(), value.Day(),
		)
	}
	return true, ""
}

func (u *Validator) Phone(value *string) (bool, string) {
	if value == nil {
		return true, ""
	}
	if *value == "" {
		return true, ""
	}
	if !rgxPhone.MatchString(*value) {
		return false, fmt.Sprintf("Phone is invalid: %s", *value)
	}
	return true, ""
}
