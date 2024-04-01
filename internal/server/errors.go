package server

import "errors"

var (
	ErrHashcashExists = errors.New("hashcash already exists in the store")
	ErrAddrMismatch   = errors.New("hashcash resource and client addr aren't matching")
	ErrDateCheck      = errors.New("didn't pass date check")
	ErrCheckHash      = errors.New("hash does not meet the difficulty criteria")
)
