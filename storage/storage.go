package storage

import (
	"context"
)

type Storage interface {
	// Name() string
	// Type() storeType
	Get(ctx context.Context, id string) (Item, error)
	GetMulti(ctx context.Context, ids ...string) ([]Item, error)
	GetBy(id, op string, value interface{}, limit int) ([]Item, error)

	Set(id string, i Item, sk map[string]interface{}) error

	Delete(id string) error
	DeleteChangelogs(ids ...string) error

	Iterator() (Iter, error)

	ChangelogIterator() (ChangelogIter, error)
	GetChangelogsForObject(id string) ([]Changelog, error)
	GetLatestChangelogForObject(id string) (*Changelog, error)
}

type Item interface {
	ID() string
	Value() []byte
	Timestamp() int64

	MarshalBinary() (data []byte, err error)
	UnmarshalBinary(data []byte) error
}

type Iter interface {
	Next() (Item, error)
}

type Result struct {
	Item Item
	Err  error
}
