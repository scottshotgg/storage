package datastore

import (
	"context"

	dstore "cloud.google.com/go/datastore"
	"github.com/scottshotgg/storage/storage"
)

type ChangelogIter struct {
	I *dstore.Iterator
}

func (db *DB) ChangelogIterator() (storage.ChangelogIter, error) {
	return &ChangelogIter{
		I: db.Instance.Run(context.Background(), db.Instance.NewQuery("changelog")),
	}, nil
}

func (i *ChangelogIter) Next() (*storage.Changelog, error) {
	var cl storage.Changelog

	_, err := i.I.Next(&cl)
	if err != nil {
		return nil, err
	}

	return &cl, nil
}

func getLatest(cls []storage.Changelog) (*storage.Changelog, error) {
	var latest storage.Changelog

	for _, cl := range cls {
		if cl.Timestamp > latest.Timestamp {
			latest = cl
		}
	}

	return &latest, nil
}

func (db *DB) GetLatestChangelogForObject(id string) (*storage.Changelog, error) {
	return nil, ErrNotImplemented
}

func (db *DB) GetChangelogsForObject(id string) ([]storage.Changelog, error) {
	var (
		ctx   = context.Background()
		query = db.Instance.NewQuery("changelog").Filter("ObjectID=", id)
		// iter  = db.Instance.Client().Run(ctx, query)
		cls []storage.Changelog
	)

	err := db.Instance.GetDocuments(ctx, query, &cls)
	if err != nil {
		return nil, err
	}

	// for {
	// 	var s dstore.PropertyList
	// 	// var s pb.Item
	// 	_, err = iter.Next(&s)
	// 	if err != nil {
	// 		if err == iterator.Done {
	// 			break
	// 		}

	// 		return nil, err
	// 	}

	// 	// items = append(items, object.FromProto(&s))
	// 	items = append(items, object.FromProps(s))
	// }

	return cls, err
}

func (db *DB) DeleteChangelogs(ids ...string) error {
	return db.Instance.DeleteDocuments(context.Background(), "changelog", ids)
}
