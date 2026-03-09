package validate

import (
	"net/mail"
	"regexp"
)

var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func Email(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func MinLength(s string, min int) bool {
	return len(s) >= min
}

func Slug(s string) bool {
	return slugRegex.MatchString(s)
}
