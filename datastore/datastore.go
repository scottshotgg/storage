package datastore

import (
	"context"
	"errors"

	dstore "cloud.google.com/go/datastore"
	"github.com/pizzahutdigital/datastore"
	pb "github.com/pizzahutdigital/storage/protobufs"
	"github.com/pizzahutdigital/storage/store"
)

type DB struct {
	Instance *datastore.DSInstance
}

func (db *DB) Get(id string) (store.Item, error) {
	var s pb.Item

	err := db.Instance.GetDocument(context.Background(), "something", id, &s)
	if err != nil {
		return nil, err
	}

	return store.NewObject(s.GetId(), s.GetValue()), nil
}

func (db *DB) Set(id string, i store.Item) error {
	ctx := context.Background()

	cl := store.GenInsertChangelog(i)
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

	cl := store.GenDeleteChangelog(id)
	err := db.Instance.UpsertDocument(ctx, "changelog", cl.ID, cl)
	if err != nil {
		return err
	}

	return db.Instance.DeleteDocument(ctx, "something", id)
}

func (db *DB) Iterator() (store.Iter, error) {
	return &Iter{
		I: db.Instance.Run(context.Background(), db.Instance.NewQuery("something")),
	}, nil
}

func (db *DB) ChangelogIterator() (store.ChangelogIter, error) {
	return &ChangelogIter{
		I: db.Instance.Run(context.Background(), db.Instance.NewQuery("changelog")),
	}, nil
}

type Iter struct {
	I *dstore.Iterator
}

type ChangelogIter struct {
	I *dstore.Iterator
}

func (i *Iter) Next() (store.Item, error) {
	var s pb.Item

	_, err := i.I.Next(&s)
	if err != nil {
		return nil, err
	}

	return store.NewObject(s.GetId(), s.GetValue()), nil
}

func (i *ChangelogIter) Next() (*store.Changelog, error) {
	var cl store.Changelog

	_, err := i.I.Next(&cl)
	if err != nil {
		return nil, err
	}

	return &cl, nil
}

func (db *DB) GetLatestChangelogForObject(id string) (*store.Changelog, error) {
	return nil, errors.New("Not implemented")
}
