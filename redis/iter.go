package redis

import (
	"errors"

	redigo "github.com/go-redis/redis"
	"github.com/scottshotgg/storage/object"
	"github.com/scottshotgg/storage/storage"
	"google.golang.org/api/iterator"
)

type Iter struct {
	I *redigo.ScanIterator
}

func (i *Iter) Next() (storage.Item, error) {
	ok := i.I.Next()
	if !ok {
		return nil, iterator.Done
	}

	err := i.I.Err()
	if err != nil {
		return nil, err
	}

	ok = i.I.Next()
	if !ok {
		return nil, errors.New("Did not have value for key")
	}

	err = i.I.Err()
	if err != nil {
		return nil, err
	}

	var (
		o object.Object
	)

	return &o, o.UnmarshalBinary([]byte(i.I.Val()))
}
