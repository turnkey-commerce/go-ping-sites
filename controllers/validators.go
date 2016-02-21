package controllers

import (
	"strings"
	"unicode/utf8"
)

func validatePassword(allowMissingPassword bool, password string, valErrors map[string]string) {
	if !allowMissingPassword && utf8.RuneCountInString(strings.TrimSpace(password)) < 6 {
		valErrors["Password"] = "Password must be at least 6 characters in length"
	} else if allowMissingPassword && utf8.RuneCountInString(strings.TrimSpace(password)) < 6 &&
		utf8.RuneCountInString(strings.TrimSpace(password)) > 0 {
		valErrors["Password"] = "Password must be at least 6 characters in length"
	}
}
