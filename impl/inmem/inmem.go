package inmem

import (
	"errors"
	"sync"

	"github.com/scottshotgg/storage/storage"
)

type NonSyncStore struct {
	sync.RWMutex
	Data map[string]storage.Item
}

type DB struct {
	Instance *NonSyncStore
}

type Iter struct {
	// I *dstorage.Iterator
}

func (db *DB) Get(id string) (storage.Item, error) {
	db.Instance.RLock()
	defer db.Instance.RUnlock()

	return db.Instance.Data[id], nil
}

func (db *DB) Set(id string, i storage.Item) error {
	db.Instance.Lock()

	db.Instance.Data[id] = i

	db.Instance.Unlock()

	return nil
}

func (db *DB) Iterator() (storage.Iter, error) {
	return nil, errors.New("Not implemented")
}

func (db *DB) ChangelogIterator() (storage.ChangelogIter, error) {
	return nil, errors.New("Not implemented")
}

func (db *DB) Delete(id string) error {
	return errors.New("Not implemented")
}

func (db *DB) GetLatestChangelogForObject(id string) (*storage.Changelog, error) {
	return nil, errors.New("Not implemented")
}

// func (i *Iter) Next() (storage.Item, error) {
// 	var s pb.Item

// 	_, err := i.I.Next(&s)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return store.NewObject(s.GetId(), s.GetValue()), nil
// }
