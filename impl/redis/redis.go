package redis

import (
	"context"
	"errors"
	"fmt"

	redigo "github.com/go-redis/redis"
	"google.golang.org/api/iterator"

	"github.com/google/uuid"
	"github.com/scottshotgg/storage/audit"
	dberrors "github.com/scottshotgg/storage/errors"
	"github.com/scottshotgg/storage/object"
	"github.com/scottshotgg/storage/storage"
)

// DB implements Storage from the storage package
type DB struct {
	id       string
	Instance *redigo.Client
}

type ChangelogIter struct {
	I  *redigo.ScanIterator
	DB *DB
}

func validUUID(id string) bool {
	var _, err = uuid.Parse(id)
	if err != nil {
		// log
		return false
	}

	return true
}

func createMetadata() error {
	return nil
}

func initDB(client *redigo.Client) (*DB, error) {
	/*
		Things that should be stores in meta:
		- id
		- creation date
		- last connected
		- keys
		- maybe db structure
		- maybe count if we can pull that off
	*/

	var meta, err = client.HGetAll("__meta__").Result()
	if err != nil {
		return nil, err
	}

	// // TODO:
	// if len(meta) == 0 {
	// 	createMetadata()
	// }

	var id = meta["id"]

	// Parse the UUID to validate that the ID in the database is correct
	_, err = uuid.Parse(id)

	// If we got an error validating the last ID then regenerate one
	if err != nil {
		// log
		// TODO: somehow the DB does not have an ID
		// create one

		v4, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}

		ok, err := client.HSet("__meta__", "id", v4.String()).Result()
		if err != nil || !ok {
			// return nil, dberrors.ErrMetaCouldNotBeCreated
			return nil, errors.New("ErrMetaCouldNotBeCreated")
		}

		id = v4.String()
	}

	return &DB{
		id:       id,
		Instance: client,
	}, nil
}

// Might consider moving these function to another package or into the main package

// New creates a new Redis DB and checks it's connection
func New(opts *redigo.Options) (*DB, error) {
	var (
		client = redigo.NewClient(opts)

		// Check that Redis is connected
		err = client.Ping().Err()
	)

	if err != nil {
		// log
		return nil, err
	}

	return initDB(client)
}

func (db *DB) ID() string {
	return db.id
}

// Audit uses the default audit function
func (db *DB) Audit() (map[string]*storage.Changelog, error) {
	return audit.Audit(db)
}

// QuickSync doesn't make sense in a single DB context but is here to satisfy the interface
func (db *DB) QuickSync() error {
	return nil
}

// Sync doesn't make sense in a single DB context but is here to satisfy the interface
func (db *DB) Sync() error {
	return nil
}

func (db *DB) Get(_ context.Context, id string) (storage.Item, error) {
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

func (db *DB) GetAll(_ context.Context) ([]storage.Item, error) {
	var results, err = db.Instance.HGetAll("something").Result()
	if err != nil {
		return nil, err
	}

	var (
		i     int
		items = make([]storage.Item, len(results))
	)

	for _, result := range results {
		items[i] = &object.Object{}
		err = items[i].UnmarshalBinary([]byte(result))
		if err != nil {
			return items, err
		}

		i++
	}

	return items, nil
}

func buildZIndexSecondaryKey(i storage.Item) (m redigo.Z) {
	var sk = "id::" + i.GetID() + "::"

	for key, value := range i.GetKeys() {
		sk += fmt.Sprintf("%s::%v::", key, value)
	}

	return redigo.Z{
		Score:  0,
		Member: sk[:len(sk)-2],
	}
}

func (db *DB) Set(_ context.Context, i storage.Item) (err error) {
	for key, value := range i.GetKeys() {
		_, err = db.Instance.SAdd("::something::"+key, i.GetID()+"::"+fmt.Sprintf("%v", value)).Result()

		// _, err = db.Instance.ZAdd("::something_sk::", buildZIndexSecondaryKey(i)).Result()
		if err != nil {
			return err
		}
	}

	var changelog = storage.GenInsertChangelog(i)
	_, err = db.Instance.HSet("changelog", changelog.ID, changelog).Result()
	if err != nil {
		return err
	}

	i.(*object.Object).SetTimestamp(changelog.Timestamp)
	_, err = db.Instance.HSet("something", i.GetID(), i).Result()
	if err != nil {
		// delete the changelog
		return err
	}

	return nil
}

func (db *DB) SetMulti(_ context.Context, items []storage.Item) error {
	var (
		// Build two maps for the changelog and the items
		clMap   = map[string]interface{}{}
		itemMap = map[string]interface{}{}
		// sks     []redigo.Z
		// sksRemove []redigo.Z
		pipe = db.Instance.Pipeline()
	)

	for i := range items {
		for key, value := range items[i].GetKeys() {
			var _, err = pipe.SAdd("::something::"+key, items[i].GetID()+"::"+fmt.Sprintf("%v", value)).Result()

			// _, err = db.Instance.ZAdd("::something_sk::", buildZIndexSecondaryKey(i)).Result()
			if err != nil {
				return err
			}
		}

		// Only insert unique items
		if itemMap[items[i].GetID()] == nil {
			itemMap[items[i].GetID()] = items[i]

			var cl = storage.GenInsertChangelog(items[i])
			clMap[cl.ID] = cl
			// sks = append(sks, buildZIndexSecondaryKey(items[i]))
		}
	}

	// Send off the changelogs
	var err = pipe.HMSet("changelog", clMap).Err()
	if err != nil {
		// delete all the changelogs
		return err
	}

	// Send off the items
	err = pipe.HMSet("something", itemMap).Err()
	if err != nil {
		// do something about this
		return err
	}

	// db.Instance.ZRem()

	// _, err = db.Instance.ZAdd("::something_sk::", sks...).Result()
	// if err != nil {
	// 	// do something about this
	// 	return err
	// }

	_, err = pipe.Exec()

	// might have to check the errors from each command

	return err
}

func (db *DB) GetMulti(_ context.Context, ids []string) ([]storage.Item, error) {
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

func (db *DB) GetBy(_ context.Context, key string, op string, value interface{}, limit int) ([]storage.Item, error) {
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

	return db.GetMulti(nil, ids)
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

func (db *DB) IteratorBy(key, op string, value interface{}) (storage.Iter, error) {
	return nil, dberrors.ErrNotImplemented
}

func (db *DB) ChangelogIterator() (storage.ChangelogIter, error) {
	return &ChangelogIter{
		DB: db,
		I:  db.Instance.HScan("changelog", 0, "", 1000000).Iterator(),
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

	var cl = storage.Changelog{
		DBID: i.DB.ID(),
	}

	return &cl, cl.UnmarshalBinary([]byte(i.I.Val()))
}

func (db *DB) GetLatestChangelogForObject(id string) (*storage.Changelog, error) {
	var keys, _, err = db.Instance.HScan("changelog", 0, id+"-*", 1000000).Result()
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

		var cl = &storage.Changelog{
			DBID: db.ID(),
		}

		err = cl.UnmarshalBinary([]byte(keys[i]))
		if err != nil {
			// log
			continue
		}

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
	var keys, _, err = db.Instance.HScan("changelog", 0, id+"-*", 1000000).Result()
	if err != nil {
		return nil, err
	}

	var changelogs []storage.Changelog

	for i := range keys {
		if i%2 == 0 {
			continue
		}

		var cl = &storage.Changelog{
			DBID: db.ID(),
		}

		err = cl.UnmarshalBinary([]byte(keys[i]))
		if err != nil {
			// log
			continue
		}

		changelogs = append(changelogs, *cl)
	}

	return changelogs, nil
}
