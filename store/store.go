package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pizzahutdigital/storage/storage"
	"google.golang.org/api/iterator"
)

// Store implements Storage with helpers and some ochestration
type Store struct {
	Stores []storage.Storage
}

var (
	ErrTimeout = errors.New("timeout hit")
)

// TODO: Make these configurable per store
// will probably have to move the async stuff into the individual stores then
const (
	ReadTimeout   = 1 * time.Second
	WriteTimeout  = 2 * time.Second
	DeleteTimeout = WriteTimeout
)

func waitgroupOrTimeout(timeout time.Duration, wg *sync.WaitGroup, closeChan chan struct{}) {
	go func() {
		wg.Wait()

		select {
		case closeChan <- struct{}{}:
		}
	}()

	select {
	case <-closeChan:
		return

	case <-time.After(timeout):
		return
	}
}

func drainErrs(errChan chan error) (merr *multierror.Error) {
	close(errChan)

	for err := range errChan {
		merr = multierror.Append(merr, err)
	}

	return merr
}

// type result struct {
// 	Item *pb.Item
// 	Err  error
// }

// func getFromStore(store *storage.Storage) (storage.Item, error) {
// 	var (
// 		s   pb.Item
// 		res *result
// 		ErrTimeout = 2 * time.Second
// 	)

// 	go func() {
// 		select {
// 		case resChan <- &result{
// 			Item: &s,
// 			Err:  db.Instance.GetDocument(ctx, "something", id, &s),
// 		}:
// 		}
// 	}()

// 	for {
// 		select {
// 		case res = <-resChan:
// 			if res.Err != nil {
// 				return nil, res.Err
// 			}

// 			return object.New(res.Item.GetId(), res.Item.GetValue()), nil

// 		case <-time.After(GetTimeout):
// 			return nil, ErrTimeout
// 		}
// 	}
// }

/*
	cases := make([]reflect.SelectCase, len(chans))
	for i, ch := range chans {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	chosen, value, ok := reflect.Select(cases)
	# ok will be true if the channel has not been closed.
	ch := chans[chosen]
	msg := value.String()
*/

// Make:
//	- GetWithTimeout
//	- GetAsync
//	- do something with/for streams
func (s *Store) Get(ctx context.Context, id string) (storage.Item, error) {
	// We only _need_ a channel of size 1 but making it the len of all the
	// stores ensures that we don't block on writing
	ctx, cancelFunc := context.WithCancel(ctx)

	// Defer the cancel if for some reason it is not canceled already
	defer func() {
		select {
		case <-ctx.Done():
		default:
			cancelFunc()
		}
	}()

	var getChan = make(chan storage.Item, len(s.Stores))

	for _, store := range s.Stores {
		go func(store storage.Storage) {
			var item, err = store.Get(ctx, id)
			if err != nil {
				if !strings.Contains(err.Error(), context.Canceled.Error()) {
					// log
				}

				return
			}

			select {
			case <-ctx.Done():
				return

			default:
				cancelFunc()
				getChan <- item
				close(getChan)
			}
		}(store)
	}

	select {
	// TODO: look and see if this allocas shit everytime
	case item := <-getChan:
		// TODO: log which one won
		return item, nil

	case <-time.After(ReadTimeout):
		// TODO: log
		return nil, ErrTimeout
	}
}

func (s *Store) Set(id string, i storage.Item, sk map[string]interface{}) error {
	var (
		wg        sync.WaitGroup
		errChan   = make(chan error, len(s.Stores))
		closeChan = make(chan struct{})
	)

	defer close(closeChan)

	for _, store := range s.Stores {
		wg.Add(1)

		go func(store storage.Storage) {
			defer wg.Done()

			select {
			// If you can't write the channel then just move on
			case errChan <- store.Set(id, i, nil):

				// TODO: use an error here and lock/append to a slice
				//default:
			}
		}(store)
	}

	waitgroupOrTimeout(WriteTimeout, &wg, closeChan)

	return drainErrs(errChan)
}

func (s *Store) Delete(id string) error {
	var (
		wg        sync.WaitGroup
		errChan   = make(chan error, len(s.Stores))
		closeChan = make(chan struct{})
	)

	defer close(closeChan)

	for _, store := range s.Stores {
		wg.Add(1)

		go func(store storage.Storage) {
			defer wg.Done()

			select {
			// If you can't write the channel then just move on
			// TODO: this should not actually delete
			case errChan <- store.Delete(id):

				// TODO: use an error here and lock/append to a slice
				//default:
			}
		}(store)
	}

	waitgroupOrTimeout(DeleteTimeout, &wg, closeChan)

	return drainErrs(errChan)
}

func (s *Store) Next() (item storage.Item, err error) {
	return nil, errors.New("Not implemented")
}

// func (s *Store) Next() (item storage.Item, err error) {
// 	var resChan = make(chan *storage.Result, len(s.Stores))

// 	for _, store := range s.Stores {
// 		go func() {
// 			item, err := store.Next()

// 			select {
// 			case resChan <- &storage.Result{
// 				Item: item,
// 				Err:  err,
// 			}:
// 				if err != nil {
// 					close(resChan)
// 				}
// 			}
// 		}()
// 	}

// 	for res := range resChan {
// 		if res.Err != nil {
// 			// log
// 			continue
// 		}

// 		item = res.Item
// 		err = res.Err
// 	}

// 	return item, err
// }

func (s *Store) Iterator() (storage.Iter, error) {
	// create iterators for all the stores
	return nil, errors.New("Not implemented")

}

func (s *Store) ChangelogIterator() (storage.ChangelogIter, error) {
	// Async over all stores here and wait with a channel for the first one
	return nil, errors.New("Not implemented")
}

