package audit

import (
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/scottshotgg/storage/storage"
	"google.golang.org/api/iterator"
)

/*
	Audit has been provisioned a separate package because MOST implementations can just use the standard
	audit function instead of creating their own since audit does not rely on any database specific utilities.

	IF you need database specific utilities to perform your audit (ideally you should not),
	you should implement your own Audit function using this one as a reference guide.
*/

// Audit will go through the changelogs of a database and reduce them down to only the unique changelogs
// It should *ALWAYS* be called *BEFORE* sync
func Audit(db storage.Storage) (map[string]*storage.Changelog, error) {
	// Loop over all the changelogs and find the unique ones

	// Get a changelog iterator
	var iter, err = db.ChangelogIterator()
	if err != nil {
		return nil, err
	}

	const (
		workerAmount = 10
		deleteLen    = 1000
		deleteLenM1  = deleteLen - 1
	)

	var (
		cl    *storage.Changelog
		clMap = map[string]*storage.Changelog{}

		deleteChan = make(chan string, deleteLen)

		merr     *multierror.Error
		errChan  = make(chan error, deleteLen)
		doneChan = make(chan struct{})

		wg sync.WaitGroup

		mapCL *storage.Changelog
	)

	go func(err error) {
		for err = range errChan {
			merr = multierror.Append(merr, err)
		}

		close(doneChan)
	}(err)

	// Spin off [workerAmount] workers to process the deletes
	for j := 0; j < workerAmount; j++ {
		wg.Add(1)

		go func(err error) {
			defer wg.Done()

			var (
				i    int
				clID string

				deleteCLs = make([]string, deleteLen)
			)

			for clID = range deleteChan {
				// Check if we need to delete the block of changelogs
				if i == deleteLenM1 {

					// Reset i
					i = 0

					// // Fire a delete
					// go func() {
					// 	defer wg.Done()

					// TODO: check the amount at some point
					// _, err = db.Instance.HDel("changelog", deleteCLs...).Result()
					err = db.DeleteChangelogs(deleteCLs...)
					if err != nil {
						// pump into err channel
						errChan <- err
					}
				}

				deleteCLs[i] = clID
				i++
			}

			// Delete the last few if there are any
			if len(deleteCLs) != 0 {
				// TODO: check the amount at some point
				err = db.DeleteChangelogs(deleteCLs...)
				if err != nil {
					// pump into err channel
					errChan <- err
				}
			}
		}(err)
	}

	for {
		cl, err = iter.Next()
		if err != nil {
			// log

			if err == iterator.Done {
				break
			}

			return clMap, err
		}

		// Find some way to do this when not so lazy
		mapCL = clMap[cl.ObjectID]

		// If that Object ID is not already in the map then just add it and continue on
		if mapCL == nil {
			clMap[cl.ObjectID] = cl
			continue
		}

		// If the timestamp is greater than what is already in the map then that is a new timestamp
		if cl.Timestamp > mapCL.Timestamp {
			// Queue the old changelog from the map for deletion
			deleteChan <- mapCL.ID

			// Overwrite the changelog in the map
			*mapCL = *cl
			continue
		}

		// Queue that changelog for deletion
		deleteChan <- cl.ID
	}

	// Close the channel since we should be done at this point
	close(deleteChan)

	// Wait on all workers to finish
	wg.Wait()

	close(errChan)

	// Wait until the error worker is done
	<-doneChan

	return clMap, nil
}
