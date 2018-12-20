package dberrors

import "errors"

const (
	TimeoutString             = "Timed out"
	NotImplementedString      = "Not implemented"
	UnableToParseRPCResString = "Unable to parse RPC response"
	EncodingValueString       = "Could not encode value"
	CouldNotOpenDBString      = "Could not open DB"
	NilDBString               = "Cannot use nil DB"
)

var (
	ErrTimeout             = errors.New(TimeoutString)
	ErrNotImplemented      = errors.New(NotImplementedString)
	ErrUnableToParseRPCRes = errors.New(UnableToParseRPCResString)
	ErrEncodingValue       = errors.New(EncodingValueString)
	ErrCouldNotOpenDB      = errors.New(CouldNotOpenDBString)
	ErrNilDB               = errors.New(NilDBString)
)
