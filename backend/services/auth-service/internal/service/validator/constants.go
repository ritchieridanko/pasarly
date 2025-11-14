package validator

import "regexp"

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
