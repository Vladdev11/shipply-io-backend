package validation

import (
	"errors"
	"strconv"
)

var (
	//ErrInvalidJSON is returned when the JSON is invalid
	ErrInvalidJSON = errors.New("invalid JSON")

	// ErrUnableToReadBody is returned when the request body cannot be read
	ErrUnableToReadBody = errors.New("unable to read request body")
)

// ErrMissingPathParameter is returned when a path parameter is missing
type ErrMissingPathParameter struct {
	Parameter string
}

func (e ErrMissingPathParameter) Error() string {
	return e.Parameter + " not found in URL path"
}

// ErrInvalidPathParameterType is returned when a path parameter is not the expected type
type ErrInvalidPathParameterType struct {
	Parameter    string
	ExpectedType string
}

func (e ErrInvalidPathParameterType) Error() string {
	return e.Parameter + " must be of type " + e.ExpectedType
}

// ErrInvalidPathParameterLength is returned when a path parameter is not the expected length
type ErrInvalidPathParameterLength struct {
	Parameter string
	Min       int
	Max       int
}

func (e ErrInvalidPathParameterLength) Error() string {

	if e.Max == 0 {
		return e.Parameter + " must be at least " + strconv.Itoa(e.Min) + " characters"
	}

	return e.Parameter + " must be between " + strconv.Itoa(e.Min) + " and " + strconv.Itoa(e.Max) + " characters"
}

// ErrInvalidField is returned when a field is invalid
type ErrInvalidField struct {
	Field string
}

func (e ErrInvalidField) Error() string {
	return e.Field + " is invalid"
}

// ErrMissingField is returned when a field is missing
type ErrMissingField struct {
	Field string
}

func (e ErrMissingField) Error() string {
	return e.Field + " is required"
}

// ErrInvalidFieldType is returned when a field is not the expected type
type ErrInvalidFieldType struct {
	Field        string
	ExpectedType string
}

func (e ErrInvalidFieldType) Error() string {
	return e.Field + " must be of type " + e.ExpectedType
}

// ErrInvalidIntFieldRange is returned when an integer field is less than the minimum or greater than the maximum
type ErrInvalidIntFieldRange struct {
	Field string
	Min   int
	Max   int
}

func (e ErrInvalidIntFieldRange) Error() string {

	if e.Max == 0 {
		return e.Field + " must be at least " + strconv.Itoa(e.Min)
	}

	return e.Field + " must be between " + strconv.Itoa(e.Min) + " and " + strconv.Itoa(e.Max)
}

// ErrInvalidFloatFieldRange is returned when a float field is less than the minimum or greater than the maximum
type ErrInvalidFloatFieldRange struct {
	Field string
	Min   float64
	Max   float64
}

func (e ErrInvalidFloatFieldRange) Error() string {

	if e.Max == 0 {
		return e.Field + " must be at least " + strconv.FormatFloat(e.Min, 'f', -1, 64)
	}

	return e.Field + " must be between " + strconv.FormatFloat(e.Min, 'f', -1, 64) + " and " + strconv.FormatFloat(e.Max, 'f', -1, 64)
}

// ErrInvalidStringFieldLength is returned when a string field is less than the minimum or greater than the maximum
type ErrInvalidStringFieldLength struct {
	Field string
	Min   int
	Max   int
}

func (e ErrInvalidStringFieldLength) Error() string {

	if e.Max == 0 {
		return e.Field + " must be at least " + strconv.Itoa(e.Min) + " characters"
	}

	return e.Field + " must be between " + strconv.Itoa(e.Min) + " and " + strconv.Itoa(e.Max) + " characters"
}

// ErrInvalidArrayFieldLength is returned when an array field is less than the minimum or greater than the maximum
type ErrInvalidArrayFieldLength struct {
	Field string
	Min   int
	Max   int
}

func (e ErrInvalidArrayFieldLength) Error() string {

	if e.Max == 0 {
		return e.Field + " must have at least " + strconv.Itoa(e.Min) + " items"
	}

	return e.Field + " must have between " + strconv.Itoa(e.Min) + " and " + strconv.Itoa(e.Max) + " items"
}
