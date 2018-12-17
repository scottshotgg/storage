package inmem

import (
	"errors"
	"sync"

	"github.com/scottshotgg/storage/storage"
)

type DB struct {
	Instance sync.Map
}

func (db *DB) Get(id string) (storage.Item, error) {
	value, _ := db.Instance.Load(id)

	itemValue, ok := value.(storage.Item)
	if !ok {
		return nil, errors.New("WTF")
	}

	return itemValue, nil
}

func (db *DB) Set(id string, i storage.Item) error {
	db.Instance.Store(id, i)

	return nil
}

func (db *DB) Iterator() (storage.Iter, error) {
	return nil, errors.New("Not implemented")
}

type Iter struct {
	// I *dstorage.Iterator
}

// func (i *Iter) Next() (storage.Item, error) {
// 	var s pb.Item

// 	_, err := i.I.Next(&s)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return store.NewObject(s.GetId(), s.GetValue()), nil
// }
