package OrganizationHandlers

import "errors"

var (

	// ErrProductDoesNotBelongToOrganization Returned if the user's organization doesn't own the requested product
	ErrProductDoesNotBelongToOrganization = errors.New("product does not belong to organization")
)
