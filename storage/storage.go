package storage

import (
	"context"

	pb "github.com/scottshotgg/storage/protobufs"
)

// Might try doing this

type Metadata interface {
	ID() string
	Created() int64
	LastAccessed() int64
	LastAudit() int64
	LastSync() int64
	Count() int64
	Tables() []string
	Keys() []string
	Accuracy() float64
	ErrorRate() float64

	Values() map[string]interface{}
}

// Think about this for a bit

type Value interface {
	// Require Gob encoding/decoding to be implemented
	GobEncode() ([]byte, error)
	GobDecode([]byte) error
}

type Filter struct {
	key   string
	op    string
	value Value

	keysOnly bool
}

type Query struct {
	IDs     []string
	Filters []Filter
	Limit   int64
}

type Storage interface {
	ID() string
	// Metadata() *Metadata
	// Open() error
	// OpenWith() error
	// Close() error
	// New() (Storage, error)
	// NewWith(config Config) (Storage, error)

	// TODO: try doing this to encompass all of the implementations
	// Get(ctx context.Context, query Query) ([]Item, error)

	Get(ctx context.Context, itemIDs []string) (pb.Item, error)
	// GetBy(ctx context.Context, key, op string, value interface{}, limit int) ([]Item, error)
	// GetMulti(ctx context.Context, ids []string) ([]Item, error)
	// GetAll(ctx context.Context) ([]Item, error)

	Set(ctx context.Context, items []pb.Item) error
	// SetMulti(ctx context.Context, items []Item) error

	Delete(ctx context.Context, itemIDs []string) error
	// DeleteBy
	// DeleteMulti
	// DeleteAll() error

	Iterator(ctx context.Context, q Query) (Iter, error)
	// IteratorBy(key, op string, value interface{}) (Iter, error)

	// // Changelog stuff: move this to it's own file
	// GetChangelogsForObject(id string) ([]Changelog, error)
	// GetLatestChangelogForObject(id string) (*Changelog, error)

	// DeleteChangelogs(ids ...string) error

	// ChangelogIterator() (ChangelogIter, error)

	// Audit() (map[string]*Changelog, error)
	// QuickSync() error

	// Sync() error
}

type Result struct {
	Item Item
	Err  error
}
