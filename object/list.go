package object

import (
	"errors"

	"github.com/pizzahutdigital/storage/storage"
)

type ObjectList struct {
	items   []storage.Item
	counter int
}

func (ol *ObjectList) Append(i storage.Item) {
	ol.items = append(ol.items, i)

	return
}

func (ol *ObjectList) Length() int {
	return len(ol.items)
}

func (ol *ObjectList) Each() (storage.Item, error) {
	if ol.counter == len(ol.items) {
		return nil, errors.New("at the end")
	}

	return ol.items[ol.counter], nil
}

func (ol *ObjectList) Reset() {
	ol.counter = 0
}
