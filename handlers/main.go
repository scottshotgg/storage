package handlers

import (
	"context"

	pb "github.com/scottshotgg/storage/protobufs"
)

// Putting this here because I need something for the package

// TODO: change this to a struct later
type StorageHTTP interface {
	Get(ctx context.Context, in *pb.GetReq) (*pb.GetRes, error)
	GetBy(ctx context.Context, in *pb.GetByReq) (*pb.GetByRes, error)
	GetMulti(ctx context.Context, in *pb.GetMultiReq) (*pb.GetMultiRes, error)
	GetAll(in *pb.GetAllReq, server pb.Storage_GetAllServer) error

	Set(ctx context.Context, in *pb.SetReq) (*pb.SetRes, error)
	SetMulti(ctx context.Context, in *pb.SetMultiReq) (*pb.SetMultiRes, error)

	Delete(ctx context.Context, in *pb.DeleteReq) (*pb.DeleteRes, error)

	Iterator(in *pb.IteratorReq, server pb.Storage_IteratorServer) error
	IteratorBy(in *pb.IteratorByReq, server pb.Storage_IteratorByServer) error

	Audit(ctx context.Context, in *pb.AuditReq) (*pb.AuditRes, error)
	QuickSync(ctx context.Context, in *pb.QuickSyncReq) (*pb.QuickSyncRes, error)
	Sync(ctx context.Context, in *pb.SyncReq) (*pb.SyncRes, error)
}

// TODO: implement later
func New() (StorageHTTP, error) {
	return nil, nil
}
