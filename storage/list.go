package storage

type ItemList interface {
	Each() Item
	Length() int
	Append(Item) int
	Reset()
}
