package sql

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	dberrors "github.com/scottshotgg/storage/errors"
	"github.com/scottshotgg/storage/storage"
)

// DB implements Storage from the storage package
type DB struct {
	Instance *sql.DB
}

func (db *DB) Close() error {
	return db.Instance.Close()
}

func New(name, connString string) (*DB, error) {
	var db, err = sql.Open(name, connString)
	if err != nil {
		return nil, err
	}

	return &DB{
		Instance: db,
	}, nil
}

func NewFrom(db *sql.DB) *DB {
	return &DB{
		Instance: db,
	}
}

func (db *DB) Get(ctx context.Context, id string) (storage.Item, error) {
	return nil, dberrors.ErrNotImplemented
}

func (db *DB) GetBy(ctx context.Context, id, op string, value interface{}, limit int) ([]storage.Item, error) {
	return nil, dberrors.ErrNotImplemented
}

func (db *DB) GetMulti(ctx context.Context, ids ...string) ([]storage.Item, error) {
	return nil, dberrors.ErrNotImplemented
}

func (db *DB) GetAll(ctx context.Context) ([]storage.Item, error) {
	return nil, dberrors.ErrNotImplemented
}

func (db *DB) Set(ctx context.Context, item storage.Item) error {
	return dberrors.ErrNotImplemented
}

func (db *DB) SetMulti(ctx context.Context, items []storage.Item) error {
	return dberrors.ErrNotImplemented
}

func (db *DB) Delete(id string) error {
	return dberrors.ErrNotImplemented
}

// DeleteBy
// DeleteMulti
// DeleteAll() error

func (db *DB) Iterator() (storage.Iter, error) {
	return nil, dberrors.ErrNotImplemented
}

func (db *DB) IteratorBy(key, op string, value interface{}) (storage.Iter, error) {
	return nil, dberrors.ErrNotImplemented
}

// storage.Changelog stuff: move this to it's own file

func (db *DB) GetChangelogsForObject(id string) ([]storage.Changelog, error) {
	return nil, dberrors.ErrNotImplemented
}

func (db *DB) GetLatestChangelogForObject(id string) (*storage.Changelog, error) {
	return nil, dberrors.ErrNotImplemented
}

func (db *DB) DeleteChangelogs(ids ...string) error {
	return dberrors.ErrNotImplemented
}

func (db *DB) ChangelogIterator() (storage.ChangelogIter, error) {
	return nil, dberrors.ErrNotImplemented
}
