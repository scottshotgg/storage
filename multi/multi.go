package multi

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	"google.golang.org/api/iterator"

	"github.com/scottshotgg/storage/storage"
)

// Store implements Storage with helpers and some ochestration
type DB struct {
	sync.RWMutex

	// Stores    []storage.Storage
	Stores map[string]storage.Storage
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

// // NewRedis creates a new Redis DB and checks it's connection
// func NewRedis(opts *redigo.Options) (Storage, error) {
// 	return redis.New(opts)
// }

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

/*
	Have a meta table that has all of the IDs for the DB
*/

func New(db storage.Storage) (*DB, error) {
	// var id = v4.String()
	// if id == "" {
	// 	// TODO: not sure what to do here
	// 	return nil,
	// }

	return &DB{
		Stores: map[string]storage.Storage{
			db.ID(): db,
		},
	}, nil
}

// TODO: need to add mutexes into other places

func (db *DB) Add(s storage.Storage) {
	db.Lock()
	db.Unlock()

	if db.Stores[s.ID()] != nil {
		return
	}

	db.Stores[s.ID()] = s
}

// Make:
//	- GetWithTimeout
//	- GetAsync
//	- do something with/for streams
func (db *DB) Get(ctx context.Context, id string) (storage.Item, error) {
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

	var getChan = make(chan storage.Item, len(db.Stores))

	for _, store := range db.Stores {
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

func (db *DB) Set(ctx context.Context, i storage.Item) error {
	var (
		wg        sync.WaitGroup
		errChan   = make(chan error, len(db.Stores))
		closeChan = make(chan struct{})
	)

	defer close(closeChan)

	for _, store := range db.Stores {
		wg.Add(1)

		go func(store storage.Storage) {
			defer wg.Done()

			select {
			// If you can't write the channel then just move on
			case errChan <- store.Set(ctx, i):

				// TODO: use an error here and lock/append to a slice
				//default:
			}
		}(store)
	}

	waitgroupOrTimeout(WriteTimeout, &wg, closeChan)

	return drainErrs(errChan)
}

func (db *DB) Delete(id string) error {
	var (
		wg        sync.WaitGroup
		errChan   = make(chan error, len(db.Stores))
		closeChan = make(chan struct{})
	)

	defer close(closeChan)

	for _, store := range db.Stores {
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

func (db *DB) Next() (item storage.Item, err error) {
	return nil, errors.New("Not implemented")
}

// func (db *DB) Next() (item storage.Item, err error) {
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

func (db *DB) Iterator() (storage.Iter, error) {
	// create iterators for all the stores
	return nil, errors.New("Not implemented")
}

func (db *DB) ChangelogIterator() (storage.ChangelogIter, error) {
	return nil, errors.New("Not implemented")
}

/*
	return ALL changelogs, make QuickSync find the exclusives and then sort and delete
*/
func (db *DB) GetLatestChangelogForObject(id string) (*storage.Changelog, error) {
	// Async over all stores here and wait with a channel for the first one
	return nil, errors.New("Not implemented")
}

func (db *DB) Audit() (map[string]*storage.Changelog, error) {
	var (
		wg sync.WaitGroup

		// TODO: could run into problems here if one comes up while we are auditing
		errChan   = make(chan error, len(db.Stores)+10)
		clMapChan = make(chan map[string]*storage.Changelog, len(db.Stores)+10)

		clMaps = map[string]*storage.Changelog{}
	)

	// Process all of the changelogs to find an overall unique set
	go func() {
		var (
			clMap map[string]*storage.Changelog

			// id string
			cl *storage.Changelog

			// mapCL *storage.Changelog
		)

		// for clMap = range clMapChan {
		// 	for id, cl = range clMap {
		// 		mapCL = clMaps[id]

		// 		// If that Object ID is not already in the map then just add it and continue on
		// 		if mapCL == nil {
		// 			clMap[cl.ObjectID] = cl
		// 			continue
		// 		}

		// 		// If the timestamp is greater than what is already in the map then that is a new timestamp
		// 		if cl.Timestamp > mapCL.Timestamp {
		// 			// Overwrite the changelog in the map
		// 			*mapCL = *cl
		// 			continue
		// 		}
		// 	}
		// }

		// Just add all of the changelogs together, keyed by changelog ID instead of object ID
		for clMap = range clMapChan {
			for _, cl = range clMap {
				// mapCL = clMaps[cl.ID]
				clMap[cl.ID] = cl

				// // If that Object ID is not already in the map then just add it and continue on
				// if mapCL == nil {
				// 	clMap[cl.ObjectID] = cl
				// 	continue
				// }

				// // If the timestamp is greater than what is already in the map then that is a new timestamp
				// if cl.Timestamp > mapCL.Timestamp {
				// 	// Overwrite the changelog in the map
				// 	*mapCL = *cl
				// 	continue
				// }
			}
		}
	}()

	// Launch audit for each store asynchronously
	for _, store := range db.Stores {
		wg.Add(1)

		go func(store storage.Storage) {
			defer wg.Done()

			// Tell the store to audit itself
			var clMap, err = store.Audit()
			if err != nil {
				errChan <- err
				return
			}

			// Put the map into the channel for processing
			clMapChan <- clMap
		}(store)
	}

	// TODO:
	// reduce down changelogs to just the unique items
	// have some way to link it back to the original db
	// if someone needs something then copy over from that db

	// Wait on audit workers
	wg.Wait()

	// Drain the errors from the channel
	return clMaps, drainErrs(errChan).ErrorOrNil()
}

func (db *DB) QuickSync() error {
	// Audit every store first to reduce the number of changelogs we have to look through
	var _, err = db.Audit()
	if err != nil {
		// log
		return err
	}

	var storesCopy []storage.Storage

	// Copy the stores incase one is added later on
	for _, store := range db.Stores {
		storesCopy = append(storesCopy, store)
	}

	// Iterate over all of the stores multiple times as changelogs are copied back and forth
	for i := 0; i < len(storesCopy)*len(storesCopy); i++ {
		// Assume the first store as the master for this iteration
		var (
			master     = storesCopy[0]
			waitLength time.Duration
		)

		// Range over all "slaves" of that master storage
		for _, slave := range storesCopy[1:] {
			// Iterate over all the changelogs
			var clIter, err = master.ChangelogIterator()
			if err != nil {
				return err
			}

			var (
				wg          sync.WaitGroup
				workerChan  = make(chan struct{}, 100)
				objectIDMap = map[string]*struct{}{}
			)

			// While we still have more changelogs ...
			// TODO: to speed this up we could define a time range, get all changelogs in that
			// time range and dispatch them
			for {
				// Get a changelog
				var cl, err = clIter.Next()
				if err != nil {
					if err == iterator.Done {
						waitLength = time.Duration((len(objectIDMap) / 1000)) * time.Second
						break
					}

					return err
				}

				// Check if we have already seen this objectID
				if objectIDMap[cl.ObjectID] != nil {
					continue
				}

				// Mark that we have seen this objectID
				objectIDMap[cl.ObjectID] = &struct{}{}

				wg.Add(1)
				workerChan <- struct{}{}
				go func() {
					defer func() {
						<-workerChan
						wg.Done()
					}()

					var err = processChangelogs(&wg, cl.ObjectID, master, slave)
					if err != nil {
						// log
					}
				}()
			}

			wg.Wait()
		}

		// Ring cycle through the stores
		storesCopy = append(storesCopy[1:len(storesCopy)], storesCopy[0])

		fmt.Println("Waiting ", waitLength, " to catch up")

		// Sleep for a bit to let the other store catch up with the load before we start again
		time.Sleep(waitLength)
	}

	return nil
}

// Maybe we should make an AuditFrom ...

func (db *DB) QuickSyncWith() error {
	// Audit every store first to reduce the number of changelogs we have to look through
	// clMap will contain ALL unique changelogs from ALL stores
	var clMap, err = db.Audit()
	if err != nil {
		// log
		return err
	}

	var (
		uniqueCLs = map[string]*storage.Changelog{}
		deleteCLs = map[string]*storage.Changelog{}
		// storesCopy []storage.Storage

		mapCL *storage.Changelog
	)

	// Iterate over all of the unique changelogs from all of the stores
	for _, cl := range clMap {
		mapCL = uniqueCLs[cl.ObjectID]

		// If that Object ID is not already in the map then just add it and continue on
		if mapCL == nil {
			uniqueCLs[cl.ObjectID] = cl
			continue
		}

		// If the timestamp is greater than what is already in the map then that is a new timestamp
		if cl.Timestamp > mapCL.Timestamp {
			// Queue the old changelog from the map for deletion
			deleteCLs[mapCL.ID] = mapCL

			// Overwrite the changelog in the map
			*mapCL = *cl
			continue
		}

		// Queue that changelog for deletion
		deleteCLs[cl.ID] = cl
	}

	// Just delete them all at the end
	// need to fix this
	defer func() {
		// put an if check here
		var err = db.deleteAllChangelogs(deleteCLs)
		if err != nil {
			// log
			// Need to put something here to not delete this changelog
		}
	}()

	const (
		// Only allow 1000 total workers
		totalWorkerAmount = 1000

		// // Calculate the amount of workers that will be allocated to each store
		// storeWorkerAmount / len(s.Stores)
	)

	var (
		wg sync.WaitGroup
		// workerChanMap = map[string]chan struct{}{}
		workerChan = make(chan struct{}, totalWorkerAmount)

		// Cache the items as we go through the stores
		// Consider doing a time.After to clear this
		itemCacheMap = map[string]storage.Item{}
	)

	// Make a map of the unprocessable items and then loop over the deleteChangelogs and any
	// with the same objectID should be kept (deleted from the deleteMap)

	// Process the changelogs for each store
	for dbID, store := range db.Stores {
		// Declare a context for each store
		var ctx = context.Background()
		wg.Add(1)

		// go func(ctx context.Context, dbID string, store storage.Storage) {
		// 	defer wg.Done()

		// Iterate over all of the unique changelogs
		for _, cl := range uniqueCLs {
			// Don't check against ourselves... for now
			if cl.DBID == dbID {
				continue
			}

			wg.Add(1)
			workerChan <- struct{}{}

			// Go func each retrieval
			go func(cl *storage.Changelog) {
				defer func() {
					wg.Done()
					<-workerChan
				}()

				// Get the item from the store we are comparing against
				var item, err = store.Get(ctx, cl.ObjectID)
				if err != nil {
					// log
					// Need to put something here to not delete this changelog
					return
				}

				// If the item that we got is older than the changelog for that item
				// then go get the item from the store that recorded the changelog
				if item.GetTimestamp() < cl.Timestamp {
					// Check the cache for the item
					item = itemCacheMap[cl.ObjectID]

					// If the item is nil then the item needs to be cached
					if item == nil {
						// Get the item to be cached from the store that recorded the changelog
						item, err = db.Stores[cl.DBID].Get(ctx, cl.ObjectID)
						if err != nil {
							// log
							// Need to put something here to not delete this changelog
							// This is BAD
							return
						}

						// Set the item in the cache
						itemCacheMap[cl.ObjectID] = item
					}

					// Compare for the item that we got back to the changelog for good measure
					// The item timestamp SHOULD be either:
					//	1) the same if the item has not changed in the meantime
					//	2) greater than if the item has, in which case we will still copy it
					if item.GetTimestamp() < cl.Timestamp {
						// wtf
						// Could put something here to check the two object against each other...
						return
					}

					// Might consider checking here right before upserting to it
					err = store.Set(ctx, item)
					if err != nil {
						// log
						// Need to put something here to not delete this changelog
						return
					}
				}
			}(cl)
			// Don't do anything if it's newer than the timestamp
		}
		// }(ctx, dbID, store)
	}

	// Just for now
	return nil
}

// Sync attempts to look through all objects of all stores and distribute the most up to date set
// This should really only be run at times when there is low writes
func (db *DB) Sync() error {
	var (
		storesCopy []storage.Storage
		wg         sync.WaitGroup
		merr       *multierror.Error

		errChan    = make(chan error, 100)
		workerChan = make(chan struct{}, 10)
	)

	// Copy the stores incase one is added later on
	for _, store := range db.Stores {
		storesCopy = append(storesCopy, store)
	}

	// Spin off a goroutine to process the errors
	go func() {
		// The assumption at this point is that all errors are non-nil
		var err error
		for err = range errChan {
			merr = multierror.Append(merr, err)
		}
	}()

	// Iterate over all of the stores
	for i := 0; i < len(storesCopy); i++ {
		// Assume the first store as the master for this iteration
		var (
			item storage.Item

			// Declare the master for this iteration
			master = storesCopy[0]

			// Get an item iterator from the master
			iter, err = master.Iterator()
		)

		if err != nil {
			errChan <- err

			// Just skip the entire verification if we can't get the iterator for some reason
			continue
		}

		// Drain the iterator
		for {
			item, err = iter.Next()
			if err != nil {
				// log

				if err == iterator.Done {
					break
				}

				return err
			}

			// Reserve a spot for processing items BEFORE spawning a goroutine
			// This will cut down on internal runtime coordination and memory/processing requirements
			wg.Add(1)
			workerChan <- struct{}{}

			// Spin off a worker
			go func(item storage.Item) {
				// Check the stores for that item
				var err = checkStoresForItem(item, storesCopy)

				// Enable another worker before processing the error
				<-workerChan
				wg.Done()

				// Process the error
				if err != nil {
					errChan <- err
				}
			}(item)
		}

		// Wait for all the workers from this store to finish
		wg.Wait()

		// Ring cycle the stores
		storesCopy = append(storesCopy[1:len(storesCopy)], storesCopy[0])
	}

	// Wait for good measure
	wg.Wait()

	// Close the channels
	close(workerChan)

	close(errChan)

	// Return either error or nil
	return merr.ErrorOrNil()
}

func checkStoresForItem(item storage.Item, storesCopy []storage.Storage) error {
	var (
		master = storesCopy[0]
		wg     sync.WaitGroup
		ctx    = context.Background()

		// Spawn 10 workers to take care of the stores
		workerChan = make(chan struct{}, 10)
		errChan    = make(chan error, 100)

		merr *multierror.Error
	)

	// Spin off a goroutine to process the errors
	go func() {
		// The assumption at this point is that all errors are non-nil
		var err error
		for err = range errChan {
			multierror.Append(merr, err)
		}
	}()

	// Range over all "slaves" of that master storage
	for i, slave := range storesCopy[1:] {
		wg.Add(1)
		workerChan <- struct{}{}

		go func(i int) {
			// Do something with the error later
			var item2, err = slave.Get(ctx, item.GetID())

			defer func() {
				wg.Done()
				<-workerChan
			}()

			if err != nil {
				// log
				errChan <- err
				return
			}

			// If the slaves timestamp is less than that of the master then copy master -> slave
			if item2.GetTimestamp() < item.GetTimestamp() {
				err = slave.Set(ctx, item)
				// If the masters timestamp is less that that of the slave then copy slave -> master
			} else if item.GetTimestamp() < item2.GetTimestamp() {
				err = master.Set(ctx, item2)
			}
			// else we are assuming they are the same and/or that we can't do anything about it

			if err != nil {
				// log
				errChan <- err
				return
			}
		}(i)
	}

	// Wait until all workers are done
	wg.Wait()

	// Close the channels
	close(workerChan)
	close(errChan)

	// Return either error or nil
	return merr.ErrorOrNil()
}

func processChangelogs(wg *sync.WaitGroup, objectID string, master, slave storage.Storage) error {
	var (
		wg2       sync.WaitGroup
		cls       []storage.Changelog
		item      storage.Item
		err, err2 error
	)

	wg2.Add(1)
	go func() {
		defer wg2.Done()

		// Get all changelogs for that objectID from the master
		cls, err = master.GetChangelogsForObject(objectID)
	}()

	wg2.Add(1)
	go func() {
		defer wg2.Done()

		// Get the item from the slave to see if checking the changelog is even applicable
		item, err2 = slave.Get(context.Background(), objectID)
	}()

	wg2.Wait()

	if err != nil {
		return err
	}

	if err2 != nil {
		return err2
	}

	// If the master does not have any changelogs relating to this object then skip for now
	if len(cls) == 0 {
		return nil
	}

	// Compare and find the latest changelog out of all the ones related to that ObjectID;
	// we will delete the other later
	var (
		latest        = getLatest(cls)
		refetchLatest bool
		ctx           = context.Background()
	)

	// If the master timestamp is greater than the slaves object timestamp then update the slave
	if item.GetTimestamp() < latest.Timestamp {
		item, err = master.Get(ctx, objectID)
		if err != nil {
			return err
		}

		// Update the slaves object
		err = slave.Set(ctx, item)
		if err != nil {
			return err
		}
	} else if item.GetTimestamp() > latest.Timestamp {
		// Update the masters object
		err = master.Set(ctx, item)
		if err != nil {
			return err
		}

		refetchLatest = true
	}

	// else {
	// they are the same
	// }

	// Delete all master timestamps from the master that were retrieved related to this object
	// We processed this object so do this regardless of whether it was used to update
	// wg.Add(1)
	// go func() {
	// defer wg.Done()

	var clIDs []string
	for i := range cls {
		clIDs = append(clIDs, cls[i].ID)
	}

	err = master.DeleteChangelogs(clIDs...)
	if err != nil {
		// TODO: probably shouldn't return here
		return err
	}
	// }()

	// TODO: Might be able to somehow pipe all of this until the end
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Fetch changelogs for that object from slave
		var cls, err = slave.GetChangelogsForObject(objectID)
		if err != nil {
			// log
			return
		}

		if refetchLatest {
			latest = getLatest(cls)
		}

		// Get only the changelogs older than the one we processed
		cls = getOlder(cls, *latest)

		var clIDs []string
		for i := range cls {
			clIDs = append(clIDs, cls[i].ID)
		}

		// Delete all older changelogs
		err = slave.DeleteChangelogs(clIDs...)
		if err != nil {
			// log
		}
	}()

	return nil
}

func getOlder(cls []storage.Changelog, compare storage.Changelog) []storage.Changelog {
	var cls2 []storage.Changelog

	for _, cl := range cls {
		if compare.Timestamp >= cl.Timestamp {
			cls2 = append(cls2, cl)
		}
	}

	return cls2
}

func getLatest(cls []storage.Changelog) *storage.Changelog {
	var latest storage.Changelog

	for _, cl := range cls {
		if cl.Timestamp > latest.Timestamp {
			latest = cl
		}
	}

	return &latest
}

func handleDiff(latest, cl storage.Changelog) int {
	fmt.Println("latest, cl", latest, cl)

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

func (db *DB) deleteAllChangelogs(clMap map[string]*storage.Changelog) error {
	var (
		dbToCLIDMap = map[string][]string{}
		wg          sync.WaitGroup

		// errChan = make(chan error, len(s.Stores))
	)

	// // Fill the map with all the stores so we don't have to check
	// for _, store := range s.Stores {
	// 	dbToCLIDMap[store.ID()] = []string{}
	// }

	for id, cl := range clMap {
		dbToCLIDMap[cl.DBID] = append(dbToCLIDMap[cl.DBID], id)
	}

	for id, cls := range dbToCLIDMap {
		wg.Add(1)

		// Async delete all the changelogs
		go func(id string, cls []string) {
			defer wg.Done()

			var err = db.Stores[id].DeleteChangelogs(cls...)
			if err != nil {
				// log
			}
		}(id, cls)
	}

	wg.Wait()

	return nil
}

// TODO: generically implement store like this and let the other db's be specific resolvers
// func (db *DB) Getq(q *query.Query) error {
// 	if q == nil {
// 		return errors.New("nil query")
// 	}

// 	if q.IsAsync() {
// 		return q.Start()
// 	}

// 	return q.Run()
// }

/*
	Changelog notes

	get changelog from db1
	find latest changelog1

	get changelogs from db2
	find latest changelog2

	compare changelog1 to changelog2
	if changelog1 is newer:
		- get objects from db1 and db2
		- compare changelog timestamp to object timestamps in db1 and db2
		- if still newer than db2:
			- update db2 with object
*/
