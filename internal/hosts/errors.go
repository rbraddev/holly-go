package hosts

import "errors"

var (
	ErrSiteRequired         = errors.New("site required")
	ErrMissingInventoryType = errors.New("inventory type is required")
)
