package datastore_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"testing"

	dstore "github.com/pizzahutdigital/datastore"
	"github.com/scottshotgg/storage/impl/datastore"
	"github.com/scottshotgg/storage/object"
	"github.com/scottshotgg/storage/storage"
	"github.com/scottshotgg/storage/test"
	"google.golang.org/api/iterator"
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

	test.DB = &db
}

func TestSet(t *testing.T) {
	var (
		wg  = &sync.WaitGroup{}
		ctx = context.Background()
	)

	for i, obj := range test.Objs {
		wg.Add(1)
		test.WorkerChan <- struct{}{}

		go func(obj *object.Object, i int) {
			defer func() {
				<-test.WorkerChan
				wg.Done()
			}()

			err := test.DB.Set(ctx, obj)
			if err != nil {
				t.Fatalf("err %+v", err)
			}
		}(obj, i)
	}

	wg.Wait()
}

func TestSetMulti(t *testing.T) {
	// For now just do this
	var ifaces = make([]storage.Item, len(test.Objs))
	for i := range test.Objs {
		ifaces[i] = test.Objs[i]
	}

	var err = test.DB.SetMulti(context.Background(), ifaces)
	if err != nil {
		t.Fatalf("err %+v", err)
	}
}

func TestGetBy(t *testing.T) {
	var ctx = context.Background()
	items, err := test.DB.GetBy(ctx, "another", "=", 0, -1)
	if err != nil {
		t.Fatalf("err %+v:", err)
	}

	log.Printf("Items: %+v\nLength: %d\n", items, len(items))
}

func TestGetMulti(t *testing.T) {
	var (
		ids = []string{
			"some_id_0",
			"some_id_65",
		}
		items, err = test.DB.GetMulti(context.Background(), ids...)
	)

	if err != nil {
		t.Fatalf("%+v", err)
	}

	fmt.Println("items", items)
}

// func TestQueryGet(t *testing.T) {
// 	// qf :=

// 	q := query.New(func() {
// 		test.DB.Instance.GetDocument(q.ctx)
// 	})
// 	test.DB.GetQ()
// }

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

			var test test.Test
			err = json.Unmarshal(item.Value(), &test)
			if err != nil {
				t.Fatalf("err %+v", err)
			}
		}(i)
	}

	wg.Wait()
}

func TestDelete(t *testing.T) {
	err := test.DB.Delete("some_id_1")
	if err != nil {
		t.Fatalf("err %+v", err)
	}
}

func TestDeleteAll(t *testing.T) {
	var wg = &sync.WaitGroup{}

	for _, obj := range test.Objs {
		wg.Add(1)
		test.WorkerChan <- struct{}{}

		go func(obj *object.Object) {
			defer func() {
				<-test.WorkerChan
				wg.Done()
			}()

			err := test.DB.Delete(obj.ID())
			if err != nil {
				t.Fatalf("err %+v", err)
			}
		}(obj)
	}

	wg.Wait()
}

func TestIter(t *testing.T) {
	iter, err := test.DB.Iterator()
	if err != nil {
		t.Fatalf("err %+v", err)
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
	fmt.Println(test.DB.GetLatestChangelogForObject("some_id_0"))
}

func TestChangelogsForObject(t *testing.T) {
	fmt.Println(test.DB.GetChangelogsForObject("some_id_0"))
}

func TestIteratorBy(t *testing.T) {
	iter, err := test.DB.IteratorBy("another", "=", 1)
	if err != nil {
		t.Fatalf("err %+v", err)
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

		fmt.Println("item", item, testt)
	}
}

func TestGetAll(t *testing.T) {
	var items, err = test.DB.GetAll(context.Background())
	if err != nil {
		t.Fatalf("err %+v", err)
	}

	fmt.Println("len", len(items))
}
