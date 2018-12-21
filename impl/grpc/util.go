package grpc

import (
	pb "github.com/scottshotgg/storage/protobufs"
	"github.com/scottshotgg/storage/storage"
)

// These might be useful in other packages as well

func protoToItems(pbItems []*pb.Item) []storage.Item {
	var items = make([]storage.Item, len(pbItems))

	for i := range pbItems {
		// items[i] = object.FromProto(pbItems[i])
		items[i] = pbItems[i]
	}

	return items
}

func itemsToProto(items []storage.Item) []*pb.Item {
	var pbItems = make([]*pb.Item, len(items))

	for i := range items {
		// pbItems[i] = items[i].ToProto()
		pbItems[i] = items.ToProto()
	}

	return pbItems
}
