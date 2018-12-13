package storage

import (
	"context"
)

// type ItemList interface {
// 	Values() ([]storage.Item, error)
// }

type Storage interface {
	// Name() string
	// Type() storeType
	Get(ctx context.Context, id string) (Item, error)
	GetBy(ctx context.Context, id, op string, value interface{}, limit int) ([]Item, error)
	GetMulti(ctx context.Context, ids ...string) ([]Item, error)
	GetAll(ctx context.Context) ([]Item, error)

	Set(ctx context.Context, item Item) error
	SetMulti(ctx context.Context, items []Item) error

	Delete(id string) error

	Iterator() (Iter, error)
	IteratorBy(key, op string, value interface{}) (Iter, error)

	DeleteChangelogs(ids ...string) error
	ChangelogIterator() (ChangelogIter, error)
	GetChangelogsForObject(id string) ([]Changelog, error)
	GetLatestChangelogForObject(id string) (*Changelog, error)
}

type Item interface {
	ID() string
	Value() []byte
	Timestamp() int64
	Keys() map[string]interface{}

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
