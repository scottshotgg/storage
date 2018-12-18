package query_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/scottshotgg/storage/query"
	"github.com/scottshotgg/storage/storage"
)

func TestNew(t *testing.T) {
	var (
		of = func(o *query.Operation, v query.Values) (query.Values, error) {
			fmt.Println("hey its me, the op func")
			time.Sleep(1 * time.Second)
			fmt.Println("query", o, v)

			return nil, nil
		}

		operations = []query.OpFunc{of, of, of, of, of}

		timeout         = 4 * time.Second
		ctx, cancelFunc = context.WithCancel(context.Background())
		output          = &[]storage.Item{}

		q = query.New().
			Async(false).
			WithTimeout(timeout).
			WithContext(ctx).
			WithCancelFunc(cancelFunc).
			WithOutput(output).
			WithOperations(operations)
	)

	if q.Err() != nil {
		t.Fatalf("err: %+v", q.Err())
	}

	var err = q.Start()
	if err != nil {
		t.Fatalf("err: %+v", err)
	}

	// err = q.WaitWithTimeout(timeout)
	// if err != nil {
	// 	t.Fatalf("err: %+v", err)
	// }

	q.Wait()

	var results = q.Results()

	if q.Err() != nil {
		t.Fatalf("err: %+v", q.Err())
	}

	log.Printf("results %+v\n", results)
}
