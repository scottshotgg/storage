package reference_test

import (
	"testing"

	"github.com/scottshotgg/storage/reference"
	"github.com/scottshotgg/storage/test"
)

var (
	db  *sqlite.DB
	err error
)

func TestNew(t *testing.T) {
	db, err = reference.New("ref", "refstring")
	if err != nil {
		t.Errorf("err %+v", err)
	}

	test.DB = &db
}
