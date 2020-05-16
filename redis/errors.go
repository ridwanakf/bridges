package redis

import "github.com/pkg/errors"

var (
	// ErrNil means that the key that you specified is empty in redis
	ErrNil    = errors.New("nil returned")
	ErrFailed = errors.New("operation failed")
)
