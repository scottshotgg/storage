package query_test

import (
	"context"
	"testing"
	"time"

	"github.com/pizzahutdigital/storage/query"
	"github.com/pizzahutdigital/storage/storage"
)

func TestNew(t *testing.T) {
	var (
		timeout         = 5 * time.Second
		ctx, cancelFunc = context.WithCancel(context.Background())
		output          = &[]storage.Item{}

		q = query.New().
			WithTimeout(timeout).
			WithContext(ctx).
			WithCancelFunc(cancelFunc).
			IsAsync(true).
			WithOutput(output)
	)

	if q.Err() != nil {
		t.Fatalf("err: %+v", q.Err())
	}
}
