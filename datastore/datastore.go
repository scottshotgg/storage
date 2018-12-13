package datastore

import (
	"context"
	"errors"
	"sync"
	"time"

	dstore "cloud.google.com/go/datastore"
	"github.com/hashicorp/go-multierror"
	"github.com/pizzahutdigital/datastore"
	"github.com/pizzahutdigital/storage/object"
	"github.com/pizzahutdigital/storage/storage"
	"google.golang.org/api/iterator"
)

// DB implements Storage from the storage package
type DB struct {
	Instance *datastore.DSInstance
}

const (
	GetTimeout = 1 * time.Second
)

var (
	ErrTimeout        = errors.New("Timeout")
	ErrNotImplemented = errors.New("Not implemented")
)

func (db *DB) Get(ctx context.Context, id string) (storage.Item, error) {
	var (
		props dstore.PropertyList
		err   = db.Instance.GetDocument(ctx, "something", id, &props)
	)

	// TODO: might need to do this
	// if err != nil {
	// 	return nil, err
	// }

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

func (db *DB) GetBy(ctx context.Context, key, op string, value interface{}, limit int) (items []storage.Item, err error) {
	var (
		query = db.Instance.NewQuery("something").Filter(key+op, value).Limit(limit)
		iter  = db.Instance.Client().Run(ctx, query)
	)

	for {
		var s dstore.PropertyList
		// var s pb.Item
		_, err = iter.Next(&s)
		if err != nil {
			if err == iterator.Done {
				break
			}

			return nil, err
		}

		// items = append(items, object.FromProto(&s))
		items = append(items, object.FromProps(s))
	}

	return items, nil
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

func (db *DB) GetAll(ctx context.Context) (items []storage.Item, err error) {
	var iter = db.Instance.Client().Run(ctx, db.Instance.NewQuery("something"))

	for {
		var s dstore.PropertyList
		// var s pb.Item
		_, err = iter.Next(&s)
		if err != nil {
			if err == iterator.Done {
				break
			}

			return items, err
		}

		// items = append(items, object.FromProto(&s))
		items = append(items, object.FromProps(s))
	}

	return items, nil
}

// TODO: could we just do interface here?
func (db *DB) Set(ctx context.Context, i storage.Item) error {
	var (
		cl  = storage.GenInsertChangelog(i)
		err = db.Instance.UpsertDocument(ctx, "changelog", cl.ID, cl)
	)

	// Could make a generic interface and then reflect over the keys at runtime and use
	// those as indexes

	if err != nil {
		return err
	}

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

	for k, v := range i.Keys() {
		if v == nil {
			continue
		}

		props = append(props, dstore.Property{
			Name:  k,
			Value: v,
		})
	}

	_, err = db.Instance.Client().Put(ctx, &key, &props)
	return err

	// return db.Instance.UpsertDocument(ctx, "something", id, &pb.Item{
	// 	Id:    i.ID(),
	// 	Value: i.Value(),
	// })
}

// func buildArray() {}

func drainErrs(errChan chan error) (merr *multierror.Error) {
	close(errChan)

	for err := range errChan {
		merr = multierror.Append(merr, err)
	}

	return merr
}

func (db *DB) SetMulti(ctx context.Context, items []storage.Item) error {
	const (
		amount   = 12
		amountM1 = amount - 1
	)

	var (
		wg      sync.WaitGroup
		errChan = make(chan error, len(items))

		clKeys  = make([]*dstore.Key, len(items))
		clProps = make([]dstore.PropertyList, len(items))

		keys  = make([]*dstore.Key, len(items))
		props = make([]dstore.PropertyList, len(items))

		// The amount of workers here made no difference past about 2.5% of the total length
		workerChan = make(chan struct{}, len(items)/50)
	)

	defer close(workerChan)

	for i := 0; i < len(items)/amount-1; i++ {
		wg.Add(1)
		workerChan <- struct{}{}

		select {
		case <-ctx.Done():
			return context.Canceled

		default:
		}

		go func(start, end int) {
			defer func() {
				wg.Done()
				<-workerChan
			}()

			for j := start; j < end; j++ {
				select {
				case <-ctx.Done():
					return

				default:
				}

				var (
					item = items[j]
					cl   = storage.GenInsertChangelog(item)
				)

				clKeys[j] = &dstore.Key{
					Kind: "changelog",
					Name: cl.ID,
					// Namespace: db.Instance.Namespace(),
				}

				clProps[j] = dstore.PropertyList{
					dstore.Property{
						Name:  "ID",
						Value: cl.ID,
					},
					dstore.Property{
						Name:  "ObjectID",
						Value: cl.ObjectID,
					},
					dstore.Property{
						Name:  "Timestamp",
						Value: cl.Timestamp,
					},
					dstore.Property{
						Name:  "Type",
						Value: cl.Type,
					},
				}

				keys[j] = &dstore.Key{
					Kind: "something",
					Name: item.ID(),
					// Namespace: db.Instance.Namespace(),
				}

				props[j] = dstore.PropertyList{
					dstore.Property{
						Name:  "id",
						Value: item.ID(),
					},
					dstore.Property{
						Name:  "timestamp",
						Value: item.Timestamp(),
					},
					dstore.Property{
						Name:  "value",
						Value: item.Value(),
					},
				}

				for k, v := range item.Keys() {
					if v == nil {
						continue
					}

					props[j] = append(props[j], dstore.Property{
						Name:  k,
						Value: v,
					})
				}
			}

			var t, err = db.Instance.Client().NewTransaction(ctx)
			if err != nil {
				errChan <- err
				return
			}

			_, err = t.PutMulti(keys[start:end], props[start:end])
			if err != nil {
				errChan <- err
				return
			}

			_, err = t.PutMulti(clKeys[start:end], clProps[start:end])
			if err != nil {
				errChan <- err
				return
			}

			_, err = t.Commit()
			if err != nil {
				errChan <- err
				return
			}
		}(i*amount, (i+1)*amount)
	}

	wg.Wait()

	return drainErrs(errChan).ErrorOrNil()
}

func (db *DB) DeleteChangelogs(ids ...string) error {
	return db.Instance.DeleteDocuments(context.Background(), "changelog", ids)
}

func (db *DB) Delete(id string) error {
	var (
		ctx = context.Background()
		cl  = storage.GenDeleteChangelog(id)
		err = db.Instance.UpsertDocument(ctx, "changelog", cl.ID, cl)
	)

	if err != nil {
		return err
	}

	return db.Instance.DeleteDocument(ctx, "something", id)
}
