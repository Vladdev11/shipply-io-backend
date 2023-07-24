package handlers

import "errors"

var (
	//ErrFailedToReceiveCredentials returned if can't parse credentials from request
	ErrFailedToReceiveCredentials = errors.New("failed to parse credentials")
	//ErrInvalidCredentials returned if credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
	//ErrFailedTokenSignature returned if error while signing a token
	ErrFailedTokenSignature = errors.New("failed to sign token")
	//ErrUserEmailDoesNotExist returned if user with email doens't exist
	ErrUserEmailDoesNotExist = errors.New("user with this email does not exist")
	//ErrResetPasswordTokenGeneration returned if error occurs during password reset token generation
	ErrResetPasswordTokenGeneration = errors.New("failed to create password reset token")
	//ErrTokenRequired returned when token is required, but not in request
	ErrTokenRequired = errors.New("token is required")
	//ErrInvalidToken returned when token is not valid
	ErrInvalidToken = errors.New("invalid token")
	//ErrTokenExpired
	ErrTokenExpired = errors.New("token is expired")
	//ErrInvalidPassword
	ErrInvalidPassword = errors.New("invalid password")
	//ErrInvalidUser
	ErrInvalidUser = errors.New("requested user is not valid")
	//ErrFailedPasswordGeneration
	ErrFailedPasswordGeneration = errors.New("failed to generate password")
	//ErrUpdatePassword
	ErrUpdatePassword = errors.New("failed to update password")
	//ErrDeletePasswordResetToken
	ErrDeletePasswordResetToken = errors.New("failed to delete password reset token")
	//ErrInvalidUserType
	ErrInvalidUserType = errors.New("requesting user is not a valid user type")
	//ErrInvalidShopifyShopName
	ErrInvalidShopifyShopName = errors.New("invalid shop name")
	//ErrNoStoreWithShopifyShopName
	ErrNoStoreWithShopifyShopName = errors.New("no store with this shop name exists")
	//ErrFailedDeleteShopifyAccessToken
	ErrDeleteShopifyAccessToken = errors.New("failed to delete access token")
	//ErrDeleteShopifyLocations
	ErrDeleteShopifyLocations = errors.New("failed to delete shopify locations")
	//ErrShopifyDeactivateStore
	ErrShopifyDeactivateStore = errors.New("failed to deactivate store")
	//ErrRetrieveShopifyGraphqlOrder
	ErrRetrieveShopifyGraphqlOrder = errors.New("failed to retrieve order from shopify")
	//ErrConvertShopifyGraphqlOrder
	ErrConvertShopifyGraphqlOrder = errors.New("failed to convert shopify graphql order")
	//ErrShopifySyncGraphqlOrder
	ErrShopifySyncGraphqlOrder = errors.New("failed to sync shopify graphql order")
	//ErrRetrieveShopifyGraphqlProduct
	ErrRetrieveShopifyGraphqlProduct = errors.New("failed to retrieve product from shopify")
)
