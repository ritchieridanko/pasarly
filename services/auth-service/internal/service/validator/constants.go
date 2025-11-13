package validator

import "regexp"

const (
	passMinLength int    = 8
	passMaxLength int    = 50
	specialChars  string = `!@#$%^&*(),.?":{}|<>`
)

var (
	rgxLowercase    = regexp.MustCompile(`[a-z]`)
	rgxUppercase    = regexp.MustCompile(`[A-Z]`)
	rgxNumber       = regexp.MustCompile(`[0-9]`)
	rgxSpecialChars = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
)
