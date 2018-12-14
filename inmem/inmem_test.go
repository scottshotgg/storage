package inmem_test

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"unsafe"

	"github.com/scottshotgg/storage/inmem"
	"github.com/scottshotgg/storage/store"
	"github.com/scottshotgg/storage/test"
)

func init() {
	test.DB = &inmem.DB{
		Instance: &inmem.NonSyncStore{
			Data: map[string]store.Item{},
		},
	}
}

func TestSize(t *testing.T) {
	TestSet(t)

	for _, data := range test.DB.(*inmem.DB).Instance.Data {
		fmt.Println(unsafe.Sizeof(data.(*store.Object)))
	}
}

func TestSet(t *testing.T) {
	// var wg = &sync.WaitGroup{}

	for _, obj := range test.Objs {
		// wg.Add(1)
		// test.WorkerChan <- struct{}{}

		// go func(obj *store.Object) {
		// defer func() {
		// 	<-test.WorkerChan
		// 	wg.Done()
		// }()

		err := test.DB.Set(obj.ID(), obj)
		if err != nil {
			t.Fatalf("err %+v", err)
		}
		// }(obj)
	}

	// wg.Wait()
}

func TestGet(t *testing.T) {
	var (
		// item  store.Item
		// err   error
		// testt test.Test
		wg = &sync.WaitGroup{}
	)

	for i := 0; i < test.AmountOfTests; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

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
