package store_test

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/pizzahutdigital/storage/store"

	dstore "github.com/pizzahutdigital/datastore"
	"github.com/pizzahutdigital/storage/datastore"

	redigo "github.com/go-redis/redis"
	"github.com/pizzahutdigital/storage/redis"
	"github.com/pizzahutdigital/storage/test"
)

var (
	s store.Store
)

func init() {
	db := datastore.DB{
		Instance: &dstore.DSInstance{},
	}

	// initialize Datastore client session for mythor metadata
	err := db.Instance.Initialize(dstore.DSConfig{
		Context:            context.Background(),
		ServiceAccountFile: "/Users/sgg7269/Documents/serviceAccountFiles/ds-serviceaccount.json",
		ProjectID:          "phdigidev",
		Namespace:          "storage_test",
	})
	if err != nil {
		log.Fatalf("err %+v", err)
	}

	s.Stores = append(s.Stores, &db)

	db2 := redis.DB{
		Instance: redigo.NewClient(&redigo.Options{
			Addr: "localhost:6379",
			// Password:   os.Getenv("RP"),
			MaxRetries: 10,
			// TLSConfig: we should set this up
			PoolSize:    1000,
			ReadTimeout: time.Minute,
			// ReadTimeout: -1,
			IdleTimeout: -1,
			// DialTimeout:
		}),
	}
	_, err = db2.Instance.Ping().Result()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	s.Stores = append(s.Stores, &db2)

	var wg = &sync.WaitGroup{}

	// Skip the primary store for now
	for i, s := range s.Stores[1:len(s.Stores)] {
		for _, obj := range test.Objs {
			wg.Add(1)
			test.WorkerChan <- struct{}{}

			go func(i int, s store.Storage, obj *store.Object) {
				defer func() {
					wg.Done()
					<-test.WorkerChan
				}()

				err := s.Set(fmt.Sprintf("some_id_%d", i), obj)
				if err != nil {
					log.Fatalf("%+v", err)
				}
			}(i, s, obj)
		}
	}

	wg.Wait()
}

func TestSync(t *testing.T) {
	err := s.Sync()
	if err != nil {
		t.Errorf("%+v", err)
	}
}

func TestSync2(t *testing.T) {
	err := s.Sync2()
	if err != nil {
		t.Errorf("%+v", err)
	}
}
