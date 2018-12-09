package storage

import "github.com/pizzahutdigital/storage/changelog"

type Storage interface {
	// Name() string
	// Type() storeType
	Get(id string) (Item, error)
	Set(id string, i Item) error
	Delete(id string) error
	Iterator() (Iter, error)
	ChangelogIterator() (ChangelogIter, error)
	GetLatestChangelogForObject(id string) (*changelog.Changelog, error)
}

type Item interface {
	ID() string
	Value() []byte
	MarshalBinary() (data []byte, err error)
	UnmarshalBinary(data []byte) error
	Timestamp() int64
}

type ChangelogIter interface {
	Next() (*changelog.Changelog, error)
}

type Iter interface {
	Next() (Item, error)
}
