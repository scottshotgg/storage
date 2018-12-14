package inmem

import (
	"errors"
	"sync"

	"github.com/scottshotgg/storage/store"
)

type NonSyncStore struct {
	sync.RWMutex
	Data map[string]store.Item
}

type DB struct {
	Instance *NonSyncStore
}

type Iter struct {
	// I *dstore.Iterator
}

func (db *DB) Get(id string) (store.Item, error) {
	db.Instance.RLock()
	defer db.Instance.RUnlock()

	return db.Instance.Data[id], nil
}

func (db *DB) Set(id string, i store.Item) error {
	db.Instance.Lock()

	db.Instance.Data[id] = i

	db.Instance.Unlock()

	return nil
}

func (db *DB) Iterator() (store.Iter, error) {
	return nil, errors.New("Not implemented")
}

func (db *DB) ChangelogIterator() (store.ChangelogIter, error) {
	return nil, errors.New("Not implemented")
}

func (db *DB) Delete(id string) error {
	return errors.New("Not implemented")
}

func (db *DB) GetLatestChangelogForObject(id string) (*store.Changelog, error) {
	return nil, errors.New("Not implemented")
}

// func (i *Iter) Next() (store.Item, error) {
// 	var s pb.Item

// 	_, err := i.I.Next(&s)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return store.NewObject(s.GetId(), s.GetValue()), nil
// }
