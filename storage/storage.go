package storage

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

type Storage interface {
	// Name() string
	// Type() storeType
	Get(ctx context.Context, id string) (Item, error)
	Set(id string, i Item) error
	Delete(id string) error
	Iterator() (Iter, error)
	ChangelogIterator() (ChangelogIter, error)
	GetLatestChangelogForObject(id string) (*Changelog, error)
}

type Item interface {
	ID() string
	Value() []byte
	MarshalBinary() (data []byte, err error)
	UnmarshalBinary(data []byte) error
	Timestamp() int64
}

type ChangelogIter interface {
	Next() (*Changelog, error)
}

type Iter interface {
	Next() (Item, error)
}

type Changelog struct {
	ID        string
	ObjectID  string
	Type      string
	Timestamp int64
}

type Result struct {
	Item Item
	Err  error
}

func GenTimestamp() int64 {
	return time.Now().Unix()
}

func GenChangelogID() string {
	v4, err := uuid.NewRandom()

	for err != nil {
		log.Println("Could not gen uuid, trying again...")

		v4, err = uuid.NewRandom()
	}

	return v4.String()
}

func GenInsertChangelog(i Item) *Changelog {
	return &Changelog{
		ID:        i.ID() + "-" + GenChangelogID(),
		ObjectID:  i.ID(),
		Timestamp: i.Timestamp(),
	}
}

func GenDeleteChangelog(id string) *Changelog {
	return &Changelog{
		ID:        id + "-" + GenChangelogID(),
		ObjectID:  id,
		Timestamp: GenTimestamp(),
	}
}
