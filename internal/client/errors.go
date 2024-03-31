package client

import "errors"

var (
	ErrMaxIterExceed = errors.New("maximum iterations exceeded")
	ErrDateCheck     = errors.New("didn't pass date check")
)
