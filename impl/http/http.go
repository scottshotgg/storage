package http

/*
	EXPERIMENTAL: Going to play around with websockets
	https://github.com/tmc/grpc-websocket-proxy
*/

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	dberrors "github.com/scottshotgg/storage/errors"
	"github.com/scottshotgg/storage/object"
	pb "github.com/scottshotgg/storage/protobufs"
	"github.com/scottshotgg/storage/storage"
)

// TODO: Figure out how the streaming will be done

type GRPCConfig struct {
	Address string
	Opts    []grpc.DialOption
}

// DB implements Storage from the storage package using a grpc client
type DB struct {
	Instance pb.StorageClient
	client   *grpc.ClientConn
	config   GRPCConfig
}

// Close should call close internally on your implementations DB and clean up any lose ends; channels, waitgroups, etc
func (db *DB) Close() {
	// myDB.Instance.Close()
}

/*
	grpcAddress,
	grpc.WithInsecure(),
	grpc.WithStatsHandler(&ocgrpc.ClientHandler{
		StartOptions: trace.StartOptions{
			Sampler: trace.AlwaysSample(),
		},
	}))
*/

func New(config GRPCConfig) (*DB, error) {
	var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	// Dial the GRPC connection - important to not that this is a BLOCKING call
	var conn, err = grpc.DialContext(ctx, config.Address, append(config.Opts, grpc.WithBlock())...)
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial")
	}

	var sc = pb.NewStorageClient(conn)

	// Attempt to get the metadata to make sure it is up
	// Consider a `Ping` requirement on the interface

	return &DB{
		Instance: sc,
		client:   conn,
		config:   config,
	}, nil
}

var emptyObject object.Object

func (db *DB) Get(ctx context.Context, id string) (storage.Item, error) {
	// Make a call to the RPC server
	var res, err = db.Instance.Get(ctx, &pb.GetReq{
		ItemID: id,
	})
	if err != nil {
		return nil, err
	}

	// If the item from the response is nil then just return a default empty object as the item
	if res.Item == nil {
		// TODO: need to iron down what the standard response is if there is no object there
		return &emptyObject, nil
	}

	// Decode the proto into the object
	return object.FromProto(res.Item), nil
}

func (db *DB) GetBy(ctx context.Context, key, op string, value interface{}, limit int) ([]storage.Item, error) {
	var (
		// Could make value have a stringer on it...
		// Encode the value; it will be decoded on the other side
		buf bytes.Buffer
		err = gob.NewEncoder(&buf).Encode(value)
		res *pb.GetByRes
	)

	if err != nil {
		return nil, dberrors.ErrEncodingValue
	}

	res, err = db.Instance.GetBy(ctx, &pb.GetByReq{
		Key:   key,
		Op:    op,
		Value: buf.Bytes(),
		Limit: int64(limit),
	})
	if err != nil {
		return nil, err
	}

	// Either way you return err
	// put this back when we start logging
	return protoToItems(res.GetItems()), nil
}

func (db *DB) GetMulti(ctx context.Context, ids []string) ([]storage.Item, error) {
	var res, err = db.Instance.GetMulti(ctx, &pb.GetMultiReq{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}

	// Either way you return err
	// put this back when we start logging
	return protoToItems(res.GetItems()), nil
}

func (db *DB) GetAll(ctx context.Context) ([]storage.Item, error) {
	// TODO: this will be harder since it's a stream; look a ShitStreamer
	return nil, dberrors.ErrNotImplemented
}

func (db *DB) Set(ctx context.Context, item storage.Item) error {
	var _, err = db.Instance.Set(ctx, &pb.SetReq{
		Item: item.ToProto(),
	})

	// Either way you return err
	// put this back when we start logging
	return err
}

func protoToItems(pbItems []*pb.Item) []storage.Item {
	var items = make([]storage.Item, len(pbItems))

	for i := range pbItems {
		items[i] = object.FromProto(pbItems[i])
	}

	return items
}

func itemsToProto(items []storage.Item) []*pb.Item {
	var pbItems = make([]*pb.Item, len(items))

	for i := range items {
		pbItems[i] = items[i].ToProto()
	}

	return pbItems
}

func (db *DB) SetMulti(ctx context.Context, items []storage.Item) error {
	var _, err = db.Instance.SetMulti(ctx, &pb.SetMultiReq{
		Items: itemsToProto(items),
	})

	// Either way you return err
	// put this back when we start logging
	return err
}

func (db *DB) Delete(id string) error {
	var _, err = db.Instance.Delete(context.Background(), &pb.DeleteReq{
		Id: id,
	})

	// Either way you return err
	// put this back when we start logging
	return err
}

func (db *DB) Iterator() (storage.Iter, error) {
	// TODO: this will be harder since it's a stream; look a ShitStreamer
	return nil, dberrors.ErrNotImplemented
}

func (db *DB) IteratorBy(key, op string, value interface{}) (storage.Iter, error) {
	// TODO: this will be harder since it's a stream; look a ShitStreamer
	return nil, dberrors.ErrNotImplemented
}

// TODO: Haven't done any changelog stuff yet for grpc

func (db *DB) GetChangelogsForObject(id string) ([]storage.Changelog, error) {
	return nil, dberrors.ErrNotImplemented
}

func (db *DB) GetLatestChangelogForObject(id string) (*storage.Changelog, error) {
	return nil, dberrors.ErrNotImplemented
}

func (db *DB) DeleteChangelogs(ids ...string) error {
	return dberrors.ErrNotImplemented
}

func (db *DB) ChangelogIterator() (storage.ChangelogIter, error) {
	return nil, dberrors.ErrNotImplemented
}

// storage.Changelog stuff: move this to it's own file

// DeleteBy
// DeleteMulti
// DeleteAll() error

// TODO: Not sure if we want to do all this for the storagerpc
// // Attempt to parse out the status from the error
// var status, ok = status.FromError(err)
// // If we could not parse the error then return an error
// if !ok {
// 	// log
// 	return nil, dberrors.ErrUnableToParseResponse
// }

// // st.Code will now be used as the response status
// var responseStatus = runtime.HTTPStatusFromCode(status.Code())

// // Get the response status from the RPC response.
// // Catch if we got an error back from the RPC server.
// if responseStatus != 200 {

// 	// Assign error number based on the status code that is returned from the RPC service.
// 	switch responseStatus {
// 	case 400:

// 	case 404:

// 	default:

// 	}
// }
