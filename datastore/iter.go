package datastore

import (
	dstore "cloud.google.com/go/datastore"
	"github.com/pizzahutdigital/storage/object"
	pb "github.com/pizzahutdigital/storage/protobufs"
	"github.com/pizzahutdigital/storage/storage"
)

type Iter struct {
	I *dstore.Iterator
}

func (i *Iter) Next() (storage.Item, error) {
	var s pb.Item

	_, err := i.I.Next(&s)
	if err != nil {
		return nil, err
	}

	return object.New(s.GetId(), s.GetValue()), nil
}
