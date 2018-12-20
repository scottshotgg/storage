package datastore

import (
	"fmt"
	"reflect"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

// DSConfig contains configuration parameters used when constructing an instance of DSInstance
type DSConfig struct {
	Context            context.Context
	ServiceAccountFile string
	ProjectID          string
	Namespace          string
}

// DSInstance contains the configured Datastore client implementation and its targeted namespace
type DSInstance struct {
	client    datastore.Client
	namespace string
}

// NewQuery ...
func (instance *DSInstance) NewQuery(kind string) *datastore.Query {
	return datastore.NewQuery(kind).Namespace(instance.namespace)
}

// NewKey returns a new datastore name key
func (instance *DSInstance) NewKey(kind string, name string, parent *datastore.Key) *datastore.Key {
	key := datastore.NameKey(kind, name, parent)
	key.Namespace = instance.namespace
	return key
}

// Close will close the specified Datastore client contained within DSInstance
func (instance *DSInstance) Close() {
	instance.client.Close()
}

// Count counts the number of items based on a query
func (instance *DSInstance) Count(ctx context.Context, query *datastore.Query) (int, error) {
	return instance.client.Count(ctx, query)
}

// Initialize will populate the specified Datastore client contained within DSInstance and its namespace with the
//   referenced configuration parameters contained within DSConfig
func (instance *DSInstance) Initialize(configuration DSConfig) error {
	client, err := datastore.NewClient(configuration.Context,
		configuration.ProjectID,
		option.WithCredentialsFile(configuration.ServiceAccountFile))
	if err != nil {
		return errors.Wrap(err, "DSInstance.Initialize")
	}

	instance.client = *client
	instance.namespace = configuration.Namespace

	return nil
}

// GetDocumentByKey will populate result by using the provide key
func (instance *DSInstance) GetDocumentByKey(ctx context.Context, key *datastore.Key, result interface{}) error {
	err := instance.client.Get(ctx, key, result)
	if err != nil {
		return errors.Wrap(err, "DSInstance.GetDocumentByKey")
	}
	return nil
}

// GetDocument will populate the referenced result interface from the Datastore with any entity that conforms to the
//   parameters being passed via ctx, kind, and name
func (instance *DSInstance) GetDocument(ctx context.Context, kind string, name string, result interface{}) error {
	// code adapted from MGO code at http://bazaar.launchpad.net/+branch/mgo/v2/view/head:/session.go#L2769

	// check to ensure result argument is an address
	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr {
		panic("result argument must be an address")
	}

	// construct key for Datastore query
	key := &datastore.Key{
		Kind:      kind,
		Name:      name,
		Namespace: instance.namespace,
	}

	// attempt to populate result from Datastore client
	err := instance.client.Get(ctx, key, result)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil
		}

		return errors.Wrap(err, fmt.Sprintf("DSInstance.GetDocument: key: %v", key))
	}

	return nil
}

// GetDocuments will populate the referenced result interface from the Datastore with any entities that conform to the
//   parameters being passed via ctx and query
func (instance *DSInstance) GetDocuments(ctx context.Context, query *datastore.Query, result interface{}) error {
	// code adapted from MGO code at http://bazaar.launchpad.net/+branch/mgo/v2/view/head:/session.go#L2769

	// check to ensure result argument is an address
	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr {
		panic("result argument must be an address")
	}

	// attempt to populate results from Datastore client
	_, err := instance.client.GetAll(ctx, query, result)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil
		}

		return errors.Wrap(err, "DSInstance.GetDocuments")
	}

	return nil
}

// GetKeys will populate the referenced result interface from the Datastore with any entities that conform to the
//   parameters being passed via ctx and query
func (instance *DSInstance) GetKeys(ctx context.Context, query *datastore.Query) ([]*datastore.Key, error) {
	// attempt to populate results from Datastore client
	return instance.client.GetAll(ctx, query, nil)
}

// UpsertDocumentByKey will insert the specified document to the datastore by key
func (instance *DSInstance) UpsertDocumentByKey(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	res, err := instance.client.Put(ctx, key, src)
	if err != nil {
		return nil, errors.Wrap(err, "DSInstance.UpsertDocumentByKey")
	}
	return res, nil
}

