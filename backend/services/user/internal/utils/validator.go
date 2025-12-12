package utils

import (
	"fmt"
	"regexp"
	"time"
)

const (
	bioMaxLength      int = 200
	birthdateMinAge   int = 17
	countryMaxLength  int = 50
	labelMaxLength    int = 50
	nameMaxLength     int = 100
	nameMinLength     int = 3
	notesMaxLength    int = 100
	postcodeMaxLength int = 15
	streetMaxLength   int = 250
	subdivMaxLength   int = 250

	maxLatitude  float64 = 90
	minLatitude  float64 = -90
	maxLongitude float64 = 180
	minLongitude float64 = -180
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
		return false, "Name is empty"
	}
	if len(*value) < nameMinLength {
		return false, fmt.Sprintf("Name must be at least %d characters", nameMinLength)
	}
	if len(*value) > nameMaxLength {
		return false, fmt.Sprintf("Name must not exceed %d characters", nameMaxLength)
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
		return false, fmt.Sprintf("Bio must not exceed %d characters", bioMaxLength)
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

	now := time.Now().UTC()
	if value.After(now) {
		return false, fmt.Sprintf("Birthdate is invalid: %s", value.Format("2006-01-02"))
	}

	age := now.Year() - value.Year()
	if now.Month() < value.Month() || (now.Month() == value.Month() && now.Day() < value.Day()) {
		age--
	}
	if age < birthdateMinAge {
		return false, fmt.Sprintf("Must be at least %d years old", birthdateMinAge)
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

func (u *Validator) AddrPhone(value *string, optional bool) (bool, string) {
	if optional && value == nil {
		return true, ""
	}
	if !optional && value == nil {
		return false, "Phone is not provided"
	}
	if *value == "" {
		return false, "Phone is empty"
	}
	if !rgxPhone.MatchString(*value) {
		return false, fmt.Sprintf("Phone is invalid: %s", *value)
	}
	return true, ""
}

func (u *Validator) AddrLabel(value *string, optional bool) (bool, string) {
	if optional && value == nil {
		return true, ""
	}
	if !optional && value == nil {
		return false, "Label is not provided"
	}
	if *value == "" {
		return false, "Label is empty"
	}
	if len(*value) > labelMaxLength {
		return false, fmt.Sprintf("Label must not exceed %d characters", labelMaxLength)
	}
	return true, ""
}

func (u *Validator) AddrNotes(value *string) (bool, string) {
	if value == nil {
		return true, ""
	}
	if *value == "" {
		return true, ""
	}
	if len(*value) > notesMaxLength {
		return false, fmt.Sprintf("Notes must not exceed %d characters", notesMaxLength)
	}
	return true, ""
}

func (u *Validator) AddrCountry(value *string, optional bool) (bool, string) {
	if optional && value == nil {
		return true, ""
	}
	if !optional && value == nil {
		return false, "Country is not provided"
	}
	if *value == "" {
		return false, "Country is empty"
	}
	if len(*value) > countryMaxLength {
		return false, fmt.Sprintf("Country must not exceed %d characters", countryMaxLength)
	}
	return true, ""
}

func (u *Validator) AddrSubdivision(value *string) (bool, string) {
	if value == nil {
		return true, ""
	}
	if *value == "" {
		return true, ""
	}
	if len(*value) > subdivMaxLength {
		return false, fmt.Sprintf("Subdivision must not exceed %d characters", subdivMaxLength)
	}
	return true, ""
}

func (u *Validator) AddrStreet(value *string, optional bool) (bool, string) {
	if optional && value == nil {
		return true, ""
	}
	if !optional && value == nil {
		return false, "Street is not provided"
	}
	if *value == "" {
		return false, "Street is empty"
	}
	if len(*value) > streetMaxLength {
		return false, fmt.Sprintf("Street must not exceed %d characters", streetMaxLength)
	}
	return true, ""
}

func (u *Validator) AddrPostcode(value *string, optional bool) (bool, string) {
	if optional && value == nil {
		return true, ""
	}
	if !optional && value == nil {
		return false, "Postcode is not provided"
	}
	if *value == "" {
		return false, "Postcode is empty"
	}
	if len(*value) > postcodeMaxLength {
		return false, fmt.Sprintf("Postcode must not exceed %d characters", postcodeMaxLength)
	}
	return true, ""
}

func (u *Validator) AddrLatitude(value *float64, optional bool) (bool, string) {
	if optional && value == nil {
		return true, ""
	}
	if !optional && value == nil {
		return false, "Latitude is not provided"
	}
	if *value < minLatitude || *value > maxLatitude {
		return false, fmt.Sprintf("Latitude is invalid: %f", *value)
	}
	return true, ""
}

func (u *Validator) AddrLongitude(value *float64, optional bool) (bool, string) {
	if optional && value == nil {
		return true, ""
	}
	if !optional && value == nil {
		return false, "Longitude is not provided"
	}
	if *value < minLongitude || *value > maxLongitude {
		return false, fmt.Sprintf("Longitude is invalid: %f", *value)
	}
	return true, ""
}
