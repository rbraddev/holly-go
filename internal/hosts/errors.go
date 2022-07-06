package hosts

import "errors"

var (
	ErrSiteRequired         = errors.New("site required")
	ErrMissingInventoryType = errors.New("inventory type is required")
	ErrDriverNotSet         = errors.New("network driver not set")
)
