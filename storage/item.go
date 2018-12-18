package storage

import (
	pb "github.com/scottshotgg/storage/protobufs"
)

type Item interface {
	ID() string
	Value() []byte
	Timestamp() int64
	Keys() map[string]interface{}

	MarshalBinary() (data []byte, err error)
	UnmarshalBinary(data []byte) error

	ToProto() *pb.Item
}
