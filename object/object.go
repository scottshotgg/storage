package object

import (
	"github.com/golang/protobuf/proto"
	pb "github.com/pizzahutdigital/storage/protobufs"
	"github.com/pizzahutdigital/storage/storage"
)

// Object implements Item
type Object struct {
	id        string
	value     []byte
	disable   bool
	timestamp int64
}

func New(id string, value []byte) *Object {
	return &Object{
		id:        id,
		value:     value,
		timestamp: storage.GenTimestamp(),
	}
}

func (o *Object) ID() string {
	return o.id
}

func (o *Object) Value() []byte {
	return o.value
}

func (o *Object) Timestamp() int64 {
	return o.timestamp
}

// MarshalBinary implements encoding.BinaryMarshaler
func (o *Object) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(&pb.Item{
		Id:        o.ID(),
		Value:     o.Value(),
		Timestamp: o.Timestamp(),
	})
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *Object) UnmarshalBinary(data []byte) error {
	var s pb.Item

	err := proto.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	o.id = s.GetId()
	o.value = s.GetValue()
	o.timestamp = s.GetTimestamp()

	return nil
}
