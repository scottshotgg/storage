package storage

type Item interface {
	ID() string
	Value() []byte
	Timestamp() int64
	Keys() map[string]interface{}

	MarshalBinary() (data []byte, err error)
	UnmarshalBinary(data []byte) error
}
