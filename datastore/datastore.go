package datastore

import (
	"context"
	"errors"
	"sync"
	"time"

	dstore "cloud.google.com/go/datastore"
	"github.com/pizzahutdigital/datastore"
	"github.com/pizzahutdigital/storage/object"
	pb "github.com/pizzahutdigital/storage/protobufs"
	"github.com/pizzahutdigital/storage/storage"
	"google.golang.org/api/iterator"
)

type DB struct {
	Instance *datastore.DSInstance
}

type Iter struct {
	I *dstore.Iterator
}

type ChangelogIter struct {
	I *dstore.Iterator
}

const (
	GetTimeout = 1 * time.Second
)

var (
	ErrTimeout = errors.New("Timeout")
)

func (db *DB) Get(ctx context.Context, id string) (storage.Item, error) {
	var (
		props dstore.PropertyList
		err   = db.Instance.GetDocument(ctx, "something", id, &props)
	)

	return object.FromProps(props), err
}

func (db *DB) GetWithTimeout(ctx context.Context, id string, timeout time.Duration) (storage.Item, error) {
	if timeout < 1 {
		return db.Get(ctx, id)
	}

	var (
		o       object.Object
		resChan = make(chan *storage.Result)
		res     *storage.Result
	)

	defer close(resChan)

	go func() {
		select {
		case resChan <- &storage.Result{
			Item: &o,
			Err:  db.Instance.GetDocument(ctx, "something", id, &o),
		}:
		}
	}()

	for {
		select {
		case res = <-resChan:
			if res.Err != nil {
				return nil, res.Err
			}

			return object.FromResult(res), nil

		case <-time.After(timeout):
			return nil, ErrTimeout
		}
	}
}

// Use a builder pattern or `query` to make these
func (db *DB) GetAsync(ctx context.Context, id string, timeout time.Duration) <-chan *storage.Result {
	var resChan = make(chan *storage.Result)

	go func() {
		item, err := db.GetWithTimeout(ctx, id, timeout)

		select {
		case resChan <- &storage.Result{
			Item: item,
			Err:  err,
		}:
		}
		// TODO: do something like this with a custom datastructure
		// attemptChanWrite(resChan, res)
	}()

	return resChan
}

func (db *DB) GetMulti(ctx context.Context, ids ...string) (items []storage.Item, err error) {
	var (
		itemChan = make(chan storage.Item, len(ids))
		wg       sync.WaitGroup
		item     storage.Item
	)

	for i := range ids {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			// TODO: could use GetWithTimeout here or wait till we have our Query interface
			var item, err = db.Get(ctx, ids[i])
			if err != nil {
				// log
				return
			}

			// TODO: how do we alert to the end user that it failed
			itemChan <- item
		}(i)
	}

	wg.Wait()
	close(itemChan)

	for item = range itemChan {
		items = append(items, item)
	}

	return items, nil
}

// TODO: could we just do interface here?
func (db *DB) Set(id string, i storage.Item, sk map[string]interface{}) error {
	ctx := context.Background()

	// cl := storage.GenInsertChangelog(i)
	// err := db.Instance.UpsertDocument(ctx, "changelog", cl.ID, cl)
	// if err != nil {
	// 	return err
	// }

	var (
		key = dstore.Key{
			Kind:      "something",
			Name:      i.ID(),
			Namespace: db.Instance.Namespace(),
		}

		props = dstore.PropertyList{
			dstore.Property{
				Name:  "id",
				Value: i.ID(),
			},
			dstore.Property{
				Name:  "timestamp",
				Value: i.Timestamp(),
			},
			dstore.Property{
				Name:  "value",
				Value: i.Value(),
			},
		}
	)

	for key, value := range sk {
		if value == nil {
			continue
		}

		props = append(props, dstore.Property{
			Name:  key,
			Value: value,
		})
	}

	_, err := db.Instance.Client().Put(ctx, &key, &props)
	return err

	// return db.Instance.UpsertDocument(ctx, "something", id, &pb.Item{
	// 	Id:    i.ID(),
	// 	Value: i.Value(),
	// })
}

// Make the value have a "string" attached to it
func (db *DB) GetBySK(key string, op string, value interface{}, limit int) (items []storage.Item, err error) {
	var (
		ctx   = context.Background()
		query = db.Instance.NewQuery("something").Filter(key+op, value).Limit(limit)
		iter  = db.Instance.Client().Run(ctx, query)
	)

	for {
		var props dstore.PropertyList
		_, err = iter.Next(&props)
		if err != nil {
			if err == iterator.Done {
				break
			}

			return nil, err
		}

		items = append(items, object.FromProps(props))
	}

	return items, nil
}

func (db *DB) Delete(id string) error {
	ctx := context.Background()

	cl := storage.GenDeleteChangelog(id)
	err := db.Instance.UpsertDocument(ctx, "changelog", cl.ID, cl)
	if err != nil {
		return err
	}

	return db.Instance.DeleteDocument(ctx, "something", id)
}

func (db *DB) Iterator() (storage.Iter, error) {
	return &Iter{
		I: db.Instance.Run(context.Background(), db.Instance.NewQuery("something")),
	}, nil
}

func (db *DB) ChangelogIterator() (storage.ChangelogIter, error) {
	return &ChangelogIter{
		I: db.Instance.Run(context.Background(), db.Instance.NewQuery("changelog")),
	}, nil
}

func (i *Iter) Next() (storage.Item, error) {
	var s pb.Item

	_, err := i.I.Next(&s)
	if err != nil {
		return nil, err
	}

	return object.New(s.GetId(), s.GetValue()), nil
}

func (i *ChangelogIter) Next() (*storage.Changelog, error) {
	var cl storage.Changelog

	_, err := i.I.Next(&cl)
	if err != nil {
		return nil, err
	}

	return &cl, nil
}

func (db *DB) GetLatestChangelogForObject(id string) (*storage.Changelog, error) {
	return nil, errors.New("Not implemented")
}
