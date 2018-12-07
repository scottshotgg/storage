package redis_test

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	redigo "github.com/go-redis/redis"
	"github.com/pizzahutdigital/storage/redis"
	"github.com/pizzahutdigital/storage/store"
	"github.com/pizzahutdigital/storage/test"
	"google.golang.org/api/iterator"
)

func init() {
	db := redis.DB{
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

	_, err := db.Instance.Ping().Result()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	test.DB = &db
}

func TestSet(t *testing.T) {
	var wg = &sync.WaitGroup{}

	for _, obj := range test.Objs {
		wg.Add(1)
		test.WorkerChan <- struct{}{}

		go func(obj *store.Object) {
			defer func() {
				wg.Done()
				<-test.WorkerChan
			}()

			err := test.DB.Set(obj.ID(), obj)
			if err != nil {
				t.Errorf("err %+v", err)
			}
		}(obj)
	}

	wg.Wait()
}

func TestGet(t *testing.T) {
	var wg = &sync.WaitGroup{}

	for i := 0; i < test.AmountOfTests; i++ {
		wg.Add(1)
		test.WorkerChan <- struct{}{}

		go func(i int) {
			defer func() {
				wg.Done()
				<-test.WorkerChan
			}()

			item, err := test.DB.Get(fmt.Sprintf("some_id_%d", i))
			if err != nil {
				t.Fatalf("err %+v", err)
			}

			var testt test.Test
			err = json.Unmarshal(item.Value(), &testt)
			if err != nil {
				t.Fatalf("err %+v", err)
			}
		}(i)
	}

	wg.Wait()
}

func TestIter(t *testing.T) {
	iter, err := test.DB.Iterator()
	if err != nil {
		t.Errorf("err %+v", err)
	}

	var (
		item  store.Item
		testt test.Test
	)

	for {
		item, err = iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}

			t.Fatalf("err %+v", err)
		}

		err = json.Unmarshal(item.Value(), &testt)
		if err != nil {
			t.Fatalf("err %+v", err)
		}
	}
}
