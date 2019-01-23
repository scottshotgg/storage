package object

import (
	dstore "cloud.google.com/go/datastore"
	"github.com/golang/protobuf/proto"
	dberrors "github.com/scottshotgg/storage/errors"
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
	keys      map[string][]byte
	deleted   bool
	// keys      map[string]interface{}
}

func New(id string, value []byte, keys map[string][]byte) *Object {
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
	// if props == nil {
	// 	return nil
	// }

	// var (
	// 	item = pb.Item{
	// 		Keys: map[string]string{},
	// 	}

	// 	propMap  = map[string]interface{}{}
	// 	mapValue interface{}
	// )

	// for _, prop := range props {
	// 	propMap[prop.Name] = prop.Value
	// }

	// mapValue = propMap["id"]
	// if mapValue == nil {
	// 	return nil
	// }

	// item.ID = mapValue.(string)

	// mapValue = propMap["value"]
	// if mapValue == nil {
	// 	return nil
	// }

	// item.Value = mapValue.([]byte)

	// mapValue = propMap["timestamp"]
	// if mapValue == nil {
	// 	return nil
	// }

	// item.Timestamp = mapValue.(int64)

	// for i := range props {
	// 	item.Keys[i] = props[i].Name
	// }

	// return FromProto(&item)

	if props == nil {
		return nil
	}

	var propMap = map[string]interface{}{}
	for _, prop := range props {
		propMap[prop.Name] = prop.Value
	}

	var (
		id    = propMap["id"].(string)
		value = propMap["value"].([]byte)
		ts    = propMap["timestamp"].(int64)
	)

	delete(propMap, "id")
	delete(propMap, "value")
	delete(propMap, "timestamp")

	return &Object{
		id:        id,
		value:     value,
		timestamp: ts,
		// Fix this later
		// keys:      propMap,
	}
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

func (o *Object) GetKeys() map[string][]byte {
	return o.keys
}

func (o *Object) GetDeleted() bool {
	return o.deleted
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
