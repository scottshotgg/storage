package redis

import (
	"context"
	"errors"
	"fmt"

	redigo "github.com/go-redis/redis"
	"github.com/pizzahutdigital/storage/object"
	"github.com/pizzahutdigital/storage/storage"
	"google.golang.org/api/iterator"
)

type DB struct {
	Instance *redigo.Client
}

func (db *DB) Get(_ context.Context, id string) (storage.Item, error) {
	value, err := db.Instance.HGet("something", id).Result()
	if err != nil {
		return nil, err
	}

	var o = &object.Object{}

	return o, o.UnmarshalBinary([]byte(value))
}

func (db *DB) Set(id string, i storage.Item) error {
	_, err := db.Instance.HSet("something", id, i).Result()
	return err
}

func (db *DB) Delete(id string) error {
	return db.Instance.HDel("something", id).Err()
}

func (db *DB) Iterator() (storage.Iter, error) {
	return &Iter{
		I: db.Instance.HScan("something", 0, "", 1000000).Iterator(),
	}, nil
}

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

func (db *DB) ChangelogIterator() (storage.ChangelogIter, error) {
	return nil, errors.New("Not implemented")
}

func (db *DB) GetLatestChangelogForObject(id string) (*storage.Changelog, error) {
	fmt.Println("doin it", id)

	keys, _, err := db.Instance.HScan("changelog", 0, id+"*", 1000000).Result()
	if err != nil {
		return nil, err
	}

	cls, err := db.Instance.HMGet("changelog", keys...).Result()
	if err != nil {
		return nil, err
	}

	var (
		clValues = make([]storage.Changelog, len(cls))
		clAssert storage.Changelog
		ok       bool
	)

	for i, cl := range cls {
		clAssert, ok = cl.(storage.Changelog)
		if !ok {
			return nil, errors.New("wtf")
		}

		clValues[i] = clAssert
	}

	return findLatestTS(clValues), nil
}

func findLatestTS(clValues []storage.Changelog) *storage.Changelog {
	var latest = &storage.Changelog{}

	for _, clValue := range clValues {
		if clValue.Timestamp > latest.Timestamp {
			latest = &clValue
		}

		// might need to do something about ones that are the same
	}

	return latest
}

// TODO: implement changelog stuff for redis
