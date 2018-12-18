package http

import (
	"context"
	"net/http"
	"net/url"
	"time"

	dberrors "github.com/scottshotgg/storage/errors"
	"github.com/scottshotgg/storage/storage"
)

// This package will perform based on it's versioned Swagger spec

// WEB SOCKETS FOR ALL
// * Oprah gif *
// https://github.com/tmc/grpc-websocket-proxy

const (
	httpTimeout = 60 * time.Second
)

// DB implements Storage from the storage package
type DB struct {
	client  *http.Client
	address string
	addr    *url.URL
}

/*
	TODO: look at generating this as a Swagger client from the Swagger that was generated from the protobufs... lol
*/

type AuthFunc func() error

// Close should call close internally on your implementations DB and clean up any lose ends; channels, waitgroups, etc
func (db *DB) Close() error {
	return dberrors.ErrNotImplemented
}

func New(addr *url.URL, auth AuthFunc) (*DB, error) {
	// Check the URL
	if addr == nil {
		return nil, dberrors.ErrCouldNotOpenDB
	}

	var err error
	// Authenticate if they have set that
	if auth != nil {
		err = auth()
		if err != nil {
			return nil, err
		}
	}

	// Convert their url to a string
	var (
		address = addr.String()
		myAddr  *url.URL
	)

	// Parse the url for two reasons:
	// 1) Ensure that it is valid
	// 2) We now own the url variable
	myAddr, err = url.Parse(address)
	if err != nil {
		return nil, err
	}

	var client = &http.Client{
		Timeout: httpTimeout,
	}

	// Get the metadata from the client or something to check that it is up

	return &DB{
		client:  client,
		address: address,
		addr:    myAddr,
	}, nil
}

// All of these functions return ErrNotImplemented for now so that the interface can be satisfied

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

func (db *DB) Iterator() (storage.Iter, error) {
	return nil, dberrors.ErrNotImplemented
}

func (db *DB) IteratorBy(key, op string, value interface{}) (storage.Iter, error) {
	return nil, dberrors.ErrNotImplemented
}

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

// storage.Changelog stuff: move this to it's own file

// DeleteBy
// DeleteMulti
// DeleteAll() error
