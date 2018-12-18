package grpc

import (
	dberrors "github.com/scottshotgg/storage/errors"
	pb "github.com/scottshotgg/storage/protobufs"
	"github.com/scottshotgg/storage/storage"
)

// TODO: Need to make this more generic so same logic can be used

// Iter does a paging style
type Iter struct {
	I pb.Storage_IteratorClient
}

// Iter does a paging style
type IterBy struct {
	I pb.Storage_IteratorByClient
}

func next(i *Iter) (storage.Item, error) {
	/*
		This will have to call Recv from the iterator client; look at ShitStreamer
	*/

	return nil, dberrors.ErrNotImplemented
}

func (i *Iter) Next() (storage.Item, error) {
	return next()
}

func (i *IterBy) Next() (storage.Item, error) {
	/*
		This will have to call Recv from the iterator client; look at ShitStreamer
	*/

	return nil, dberrors.ErrNotImplemented
}
