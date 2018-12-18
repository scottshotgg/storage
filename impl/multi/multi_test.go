package multi_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/scottshotgg/storage/impl/multi"

	dstore "github.com/pizzahutdigital/datastore"
	"github.com/scottshotgg/storage/impl/datastore"

	redigo "github.com/go-redis/redis"
	"github.com/scottshotgg/storage/impl/redis"
)

var (
	s multi.DB
)

func init() {
	var (
		db = datastore.DB{
			Instance: &dstore.DSInstance{},
		}

		// initialize Datastore client session for mythor metadata
		err = db.Instance.Initialize(dstore.DSConfig{
			Context:            context.Background(),
			ServiceAccountFile: "/Users/sgg7269/Documents/serviceAccountFiles/ds-serviceaccount.json",
			ProjectID:          "phdigidev",
			Namespace:          "storage_test",
		})
	)

	if err != nil {
		log.Fatalf("err %+v", err)
	}

	s.Stores = append(s.Stores, &db)

	var db2 = redis.DB{
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

	// var wg = &sync.WaitGroup{}

	// // Skip the primary store for now
	// for i, s := range s.Stores[1:len(s.Stores)] {
	// 	for _, obj := range test.Objs {
	// 		wg.Add(1)
	// 		test.WorkerChan <- struct{}{}

	// 		go func(i int, s storage.Storage, obj *object.Object) {
	// 			defer func() {
	// 				wg.Done()
	// 				<-test.WorkerChan
	// 			}()

	// 			err := s.Set(fmt.Sprintf("some_id_%d", i), obj, nil)
	// 			if err != nil {
	// 				log.Fatalf("%+v", err)
	// 			}
	// 		}(i, s, obj)
	// 	}
	// }

	// wg.Wait()
}

func TestGet(t *testing.T) {
	var ctx = context.Background()

	var item, err = s.Get(ctx, "some_id_59")
	if err != nil {
		t.Fatalf("err: %+v", err)
	}

	fmt.Println("item", item)
}

// func TestSync(t *testing.T) {
// 	var err = s.Sync()
// 	if err != nil {
// 		t.Errorf("%+v", err)
// 	}
// }

// func TestSync2(t *testing.T) {
// 	err := s.Sync2()
// 	if err != nil {
// 		t.Errorf("%+v", err)
// 	}
// }

func TestAudit(t *testing.T) {
	var err = s.Audit()
	if err != nil {
		t.Errorf("%+v", err)
	}
}

// func TestSync(t *testing.T) {
// 	err := s.Sync()
// 	if err != nil {
// 		t.Errorf("%+v", err)
// 	}
// }
