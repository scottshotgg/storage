package redis

import (
	"context"
	"errors"
	"fmt"
	"strings"

	redigo "github.com/go-redis/redis"
	"github.com/pizzahutdigital/storage/object"
	"github.com/pizzahutdigital/storage/storage"
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

func (db *DB) Set(id string, i storage.Item, sk map[string]interface{}) (err error) {
	for key, value := range sk {
		fmt.Println("sk", key, value)
		_, err = db.Instance.SAdd("::something::"+key, id+"::"+fmt.Sprintf("%v", value)).Result()
		// _, err = db.Instance.HSet("::something::"+key, id+"::"+fmt.Sprintf("%v", value), value).Result()
		if err != nil {
			fmt.Println("err", err)
			return err
		}
	}

	_, err = db.Instance.HSet("something", id, i).Result()
	return err
}

func (db *DB) GetMulti(_ context.Context, ids ...string) (items []storage.Item, err error) {
	if len(ids) == 0 {
		return items, nil
	}

	values, err := db.Instance.HMGet("something", ids...).Result()
	if err != nil {
		return nil, err
	}

	for _, value := range values {
		var o = &object.Object{}
		err = o.UnmarshalBinary([]byte(value.(string)))
		if err != nil {
			// log
		}

		items = append(items, o)
	}

	return items, nil
}

func (db *DB) GetBySK(key string, op string, value interface{}, limit int) ([]storage.Item, error) {
	amount, err := db.Instance.SCard("::something::" + key).Result()
	if err != nil {
		return nil, err
	}

	// TODO: this might need to return an error instead... need to look at how that will define the usage
	if amount == 0 {
		return []storage.Item{}, nil
	}

	if limit == -1 {
		limit = int(amount)
	}

	// ids, _, err := db.Instance.HScan("::something::"+key, 0, fmt.Sprintf("%v", value), int64(limit)).Result()
	keys, _, err := db.Instance.SScan("::something::"+key, 0, "*::"+fmt.Sprintf("%v", value), int64(limit)).Result()
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, key := range keys {
		ids = append(ids, extractIDFromKey(key))
	}

	return db.GetMulti(nil, ids...)
}

func extractIDFromKey(key string) string {
	split := strings.Split(key, "::")
	if len(split) > 1 {
		return split[0]
	}

	return ""
}

func (db *DB) Delete(id string) error {
	return db.Instance.HDel("something", id).Err()
}

func (db *DB) Iterator() (storage.Iter, error) {
	return &Iter{
		I: db.Instance.HScan("something", 0, "", 1000000).Iterator(),
	}, nil
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
