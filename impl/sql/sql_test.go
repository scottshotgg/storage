package sql_test

import (
	"testing"

	"github.com/scottshotgg/storage/impl/sql"
	"github.com/scottshotgg/storage/test"
)

var (
	db  *sql.DB
	err error
)

func TestNew(t *testing.T) {
	db, err = sql.New("sqlite3", "user:password@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		t.Errorf("err %+v", err)
	}

	test.DB = db
}
