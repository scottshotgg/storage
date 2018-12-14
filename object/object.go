package object

import (
	dstore "cloud.google.com/go/datastore"
	"github.com/golang/protobuf/proto"
	pb "github.com/scottshotgg/storage/protobufs"
	"github.com/scottshotgg/storage/storage"
)

// Object implements Item
type Object struct {
	id        string
	value     []byte
	timestamp int64
	keys      map[string]interface{}
	// disable   bool
}

func New(id string, value []byte, keys map[string]interface{}) *Object {
	return &Object{
		id:        id,
		value:     value,
		timestamp: storage.GenTimestamp(),
		keys:      keys,
	}
}

func FromProto(i *pb.Item) *Object {
	return &Object{
		id:        i.GetId(),
		value:     i.GetValue(),
		timestamp: i.GetTimestamp(),
		// keys: i.GetKeys(),
	}
}

func FromResult(res *storage.Result) *Object {
	return &Object{
		id:        res.Item.ID(),
		value:     res.Item.Value(),
		timestamp: res.Item.Timestamp(),
		// keys:      res.Item.Keys(),
	}
}

func FromProps(props dstore.PropertyList) *Object {
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
		keys:      propMap,
	}

	// var (
	// 	id        string
	// 	value     []byte
	// 	timestamp int64
	// 	ok        bool
	// )

	// TODO: this is for good measure, but if we control it and provide assurances against that then it is negligible to check this stuff
	// if propMap["id"] == nil {
	// 	// log
	// 	return nil
	// }

	// if propMap["value"] == nil {
	// 	// log
	// 	return nil
	// }

	// if propMap["timestamp"] == nil {
	// 	// log
	// 	return nil
	// }

	// id, ok = propMap["id"].(string)
	// if !ok {
	// 	// log
	// 	return nil
	// }

	// value, ok = propMap["value"].([]byte)
	// if !ok {
	// 	// log
	// 	return nil
	// }

	// timestamp, ok = propMap["timestamp"].(int64)
	// if !ok {
	// 	// log
	// 	return nil
	// }

	// return &Object{
	// 	id:        id,
	// 	value:     value,
	// 	timestamp: timestamp,
	// }
}

func (o *Object) SetTimestamp(ts int64) {
	o.timestamp = ts
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

func (o *Object) Keys() map[string]interface{} {
	return o.keys
}

// MarshalBinary implements encoding.BinaryMarshaler
func (o *Object) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(&pb.Item{
		Id:        o.ID(),
		Value:     o.Value(),
		Timestamp: o.Timestamp(),
		// Keys: // HOW TO DO THIS
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
	// o.keys = s.GetKeys()

	return nil
}
