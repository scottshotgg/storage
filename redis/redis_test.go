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
	"google.golang.org/api/iterator"

	"github.com/scottshotgg/storage/object"
	"github.com/scottshotgg/storage/redis"
	"github.com/scottshotgg/storage/storage"
	"github.com/scottshotgg/storage/test"
)

func init() {
	var (
		test.DB = redis.DB{
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

		_, err = db.Instance.Ping().Result()
	)

	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func TestGet(t *testing.T) {
	var (
		wg sync.WaitGroup

		itemChan = make(chan storage.Item, test.WorkerLimit)
		items    []storage.Item

		doneChan = make(chan struct{})
	)

	go func() {
		for item := range itemChan {
			items = append(items, item)
		}

		doneChan <- struct{}{}
	}()

	for i := 0; i < test.AmountOfTests; i++ {
		wg.Add(1)
		test.WorkerChan <- struct{}{}

		go func(i int) {
			defer func() {
				<-test.WorkerChan
				wg.Done()
			}()

			var item, err = test.DB.Get(context.Background(), fmt.Sprintf("some_id_%d", i))
			if err != nil {
				t.Fatalf("err %+v", err)
				return
			}

			itemChan <- item
		}(i)
	}

	wg.Wait()

	close(itemChan)

	<-doneChan

	fmt.Printf("len %d\n", len(items))
}

func TestGetAll(t *testing.T) {
	var items, err = test.DB.GetAll(context.Background())
	if err != nil {
		fmt.Printf("err %+v\n", err)
	}

	fmt.Println("items len", len(items))
}

func TestSet(t *testing.T) {
	var (
		wg  sync.WaitGroup
		ctx context.Context
	)

	for i, obj := range test.Objs {
		wg.Add(1)
		test.WorkerChan <- struct{}{}

		go func(obj *object.Object, i int) {
			defer func() {
				wg.Done()
				<-test.WorkerChan
			}()

			var err = test.DB.Set(ctx, obj)
			if err != nil {
				t.Errorf("err %+v", err)
			}
		}(obj, i)
	}

	wg.Wait()
}

func TestSetMulti(t *testing.T) {
	var ifaces = make([]storage.Item, len(test.Objs))

	for i := range test.Objs {
		ifaces[i] = test.Objs[i]
	}

	var err = test.DB.SetMulti(context.Background(), ifaces)
	if err != nil {
		fmt.Printf("err %+v\n", err)
	}

	fmt.Println("upload finished")
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
	var items, err = test.DB.GetBy(context.Background(), "another", "=", 0, -1)
	if err != nil {
		t.Fatalf("err %+v:", err)
	}

	fmt.Println("items", items)
}

func TestIter(t *testing.T) {
	var iter, err = test.DB.Iterator()
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
	var cls, err = test.DB.GetChangelogsForObject("some_id_0")
	if err != nil {
		t.Errorf("err %+v", err)
	}

	fmt.Println("cls", cls)
}
