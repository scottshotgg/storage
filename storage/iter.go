package storage

type Iter interface {
	Next() (Item, error)
}
