package datastore

import (
	"reflect"

	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	ads "google.golang.org/appengine/datastore"
)

// NewGaeClient creates a datastore client for use in Google App Engine.
func NewGaeClient(ctx context.Context) Client {
	return &gaeClient{context: ctx}
}

type gaeClient struct {
	context context.Context
}

func getGaeKey(key Key) *ads.Key {
	if key == nil {
		return nil
	}
	return key.getInternal().gae
}

func (c *gaeClient) IDKey(kind string, id int64, parent Key) Key {
	return &gaeKey{key: ads.NewKey(c.context, kind, "", id, getGaeKey(parent))}
}

func (c *gaeClient) IncompleteKey(kind string, parent Key) Key {
	return &gaeKey{key: ads.NewIncompleteKey(c.context, kind, getGaeKey(parent))}
}

func (c *gaeClient) NameKey(kind, name string, parent Key) Key {
	return &gaeKey{key: ads.NewKey(c.context, kind, name, 0, getGaeKey(parent))}
}

func (c *gaeClient) Close() error {
	return nil
}

func (c *gaeClient) Delete(key Key) error {
	return ads.Delete(c.context, getGaeKey(key))
}

func convertKeyToGaeKey(keys []Key) []*ads.Key {
	rv := make([]*ads.Key, len(keys))
	for i := range keys {
		rv[i] = getGaeKey(keys[i])
	}
	return rv
}

func convertGaeKeyToKey(keys []*ads.Key) []Key {
	rv := make([]Key, len(keys))
	for i := range keys {
		rv[i] = &gaeKey{key: keys[i]}
	}
	return rv
}

func (c *gaeClient) DeleteMulti(keys []Key) error {
	return ads.DeleteMulti(c.context, convertKeyToGaeKey(keys))
}

func (c *gaeClient) Get(key Key, dst interface{}) error {
	return ads.Get(c.context, getGaeKey(key), dst)
}

func (c *gaeClient) GetAll(q Query, dst interface{}) ([]Key, error) {
	k, err := q.getInternal().gae.GetAll(c.context, dst)
	if err != nil {
		return nil, err
	}
	return convertGaeKeyToKey(k), nil
}

func (c *gaeClient) GetMulti(keys []Key, dst interface{}) error {
	return ads.GetMulti(c.context, convertKeyToGaeKey(keys), dst)
}

func (c *gaeClient) Put(key Key, src interface{}) (Key, error) {
	k, err := ads.Put(c.context, getGaeKey(key), src)
	return &gaeKey{key: k}, err
}

func (c *gaeClient) PutMulti(keys []Key, src interface{}) ([]Key, error) {
	k, err := ads.PutMulti(c.context, convertKeyToGaeKey(keys), src)
	if err != nil {
		return nil, err
	}
	return convertGaeKeyToKey(k), nil
}

func (c *gaeClient) NewQuery(kind string) Query {
	return &gaeQuery{internal: ads.NewQuery(kind)}
}

func (c *gaeClient) Run(query Query) Iterator {
	return &gaeIterator{it: query.getInternal().gae.Run(c.context)}
}

type gaeKey struct {
	key *ads.Key
}

func (c *gaeKey) getInternal() internalKey {
	return internalKey{gae: c.key}
}

func (c *gaeKey) GetID() int64 {
	return c.key.IntID()
}

func (c *gaeKey) GetName() string {
	return c.key.StringID()
}

type gaeQuery struct {
	internal *ads.Query
}

func (q *gaeQuery) Filter(filterStr string, value interface{}) Query {
	if reflect.TypeOf(value) == reflect.TypeOf((*Key)(nil)) {
		value = value.(Key).getInternal()
	}
	return &gaeQuery{internal: q.internal.Filter(filterStr, value)}
}

func (q *gaeQuery) Limit(limit int) Query {
	return &gaeQuery{internal: q.internal.Limit(limit)}
}

func (q *gaeQuery) Offset(offset int) Query {
	return &gaeQuery{internal: q.internal.Offset(offset)}
}

func (q *gaeQuery) Order(fieldName string) Query {
	return &gaeQuery{internal: q.internal.Order(fieldName)}
}

func (q *gaeQuery) getInternal() internalQuery {
	return internalQuery{gae: q.internal}
}

type gaeIterator struct {
	it *ads.Iterator
}

func (i *gaeIterator) Next(dst interface{}) (Key, error) {
	k, err := i.it.Next(dst)
	if err == ads.Done {
		err = iterator.Done
	}
	return &gaeKey{key: k}, err
}
