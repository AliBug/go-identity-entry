package domain

import "errors"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("Internal Server Error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("Your requested Item is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("Item already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("Given Param is not valid")
	// ErrUnauthorized will throw if user not login or token expired
	ErrUnauthorized = errors.New("Unauthorized")
	// ErrForbidden will throw if user has no authz
	ErrForbidden = errors.New("Forbidden")
	// ErrConfig will throw if config has error
	ErrConfig = errors.New("Config Error")
)
