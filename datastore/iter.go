package datastore

import (
	dstore "cloud.google.com/go/datastore"
	"github.com/pizzahutdigital/storage/object"
	"github.com/pizzahutdigital/storage/storage"
)

type Iter struct {
	I *dstore.Iterator
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
