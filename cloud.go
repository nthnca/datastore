package datastore

import (
	cds "cloud.google.com/go/datastore"
	"golang.org/x/net/context"
)

// NewCloudClient creates a datastore client for use outside of Google App
// Engine.
func NewCloudClient(projectID string) (Client, error) {
	ctx := context.Background()
	client, err := cds.NewClient(ctx, projectID)
	return &cloudClient{context: ctx, client: client}, err
}

type cloudClient struct {
	context context.Context
	client  *cds.Client
}

func getCloudKey(key Key) *cds.Key {
	if key == nil {
		return nil
	}
	return key.getInternal().cloud
}

func (c *cloudClient) IDKey(kind string, id int64, parent Key) Key {
	return &cloudKey{key: cds.IDKey(kind, id, getCloudKey(parent))}
}

func (c *cloudClient) IncompleteKey(kind string, parent Key) Key {
	return &cloudKey{key: cds.IncompleteKey(kind, getCloudKey(parent))}
}

func (c *cloudClient) NameKey(kind, name string, parent Key) Key {
	return &cloudKey{key: cds.NameKey(kind, name, getCloudKey(parent))}
}

func (c *cloudClient) AllocateIDs(keys []Key) ([]Key, error) {
	k, err := c.client.AllocateIDs(c.context, convertKeyToCloudKey(keys))
	if err != nil {
		return nil, err
	}
	return convertCloudKeyToKey(k), nil
}

func (c *cloudClient) Close() error {
	return c.client.Close()
}

func (c *cloudClient) Delete(key Key) error {
	return c.client.Delete(c.context, getCloudKey(key))
}

func convertKeyToCloudKey(keys []Key) []*cds.Key {
	rv := make([]*cds.Key, len(keys))
	for i := range keys {
		rv[i] = getCloudKey(keys[i])
	}
	return rv
}

func convertCloudKeyToKey(keys []*cds.Key) []Key {
	rv := make([]Key, len(keys))
	for i := range keys {
		rv[i] = &cloudKey{key: keys[i]}
	}
	return rv
}

func (c *cloudClient) DeleteMulti(keys []Key) error {
	return c.client.DeleteMulti(c.context,
		convertKeyToCloudKey(keys))
}

func (c *cloudClient) Get(key Key, dst interface{}) error {
	return c.client.Get(c.context, getCloudKey(key), dst)
}

func (c *cloudClient) GetAll(q Query, dst interface{}) ([]Key, error) {
	k, err := c.client.GetAll(c.context, q.getInternal().cloud, dst)
	if err != nil {
		return nil, err
	}
	return convertCloudKeyToKey(k), nil
}

func (c *cloudClient) GetMulti(keys []Key, dst interface{}) error {
	return c.client.GetMulti(c.context, convertKeyToCloudKey(keys), dst)
}

func (c *cloudClient) Put(key Key, src interface{}) (Key, error) {
	k, err := c.client.Put(c.context, getCloudKey(key), src)
	return &cloudKey{key: k}, err
}

func (c *cloudClient) PutMulti(keys []Key, src interface{}) ([]Key, error) {
	k, err := c.client.PutMulti(c.context, convertKeyToCloudKey(keys), src)
	if err != nil {
		return nil, err
	}
	return convertCloudKeyToKey(k), nil
}

func (c *cloudClient) NewQuery(kind string) Query {
	return &cloudQuery{internal: cds.NewQuery(kind)}
}

func (c *cloudClient) Run(query Query) Iterator {
	return &cloudIterator{it: c.client.Run(
		c.context, query.getInternal().cloud)}
}

type cloudKey struct {
	key *cds.Key
}

func (c *cloudKey) getInternal() internalKey {
	return internalKey{cloud: c.key}
}

func (c *cloudKey) GetID() int64 {
	return c.key.ID
}

func (c *cloudKey) GetName() string {
	return c.key.Name
}

type cloudQuery struct {
	internal *cds.Query
}

func (c *cloudQuery) Limit(limit int) Query {
	return &cloudQuery{internal: c.internal.Limit(limit)}
}

func (c *cloudQuery) getInternal() internalQuery {
	return internalQuery{cloud: c.internal}
}

type cloudIterator struct {
	it *cds.Iterator
}

func (i *cloudIterator) Next(dst interface{}) (Key, error) {
	k, err := i.it.Next(dst)
	return &cloudKey{key: k}, err
}