func (s *Store) GetLatestChangelogForObject(id string) (*storage.Changelog, error) {
	// Async over all stores here and wait with a channel for the first one
	return nil, errors.New("Not implemented")
}

/*
Basic logic for sync:
	Considering stores A and B and an object O

	A) A has O while B does not
			OR
	B) A has a newer timestamp for O than B
		- copy O from A to B

	C) B has O while A does not
			OR
	D) B has a newer timestamp for O than A
		- copy object O from B to A

	... extrapolate for multiple stores ...
*/

// func (s *Store) Sync() error {
// 	var primary storage.Storage
// 	if len(s.Stores) < 1 {
// 		return errors.New("nothing here to do dummy")
// 	}

// 	// Just assume the first one is the primary for now
// 	primary = s.Stores[0]

// 	iter, err := primary.Iterator()
// 	if err != nil {
// 		return err
// 	}

// 	// I hate declaring pointer... but this is the cheapest way to do it right now
// 	var wg = &sync.WaitGroup{}

// 	// For each item from the master's iterator ...
// 	for {
// 		var value storage.Item
// 		// Get the next item
// 		value, err = iter.Next()
// 		if err != nil {
// 			// Break if we are at the end
// 			// TODO: change this to a different error later
// 			if err == iterator.Done {
// 				break
// 			}

// 			return err
// 		}

// 		wg.Add(1)

// 		go func() {
// 			defer wg.Done()

// 			// For each non-primary store attached, attempt to async retrieve the same object
// 			for _, store := range s.Stores[1:len(s.Stores)] {
// 				wg.Add(1)
// 				go func(store storage.Storage, item storage.Item) {
// 					defer wg.Done()

// 					// Get the item from the store
// 					i, err := store.Get(item.ID())
// 					if err != nil {
// 						// TODO: need to keep a channel/sync map here
// 						fmt.Println("err", err)
// 					}

// 					// Compare the item in the store to the item from the master
// 					// TODO: Just use string for now; can't compare array
// 					// We will implement a generic Compare function later
// 					if string(item.Value()) != string(i.Value()) {
// 						fmt.Println("err",
// 							errors.New("not the same"),
// 							string(item.Value()),
// 							string(i.Value()))
// 					}
// 				}(store, value)
// 			}
// 		}()
// 	}

// 	// Wait on all the comparisons to finish
// 	// Worker pool will probably be needed so that these don't all die in heat death
// 	wg.Wait()

// 	// TODO: Will need to return a "multi-error" essentially
// 	return nil
// }

func (s *Store) Sync2() error {
	/*
		1) get a changelog iter from one table
		2) get all the changelogs from the other table by ObjectID
		3) compare changelogs
		4) perform appropriate action
	*/

	// Synchronously loop through all stores and sync them with eachother based on changelogs.
	// This will gradually get easier as you iterate through the stores because the changelogs
	// are deleted from each other after each run, preventing you from "re-syncing" based on
	// old changelogs. As such, if they are decently consistent with eachother, then most
	// changelogs will be mutually inclusive, further preventing any "re-syncing".

	for _, store := range s.Stores {
		clIter, err := store.ChangelogIterator()
		if err != nil {
			// TODO: log here
			continue
		}

		var (
			cl *storage.Changelog
			wg = &sync.WaitGroup{}
		)

		for {
			cl, err = clIter.Next()
			if err != nil {
				if err != iterator.Done {
					// TODO: log here
					// probably should set up a channel for errors
				}

				break
			}

			fmt.Println("cl", cl)

			wg.Add(1)
			go func(cl storage.Changelog) {
				defer wg.Done()

				/*
					NOTES:
					- loop through all stores - act on changelogs
					- only squash a changelog if there is unanimous resolution
					- run continuously
					- run an inmem store as well; better response times, can always get to yourself
				*/

				// fmt.Println(s.Stores[1].(*redis.DB).Instance.SScan("changelog", 0, "", 1000000))
				var latestCL = &storage.Changelog{}
				for _, store := range s.Stores[1:len(s.Stores)] {
					// TODO: if there is no changelog then we need to compare the objects themselves...
					latest, err := store.GetLatestChangelogForObject(cl.ObjectID)
					if err != nil {
						// TODO: need to do something with the error
						return
					}

					// Figure out the latest one between the other stores
					handleDiff(*latest, *latestCL)
					// delete off the younger one
				}

				// // TODO: kinda dangerous to just deref raw...
				// // Compare the latest one from all other stores to the first one found
				// handleDiff(*latestCL, cl)
				/*
					- Based on the number that comes back, copy the appropriate object over
					- Delete the changelogs
				*/

			}(*cl)
		}

		// Wait on all changelog comparisons to finish before moving to the next store
		wg.Wait()

		// Ring cycle the stores
		s.Stores = append(s.Stores[1:len(s.Stores)], s.Stores[0])
	}

	return nil
}

func handleDiff(latest, cl storage.Changelog) int {
	if latest.Timestamp > cl.Timestamp {
		// Copy object from B to A
		return -1

	} else if latest.Timestamp < cl.Timestamp {
		// Copy object from A to B
		return 1

	} else {
		// they are the same, compare the types
		// this will still return -1 or 1 in the end
		// return -1 or 1
	}

	return 0

	// TODO: implement this
}

// TODO: generically implement store like this and let the other db's be specific resolvers
// func (s *Store) Getq(q *query.Query) error {
// 	if q == nil {
// 		return errors.New("nil query")
// 	}

// 	if q.IsAsync() {
// 		return q.Start()
// 	}

// 	return q.Run()
// }
