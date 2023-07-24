package validation

import "errors"

var (
	//ErrInvalidJSON is returned when the JSON is invalid
	ErrInvalidJSON = errors.New("invalid JSON")
)
