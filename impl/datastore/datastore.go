package datastore

import (
	"context"
	"errors"
	"sync"
	"time"

	dstore "cloud.google.com/go/datastore"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/scottshotgg/storage/audit"
	"github.com/scottshotgg/storage/object"
	"github.com/scottshotgg/storage/storage"
	"google.golang.org/api/iterator"
)

// DB implements Storage from the storage package
type DB struct {
	Instance *dstore.DSInstance
}

const (
	GetTimeout = 1 * time.Second
)

var (
	ErrTimeout                   = errors.New("Timeout")
	ErrNotImplemented            = errors.New("Not implemented")
	ErrTransactionAmountExceeded = errors.New("Only 12 items can be batched; this is a Google Datastore limit")
)

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

func (db *DB) SetMulti(ctx context.Context, items []storage.Item) error {
	const (
		amount   = 12
		amountM1 = amount - 1
	)

	var (
		wg sync.WaitGroup

		// The amount of workers here made no difference past about 10% of the total length
		// Adding 1 so that there is atleast 1 worker in the case that the calculation resolves to 0
		workerChan = make(chan struct{}, int64(float64(len(items))*.25)+1)
		errChan    = make(chan error, len(items))

		leftover = len(items) % amountM1
	)

	if leftover != 0 {
		// Spawn a batch to take care of the leftover
		go func() {
			var (
				start = len(items) - leftover
				end   = len(items)
				err   = db.BatchTransaction(ctx, items[start:end])
			)

			if err != nil {
				errChan <- err
			}

			<-workerChan
			wg.Done()
		}()
	}

	// Close the workerChan at the end
	defer close(workerChan)

	// Since each goroutine takes `amount` objects, at max we need the len(objects) / `amount-1`
	for i := 0; i < len(items)/amount; i++ {
		wg.Add(1)
		workerChan <- struct{}{}

		// Check the context before launching the goroutine
		select {
		case <-ctx.Done():
			return context.Canceled

		default:
		}

		// Batch each set of `amount` to be processed by a separate goroutine
		// Do NOT change this to use a channel; the memory usage will potentially be more efficient
		// but you will incur a severe performance pentalty for locking and unlocking at each read/write of the channel
		go func(i int) {
			var (
				start = i * amount
				end   = (i + 1) * amount
				err   = db.BatchTransaction(ctx, items[start:end])
			)

			if err != nil {
				errChan <- err
			}

			<-workerChan
			wg.Done()
		}(i)
	}

	// Wait for all transactions to finish
	wg.Wait()

	// If there were any errors then return that
	return drainErrs(errChan).ErrorOrNil()
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

func drainErrs(errChan chan error) (merr *multierror.Error) {
	close(errChan)

	for err := range errChan {
		merr = multierror.Append(merr, err)
	}

	return merr
}

// Batch transaction takes up to 12 items and generates changelogs, throws them all into a transaction and commits it
func (db *DB) BatchTransaction(ctx context.Context, items []storage.Item) error {
	if len(items) > 12 {
		return ErrTransactionAmountExceeded
	}

	var (
		clKeys  = make([]*dstore.Key, len(items))
		clProps = make([]dstore.PropertyList, len(items))

		keys  = make([]*dstore.Key, len(items))
		props = make([]dstore.PropertyList, len(items))
	)

	// TODO: put this in a function
	// Loop over the defined start and end
	for j, item := range items {
		// Check the context while we are looping inside the goroutine
		select {
		case <-ctx.Done():
			return context.Canceled

		default:
		}

		// Generate a new changelog for the upsert
		var cl = storage.GenInsertChangelog(item)

		// Create the changelog key
		clKeys[j] = genChangelogKey(cl)

		// Create the changelog props
		clProps[j] = genChangelogProps(cl)

		// Create the object key
		keys[j] = genItemKey(item)

		// Create the object props
		props[j] = genItemProps(item)
	}

	// Create a new transaction from the client
	var t, err = db.Instance.Client().NewTransaction(ctx)
	if err != nil {
		return err
	}

	// Put the object keys and props into the transaction
	_, err = t.PutMulti(keys, props)
	if err != nil {
		return err
	}

	// Put the changelog keys and props into the transaction
	_, err = t.PutMulti(clKeys, clProps)
	if err != nil {
		return err
	}

	// Commit our changes; this is where operations ACTUALLY take place
	// The ignored keys are essentially promises that can be resolved, similar to using Mongo through Node
	// But we don't care about resolving them
	_, err = t.Commit()

	return err
}

// genChangelogKey generates a Datastore name based key for a changelog
func genChangelogKey(cl *storage.Changelog) *dstore.Key {
	return dstore.NameKey("changelog", cl.ID, nil)
}

// genItemProps takes a changelog and formats it into a PropertyList for Datastore
func genChangelogProps(cl *storage.Changelog) dstore.PropertyList {
	return dstore.PropertyList{
		propFromKV("ID", cl.ID),
		propFromKV("ObjectID", cl.ObjectID),
		propFromKV("Timestamp", cl.Timestamp),
		propFromKV("Type", cl.Type),
	}
}

// getItemKey generates a Datastore name based key for an item
func genItemKey(item storage.Item) *dstore.Key {
	return dstore.NameKey("something", item.ID(), nil)
}

// genItemProps takes an item and formats it into a PropertyList for Datastore
func genItemProps(item storage.Item) dstore.PropertyList {
	var props = dstore.PropertyList{
		propFromKV("id", item.ID()),
		propFromKV("timestamp", item.Timestamp()),
		propFromKV("value", item.Value()),
	}

	// If any keys came along with the object then insert those into the props as well
	for k, v := range item.Keys() {
		// if v != nil {
		props = append(props, propFromKV(k, v))
		// }
	}

	return props
}

// propFromKV takes a key-value pair and returns a Datastore Property
func propFromKV(key string, value interface{}) dstore.Property {
	return dstore.Property{
		Name:  key,
		Value: value,
	}
}
