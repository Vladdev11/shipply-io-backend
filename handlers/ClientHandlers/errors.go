package ClientHandlers

import "errors"

var (
	// ErrProductDoesNotBelongToClient Returned if the user's client doesn't own the requested product
	ErrProductDoesNotBelongToClient = errors.New("product does not belong to client")
)
