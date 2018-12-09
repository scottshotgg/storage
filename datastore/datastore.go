package datastore

import (
	"context"
	"errors"
	"time"

	dstore "cloud.google.com/go/datastore"
	"github.com/pizzahutdigital/datastore"
	"github.com/pizzahutdigital/storage/object"
	pb "github.com/pizzahutdigital/storage/protobufs"
	"github.com/pizzahutdigital/storage/storage"
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

type result struct {
	Item *pb.Item
	Err  error
}

func (db *DB) Get(ctx context.Context, id string) (storage.Item, error) {
	var (
		s       pb.Item
		resChan = make(chan *result)
		res     *result
	)

	go func() {
		select {
		case resChan <- &result{
			Item: &s,
			Err:  db.Instance.GetDocument(ctx, "something", id, &s),
		}:
		}
	}()

	for {
		select {
		case res = <-resChan:
			if res.Err != nil {
				return nil, res.Err
			}

			return object.New(res.Item.GetId(), res.Item.GetValue()), nil

		case <-time.After(GetTimeout):
			return nil, ErrTimeout
		}
	}
}

func (db *DB) Set(id string, i storage.Item) error {
	ctx := context.Background()

	cl := storage.GenInsertChangelog(i)
	err := db.Instance.UpsertDocument(ctx, "changelog", cl.ID, cl)
	if err != nil {
		return err
	}

	return db.Instance.UpsertDocument(ctx, "something", id, &pb.Item{
		Id:    i.ID(),
		Value: i.Value(),
	})
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
