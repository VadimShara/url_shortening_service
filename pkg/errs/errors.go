package errs

import "errors"

var (
	ErrUrlNotFound = errors.New("url not found")
	ErrUrlExists   = errors.New("url exists")
	ErrAliasExists = errors.New("alias exists")
)
