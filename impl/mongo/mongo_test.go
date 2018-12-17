package mongo_test

import (
	"testing"

	"github.com/scottshotgg/storage/impl/mongo"
	"github.com/scottshotgg/storage/test"
)

var (
	db  *mongo.DB
	err error
)

func TestNew(t *testing.T) {
	test.DB, err = mongo.New("mongodb://localhost:27017")
	if err != nil {
		t.Errorf("err %+v", err)
	}
}
