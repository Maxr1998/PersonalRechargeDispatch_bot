package main

import (
	"unicode"
)

func isPortalCode(id string) bool {
	if len(id) > 4 {
		return false
	}
	for i, r := range id {
		if (i == 0 && !unicode.IsLetter(r)) || (i > 0 && !unicode.IsDigit(r)) {
			return false
		}
	}
	return true
}
