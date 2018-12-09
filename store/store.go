package store

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/pizzahutdigital/storage"
	"google.golang.org/api/iterator"
)

// Store implements Storage with helpers and some ochestration
type Store struct {
	Stores []storage.Storage
}

// func get(s storage.Storage, id string) (<-chan storage.Item, error) {
// 	// set up channel
// 	// call s.Get
// 	// write to channel

// 	getChan := make(chan storage.Item, 1)

// 	go func() {
// 		fmt.Println("stuff")
// 		fmt.Println(s.Get())
// 	}()
// }

func (s *Store) Get(id string) (storage.Item, error) {
	// Async over all stores here and wait with a channel for the first one
	// return nil, errors.New("Not implemented")
	// will need an errChan

	var (
		// wg      = &sync.WaitGroup{}
		getChan = make(chan storage.Item, 10)
		// resChan = make(chan storage.Item, 1)
	)

	// Read the response from all the stores
	go func() {
		for ret := range getChan {
			if ret.Err == nil {
				// TODO: log
				continue
			}

			fmt.Println("ret", ret.Item)
		}
	}()

	// Read from all stores
	var (
		ctx, cancelFunc = context.WithTimeout(context.Background(), time.Duration(30*time.Second))

		cleanup = func() {
			timeoutFunc()
			cancelFunc()
			close(getChan)
		}
	)

	for _, store := range s.Stores {
		go func(s storage.Storage) {
			item, err := s.Get()

			select {
			case <-ctx.Done():
				return
			}

			getChan <- item

		}(store)
	}

	// can't close here and you can't wait here or else
	// you will never close this if one is down
	// close(getChan)

	// return <-
}

func (s *Store) Set(id string, i storage.Item) error {
	// Async over all stores here and wait with a channel for the first one
	return errors.New("Not implemented")
}

func (s *Store) Delete(id string) error {
	// Async over all stores here and wait with a channel for the first one
	return errors.New("Not implemented")
}

func (s *Store) Iterator() (Iter, error) {
	// Async over all stores here and wait with a channel for the first one
	return nil, errors.New("Not implemented")
}

func (s *Store) ChangelogIterator() (changelog.ChangelogIter, error) {
	// Async over all stores here and wait with a channel for the first one
	return nil, errors.New("Not implemented")
}

func (s *Store) GetLatestChangelogForObject(id string) (*changelog.Changelog, error) {
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

func (s *Store) Sync() error {
	var primary storage.Storage
	if len(s.Stores) < 1 {
		return errors.New("nothing here to do dummy")
	}

	// Just assume the first one is the primary for now
	primary = s.Stores[0]

	iter, err := primary.Iterator()
	if err != nil {
		return err
	}

	// I hate declaring pointer... but this is the cheapest way to do it right now
	var wg = &sync.WaitGroup{}

	// For each item from the master's iterator ...
	for {
		var value storage.Item
		// Get the next item
		value, err = iter.Next()
		if err != nil {
			// Break if we are at the end
			// TODO: change this to a different error later
			if err == iterator.Done {
				break
			}

			return err
		}

		wg.Add(1)

		go func() {
			defer wg.Done()

			// For each non-primary store attached, attempt to async retrieve the same object
			for _, store := range s.Stores[1:len(s.Stores)] {
				wg.Add(1)
				go func(store storage.Storage, item storage.Item) {
					defer wg.Done()

					// Get the item from the store
					i, err := store.Get(item.ID())
					if err != nil {
						// TODO: need to keep a channel/sync map here
						fmt.Println("err", err)
					}

					// Compare the item in the store to the item from the master
					// TODO: Just use string for now; can't compare array
					// We will implement a generic Compare function later
					if string(item.Value()) != string(i.Value()) {
						fmt.Println("err",
							errors.New("not the same"),
							string(item.Value()),
							string(i.Value()))
					}
				}(store, value)
			}
		}()
	}

	// Wait on all the comparisons to finish
	// Worker pool will probably be needed so that these don't all die in heat death
	wg.Wait()

	// TODO: Will need to return a "multi-error" essentially
	return nil
}

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
			cl *Changelog
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
			go func(cl Changelog) {
				defer wg.Done()

				/*
					NOTES:
					- loop through all stores - act on changelogs
					- only squash a changelog if there is unanimous resolution
					- run continuously
					- run an inmem store as well; better response times, can always get to yourself
				*/

				// fmt.Println(s.Stores[1].(*redis.DB).Instance.SScan("changelog", 0, "", 1000000))
				var latestCL = &changelog.Changelog{}
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

func handleDiff(latest, cl Changelog) int {
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
