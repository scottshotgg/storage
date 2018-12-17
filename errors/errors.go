package dberrors

import "errors"

const (
	TimeoutString        = "Timed out"
	NotImplementedString = "Not implemented"
)

var (
	ErrTimeout        = errors.New(TimeoutString)
	ErrNotImplemented = errors.New(NotImplementedString)
)
