package object

import (
	dstore "cloud.google.com/go/datastore"
	"github.com/golang/protobuf/proto"
	dberrors "github.com/scottshotgg/errors"
	pb "github.com/scottshotgg/storage/protobufs"
	"github.com/scottshotgg/storage/storage"
)

// TODO: Instead of doing this; lets use the protobuf object,
// but make the setters only set if the field is a non-zero value
// Object implements Item
type Object struct {
	id        string
	value     []byte
	timestamp int64
	keys      []string
	// keys      map[string]interface{}
	// disable   bool
}

func New(id string, value []byte, keys []string) *Object {
	return &Object{
		id:        id,
		value:     value,
		timestamp: storage.GenTimestamp(),
		keys:      keys,
	}
}

func (o *Object) ToProto() *pb.Item {
	if o == nil {
		return &pb.Item{}
	}

	return &pb.Item{
		ID:        o.GetID(),
		Value:     o.GetValue(),
		Timestamp: o.GetTimestamp(),
		Keys:      o.GetKeys(),
	}
}

func FromProto(i *pb.Item) *Object {
	return &Object{
		id:        i.GetID(),
		value:     i.GetValue(),
		timestamp: i.GetTimestamp(),
		keys:      i.GetKeys(),
	}
}

func FromResult(res *storage.Result) *Object {
	return &Object{
		id:        res.Item.GetID(),
		value:     res.Item.GetValue(),
		timestamp: res.Item.GetTimestamp(),
		keys:      res.Item.GetKeys(),
	}
}

// this might need to return an error
func FromProps(props dstore.PropertyList) *Object {
	if props == nil {
		return nil
	}

	var (
		item = pb.Item{
			Keys: make([]string, len(props)),
		}

		propMap  = map[string]interface{}{}
		mapValue interface{}
	)

	for _, prop := range props {
		propMap[prop.Name] = prop.Value
	}

	mapValue = propMap["id"]
	if mapValue == nil {
		return nil
	}

	item.ID = mapValue.(string)

	mapValue = propMap["value"]
	if mapValue == nil {
		return nil
	}

	item.Value = mapValue.([]byte)

	mapValue = propMap["timestamp"]
	if mapValue == nil {
		return nil
	}

	item.Timestamp = mapValue.(int64)

	for i := range props {
		item.Keys[i] = props[i].Name
	}

	return FromProto(&item)
}

func (o *Object) SetTimestamp(ts int64) {
	o.timestamp = ts
}

func (o *Object) GetID() string {
	return o.id
}

func (o *Object) GetValue() []byte {
	return o.value
}

func (o *Object) GetTimestamp() int64 {
	return o.timestamp
}

func (o *Object) GetKeys() []string {
	return o.keys
}

func (o *Object) Marshal() ([]byte, error) {
	return proto.Marshal(&pb.Item{
		ID:        o.GetID(),
		Value:     o.GetValue(),
		Timestamp: o.GetTimestamp(),
		Keys:      o.GetKeys(),
	})
}

func (o *Object) Unmarshal(data []byte) error {
	var (
		s   pb.Item
		err = proto.Unmarshal(data, &s)
	)

	if err != nil {
		return err
	}

	o.id = s.GetID()
	o.value = s.GetValue()
	o.timestamp = s.GetTimestamp()
	o.keys = s.GetKeys()

	return nil
}

// MarshalBinary implements encoding.BinaryMarshaler by marshaling to proto form
func (o *Object) MarshalBinary() ([]byte, error) {
	return o.Marshal()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler by unmarshaling from proto form
func (o *Object) UnmarshalBinary(data []byte) error {
	return o.Unmarshal(data)
}

// func (o *Object) GetKeys() map[string]interface{} {
// 	return o.keys
// }

// TODO: these need to be implemented
func (o *Object) GobEncode() ([]byte, error) {
	return nil, dberrors.ErrNotImplemented
}

func (o *Object) GobDecode([]byte) error {
	return dberrors.ErrNotImplemented
}
