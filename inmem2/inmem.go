package inmem

import (
	"errors"
	"sync"

	"github.com/pizzahutdigital/storage/store"
)

type DB struct {
	Instance sync.Map
}

func (db *DB) Get(id string) (store.Item, error) {
	value, _ := db.Instance.Load(id)

	itemValue, ok := value.(store.Item)
	if !ok {
		return nil, errors.New("WTF")
	}

	return itemValue, nil
}

func (db *DB) Set(id string, i store.Item) error {
	db.Instance.Store(id, i)

	return nil
}

func (db *DB) Iterator() (store.Iter, error) {
	return nil, errors.New("Not implemented")
}

type Iter struct {
	// I *dstore.Iterator
}

// func (i *Iter) Next() (store.Item, error) {
// 	var s pb.Item

// 	_, err := i.I.Next(&s)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return store.NewObject(s.GetId(), s.GetValue()), nil
// }
