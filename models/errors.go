package models

import (
	"errors"
	"fmt"
)

var (
	// ErrUserNotInContext Returned if the user is not in the request context
	ErrUserNotInContext = errors.New("user not in context")

	// ErrUserIsNotAClient Returned if the user is not a client (i.e. role doesn't match the 2 client roles)
	ErrUserIsNotAClient = errors.New("expected a user that belongs to a client")

	// ErrUserIsNotAnOrganization Returned if the user is not an organization (i.e. role doesn't match the 2 organization roles)
	ErrUserIsNotAnOrganization = errors.New("expected a user that belongs to an organization")

	// ErrUserDoesNotBelongToClient Returned if the user does not belong to the client
	ErrUserDoesNotBelongToClient = errors.New("user does not belong to client")

	// ErrUserDoesNotBelongToOrganization Returned if the user does not belong to the organization
	ErrUserDoesNotBelongToOrganization = errors.New("user does not belong to organization")

	//ErrCarrierNotSupported
	ErrCarrierNotSupported = errors.New("carrier not supported")

	//ErrInvalidJSON
	ErrInvalidJSON = errors.New("failed to parse json")
)

// ErrQueryFailed Returned if the query failed, with the underlying error and the object that was queried
type ErrQueryFailed struct {
	Object string
	Err    error
}

func (e ErrQueryFailed) Error() string {
	return "failed to query db for " + e.Object
}

func (e ErrQueryFailed) Unwrap() error {
	return e.Err
}

// ErrCreateFailed Returned if creation failed, with the underlying error and the object that was supposed to be created
type ErrCreateFailed struct {
	Object string
	Err    error
}

func (e ErrCreateFailed) Error() string {
	return "failed to create " + e.Object
}

func (e ErrCreateFailed) Unwrap() error {
	return e.Err
}

// ErrUpdateFailed Returned if update failed, with the underlying error and the object that was supposed to be updated
type ErrUpdateFailed struct {
	Object string
	Err    error
}

func (e ErrUpdateFailed) Error() string {
	return "failed to update " + e.Object
}

func (e ErrUpdateFailed) Unwrap() error {
	return e.Err
}

// ErrDeleteFailed Returned if deletion failed, with the underlying error and the object that was supposed to be deleted
type ErrDeleteFailed struct {
	Object string
	Err    error
}

func (e ErrDeleteFailed) Error() string {
	return "failed to delete " + e.Object
}

func (e ErrDeleteFailed) Unwrap() error {
	return e.Err
}

// ErrRequiredField Returned if a required field is missing in the request body
type ErrRequiredField struct {
	Field string
}

func (e ErrRequiredField) Error() string {
	return e.Field + " is required"
}

// ErrExists Returned if the object already exists (e.g. a product alias by the same barcode already exists)
type ErrExists struct {
	Object string
}

func (e ErrExists) Error() string {
	return e.Object + " already exists"
}

// ErrInvalidFieldType Returned if a field is not of the expected type
type ErrInvalidFieldType struct {
	Field     string
	FieldType string
}

func (e ErrInvalidFieldType) Error() string {
	return fmt.Sprintf("%s must be of type %s", e.Field, e.FieldType)
}

// ErrInvalidFieldLength Returned if a field is not of the expected length
type ErrInvalidFieldLength struct {
	Field string
	Min   int
	Max   int
}

func (e ErrInvalidFieldLength) Error() string {
	return fmt.Sprintf("%s must be between %d and %d characters", e.Field, e.Min, e.Max)
}

// ErrInvalidFieldRange Returned if a field is not of the expected range
type ErrInvalidFieldRange struct {
	Field string
	Min   int
	Max   int
}

func (e ErrInvalidFieldRange) Error() string {
	return fmt.Sprintf("%s must be between %d and %d", e.Field, e.Min, e.Max)
}

// ErrInvalidFieldValue Returned if a field has an invalid value or does not meet certain criteria
type ErrInvalidFieldValue struct {
	Field  string
	Reason string
}

func (e ErrInvalidFieldValue) Error() string {
	return fmt.Sprintf("%s is invalid: %s", e.Field, e.Reason)
}
