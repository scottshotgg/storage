package reference_test

import (
	"testing"

	"github.com/scottshotgg/storage/impl/reference"
	"github.com/scottshotgg/storage/test"
)

func TestNew(t *testing.T) {
	var err error

	test.DB, err = reference.New("ref", "refstring")
	if err != nil {
		t.Errorf("err %+v", err)
	}
}