// UpsertDocument will insert|update the specified document to the Datastore that conforms to the parameters being
//   passed via ctx, kind, and name
func (instance *DSInstance) UpsertDocument(ctx context.Context, kind string, name string, document interface{}) error {
	// construct key for Datastore transaction
	key := &datastore.Key{
		Kind:      kind,
		Name:      name,
		Namespace: instance.namespace,
	}

	// attempt to insert|update document with Datastore client
	_, err := instance.client.Put(ctx, key, document)
	if err != nil {
		return errors.Wrap(err, "DSInstance.UpsertDocument")
	}

	return nil
}

// UpsertDocuments ...
func (instance *DSInstance) UpsertDocuments(ctx context.Context, keys []*datastore.Key, documents interface{}) error {

	// attempt to insert|update document with Datastore client
	_, err := instance.client.PutMulti(ctx, keys, documents)
	if err != nil {
		return errors.Wrap(err, "DSInstance.UpsertDocuments")
	}

	return nil
}

// DeleteDocument will delete the specified entity from the Datastore that conforms to the parameters being passed via
//   ctx, kind, and name
func (instance *DSInstance) DeleteDocument(ctx context.Context, kind string, name string) error {
	// construct key for Datastore transaction
	key := &datastore.Key{
		Kind:      kind,
		Name:      name,
		Namespace: instance.namespace,
	}

	// attempt to delete specified name|key with Datastore client
	err := instance.client.Delete(ctx, key)
	if err != nil {
		return errors.Wrap(err, "DSInstance.DeleteDocument")
	}

	return nil
}

// DeleteDocuments will delete any specified entities from the Datastore that confirm to the parameters being passed via
//   ctx, kind and names
func (instance *DSInstance) DeleteDocuments(ctx context.Context, kind string, names []string) error {
	var keys []*datastore.Key

	// construct keys for Datastore transaction
	for _, name := range names {
		keys = append(keys, &datastore.Key{
			Kind:      kind,
			Name:      name,
			Namespace: instance.namespace,
		})
	}

	// attempt to delete specified names|keys with Datastore client
	err := instance.client.DeleteMulti(ctx, keys)
	if err != nil {
		return errors.Wrap(err, "DSInstance.DeleteDocuments")
	}

	return nil
}

// DeleteDocumentByKey is a wrapper for the datastore pkg that enables delete a single entity by key (useful for ancestors)
func (instance *DSInstance) DeleteDocumentByKey(ctx context.Context, key *datastore.Key) error {
	err := instance.client.Delete(ctx, key)
	if err != nil {
		return errors.Wrap(err, "DSInstance.DeleteDocumentByKey")
	}
	return nil
}

// DeleteMulti is a wrapper for the datastore pkg that enables you to delete entities by key (useful for ancestors)
func (instance *DSInstance) DeleteMulti(ctx context.Context, keys []*datastore.Key) error {
	err := instance.client.DeleteMulti(ctx, keys)
	if err != nil {
		return errors.Wrap(err, "DSInstance.DeleteMulti")
	}
	return nil
}

// Run will return an iterator for the provided query.
func (instance *DSInstance) Run(ctx context.Context, query *datastore.Query) *datastore.Iterator {
	return instance.client.Run(ctx, query)
}

// GetMultiDocuments will retrieve multiple documents from a kind and an array of ids
func (instance *DSInstance) GetMultiDocuments(ctx context.Context, kind string, ids []string, dst interface{}) error {
	nameKeys := []*datastore.Key{}

	for _, id := range ids {
		nameKeys = append(nameKeys, &datastore.Key{
			Kind:      kind,
			Name:      id,
			Namespace: instance.namespace,
		})
	}

	err := instance.client.GetMulti(ctx, nameKeys, dst)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil
		}

		return errors.Wrap(err, fmt.Sprintf("DSInstance.GetMulti: keys: %v", nameKeys))
	}
	return nil
}

func (instance *DSInstance) Client() *datastore.Client {
	return &instance.client
}

func (instance *DSInstance) Namespace() string {
	return instance.namespace
}
