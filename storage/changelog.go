package storage

import (
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	pb "github.com/scottshotgg/storage/protobufs"
)

type ChangelogIter interface {
	Next() (*Changelog, error)
}

type Changelog struct {
	ID        string
	ObjectID  string
	Type      string
	Timestamp int64
	DBID      string
}

// MarshalBinary implements encoding.BinaryMarshaler
func (c *Changelog) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(&pb.Changelog{
		Id:        c.ID,
		Timestamp: c.Timestamp,
		ItemID:    c.ObjectID,
	})
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (c *Changelog) UnmarshalBinary(data []byte) error {
	var s pb.Changelog

	err := proto.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	c.ID = s.GetId()
	c.ObjectID = s.GetItemID()
	c.Timestamp = s.GetTimestamp()

	return nil
}

func GenTimestamp() int64 {
	return time.Now().Unix()
}

func GenChangelogID() string {
	var v4, err = uuid.NewRandom()

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
