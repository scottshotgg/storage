package redis_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	redigo "github.com/go-redis/redis"
	"github.com/pizzahutdigital/storage/object"
	"github.com/pizzahutdigital/storage/redis"
	"github.com/pizzahutdigital/storage/storage"
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

			item, err := test.DB.Get(context.Background(), fmt.Sprintf("some_id_%d", i))
			if err != nil {
				t.Fatalf("err %+v", err)
			}

			fmt.Println("item", item.ID(), item.Timestamp())

			var testt test.Test
			err = json.Unmarshal(item.Value(), &testt)
			if err != nil {
				t.Fatalf("err %+v", err)
			}
		}(i)
	}

	wg.Wait()
}

func TestSet(t *testing.T) {
	var wg = &sync.WaitGroup{}

	for i, obj := range test.Objs {
		wg.Add(1)
		test.WorkerChan <- struct{}{}

		go func(obj *object.Object, i int) {
			defer func() {
				wg.Done()
				<-test.WorkerChan
			}()

			fmt.Println("obj.Timestamp()", obj.Timestamp())
			var err = test.DB.Set(obj.ID(), obj, map[string]interface{}{
				"another": i % 10,
			})
			if err != nil {
				t.Errorf("err %+v", err)
			}
		}(obj, i)
	}

	wg.Wait()
}

func TestGetMulti(t *testing.T) {
	var (
		ids = []string{
			"some_id_0",
			"some_id_65",
		}
		items, err = test.DB.GetMulti(nil, ids...)
	)

	if err != nil {
		t.Fatalf("%+v", err)
	}

	fmt.Println("items", items)
}

func TestGetBy(t *testing.T) {
	items, err := test.DB.GetBy("another", "=", 0, -1)
	if err != nil {
		t.Fatalf("err %+v:", err)
	}

	fmt.Println("items", items)
}

func TestIter(t *testing.T) {
	iter, err := test.DB.Iterator()
	if err != nil {
		t.Errorf("err %+v", err)
	}

	var (
		item  storage.Item
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

func TestGetChangelogsForObject(t *testing.T) {
	cls, err := test.DB.GetChangelogsForObject("some_id_0")
	if err != nil {
		t.Errorf("err %+v", err)
	}

	fmt.Println("cls", cls)
}
