package changelog

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/pizzahutdigital/storage"
)

type Changelog struct {
	ID        string
	ObjectID  string
	Type      string
	Timestamp int64
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

func GenInsertChangelog(i storage.Item) *Changelog {
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
