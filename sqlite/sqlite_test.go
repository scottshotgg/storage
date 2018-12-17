package sqlite_test

import (
	"testing"

	"github.com/scottshotgg/storage/sqlite"
	"github.com/scottshotgg/storage/test"
)

var (
	db  *sqlite.DB
	err error
)

func TestNew(t *testing.T) {
	db, err = sqlite.New("sqlite3", "user:password@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		t.Errorf("err %+v", err)
	}

	test.DB = db
}
