package domain

import (
	"errors"
	"strings"
)

var ErrInvalidRequestMethod = errors.New("Invalid request method")
var ErrZeroRows = errors.New("Zero rows")
var ErrNoFields = errors.New("Failed to find appropriate field(s)")
var ErrTagsNotFound = errors.New("Tag(s) not found")

func IsDataTypeError(errMsg string) bool {
	dataTypeErrors := []string{
		"неверный синтаксис для типа date",
		"invalid input syntax for type",
		"cannot be cast",
		"invalid input value for enum",
	}

	for _, pattern := range dataTypeErrors {
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}

	return false
}
