package storage

import (
	"context"
)

type Storage interface {
	ID() string

	Get(ctx context.Context, id string) (Item, error)
	GetBy(ctx context.Context, id, op string, value interface{}, limit int) ([]Item, error)
	GetMulti(ctx context.Context, ids ...string) ([]Item, error)
	GetAll(ctx context.Context) ([]Item, error)

	Set(ctx context.Context, item Item) error
	SetMulti(ctx context.Context, items []Item) error

	Delete(id string) error
	// DeleteBy
	// DeleteMulti
	// DeleteAll() error

	Iterator() (Iter, error)
	IteratorBy(key, op string, value interface{}) (Iter, error)

	// Changelog stuff: move this to it's own file
	GetChangelogsForObject(id string) ([]Changelog, error)
	GetLatestChangelogForObject(id string) (*Changelog, error)

	DeleteChangelogs(ids ...string) error

	ChangelogIterator() (ChangelogIter, error)

	Audit() (map[string]*Changelog, error)
	QuickSync() error
	Sync() error
}

type Result struct {
	Item Item
	Err  error
}
