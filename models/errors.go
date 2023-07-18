package models

import "errors"

var (
	// ErrUserNotInContext Returned if the user is not in the request context
	ErrUserNotInContext = errors.New("user not in context")

	// ErrUserIsNotAClient Returned if the user is not a client (i.e. role doesn't match the 2 client roles)
	ErrUserIsNotAClient = errors.New("expected a user that belongs to a client")
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
