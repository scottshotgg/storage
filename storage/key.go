package storage

import "reflect"

// Some keys cannot be supported, like chan, and pointers may be wonky
type Key struct {
	Name  string
	Type  reflect.Kind
	Value interface{}
}
