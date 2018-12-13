package datastore

import (
	"context"
	"errors"

	dstore "cloud.google.com/go/datastore"
	"github.com/pizzahutdigital/storage/object"
	"github.com/pizzahutdigital/storage/storage"
)

type Iter struct {
	I *dstore.Iterator
}

func (db *DB) Iterator() (storage.Iter, error) {
	return &Iter{
		I: db.Instance.Run(context.Background(), db.Instance.NewQuery("something")),
	}, nil
}

func (db *DB) IteratorBy(key, op string, value interface{}) (storage.Iter, error) {
	var query = db.Instance.NewQuery("something")

	if len(key) != 0 {
		if len(op) == 0 {
			return nil, errors.New("Must provide an operator")
		}

		query = query.Filter(key+op, value)
	}

	// // TODO: might need to do this
	// if value == nil {
	// 	return nil, errors.New("Must provide an operator")
	// }

	return &Iter{
		I: db.Instance.Run(context.Background(), query),
	}, nil
}

func (i *Iter) Next() (storage.Item, error) {
	var (
		props  dstore.PropertyList
		_, err = i.I.Next(&props)
	)

	if err != nil {
		return nil, err
	}

	return object.FromProps(props), nil
}
