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

// DB implements Storage from the storage package
type DB struct {
	Instance *redigo.Client
}

type ChangelogIter struct {
	I *redigo.ScanIterator
}

var (
	ErrTimeout        = errors.New("Timeout")
	ErrNotImplemented = errors.New("Not implemented")
)

func (db *DB) Get(ctx context.Context, id string) (storage.Item, error) {
	var (
		value string
		err   error
	)

	// something := make(chan struct{})
	// go func() {
	// 	value, err = db.Instance.HGet("something", id).Result()
	// 	if err != nil {
	// 		// log
	// 	}

	// 	fmt.Println("i am here")
	// 	something <- struct{}{}
	// }()

	// select {
	// case <-ctx.Done():
	// 	return nil, context.Canceled

	// case <-something:
	// }

	value, err = db.Instance.HGet("something", id).Result()
	if err != nil {
		// log
		return nil, err
	}

	// select {
	// case <-ctx.Done():
	// 	return nil, context.Canceled

	// default:
	// }

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

	changelog := storage.GenInsertChangelog(i)
	_, err = db.Instance.HSet("changelog", changelog.ID, changelog).Result()
	if err != nil {
		return err
	}

	i.(*object.Object).SetTimestamp(changelog.Timestamp)
	_, err = db.Instance.HSet("something", id, i).Result()
	return err
}

func (db *DB) GetMulti(_ context.Context, ids ...string) ([]storage.Item, error) {
	if len(ids) == 0 {
		return []storage.Item{}, nil
	}

	var values, err = db.Instance.HMGet("something", ids...).Result()
	if err != nil {
		return nil, err
	}

	var items = make([]storage.Item, len(values))
	for i := range values {
		var o object.Object
		err = o.UnmarshalBinary([]byte(values[i].(string)))
		if err != nil {
			// log
			continue
		}

		items[i] = &o
	}

	return items, nil
}

func (db *DB) GetBy(key string, op string, value interface{}, limit int) ([]storage.Item, error) {
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
	for i := range keys {
		ids = append(ids, extractIDFromKey(keys[i]))
	}

	return db.GetMulti(nil, ids...)
}

func (db *DB) Delete(id string) error {
	return db.Instance.HDel("something", id).Err()
}

func (db *DB) DeleteChangelogs(ids ...string) error {
	// return db.Instance.HDel("changelog", id).Err()
	return db.Instance.HDel("changelog", ids...).Err()
}

func (db *DB) Iterator() (storage.Iter, error) {
	return &Iter{
		I: db.Instance.HScan("something", 0, "", 1000000).Iterator(),
	}, nil
}

func (db *DB) ChangelogIterator() (storage.ChangelogIter, error) {
	return &ChangelogIter{
		I: db.Instance.HScan("changelog", 0, "", 1000000).Iterator(),
	}, nil
}

func (i *ChangelogIter) Next() (*storage.Changelog, error) {
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
		cl storage.Changelog
	)

	return &cl, cl.UnmarshalBinary([]byte(i.I.Val()))
}

func (db *DB) GetLatestChangelogForObject(id string) (*storage.Changelog, error) {
	keys, _, err := db.Instance.HScan("changelog", 0, id+"-*", 1000000).Result()
	if err != nil {
		return nil, err
	}

	// var oids []string
	// for i := range keys {
	// 	split := strings.Split(keys[i], "-")

	//
	// }

	var changelogs []storage.Changelog

	for i := range keys {
		if i%2 == 0 {
			continue
		}

		var cl = &storage.Changelog{}
		cl.UnmarshalBinary([]byte(keys[i]))

		changelogs = append(changelogs, *cl)
	}

	return findLatestTS(changelogs), nil

	// cls, err := db.Instance.HMGet("changelog", keys...).Result()
	// if err != nil {
	// 	return nil, err
	// }

	// fmt.Println("cls", cls)

	// for _, cl := range cls {
	// 	fmt.Println("cl2", cl)
	// 	if cl == nil {
	// 		fmt.Println("nil")
	// 	}
	// }

	// fmt.Println("cls", cls)

	// var (
	// 	clValues = make([]storage.Changelog, len(cls))
	// 	// ok       bool
	// )

	// for i, cl := range cls {
	// 	// // cla, ok := cl.(string)
	// 	// cla := cl.(string)
	// 	// if !ok {
	// 	// 	return nil, errors.New("Could not assert value")
	// 	// }

	// 	// err = clAssert.UnmarshalBinary([]byte(cla))
	// 	// if err != nil {
	// 	// 	return nil, err
	// 	// }

	// 	// clValues[i] = clAssert

	// 	fmt.Println("cl", cl, i)
	// }

	// return findLatestTS(clValues), nil
	// return nil, nil
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

func (db *DB) GetChangelogsForObject(id string) ([]storage.Changelog, error) {
	keys, _, err := db.Instance.HScan("changelog", 0, id+"-*", 1000000).Result()
	if err != nil {
		return nil, err
	}

	var changelogs []storage.Changelog

	for i := range keys {
		if i%2 == 0 {
			continue
		}

		var cl = &storage.Changelog{}
		cl.UnmarshalBinary([]byte(keys[i]))

		changelogs = append(changelogs, *cl)
	}

	return changelogs, nil
}
