package datastore

import (
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	ads "google.golang.org/appengine/datastore"
)

type gaeClient struct {
	context context.Context
}

func NewGaeClient(ctx context.Context) Client {
	return &gaeClient{context: ctx}
}

func (c *gaeClient) NameKey(kind, name string) Key {
	return &gaeKey{key: ads.NewKey(c.context, kind, name, 0, nil)}
}

func (c *gaeClient) IncompleteKey(kind string) Key {
	return &gaeKey{key: ads.NewIncompleteKey(c.context, kind, nil)}
}

func (c *gaeClient) Get(key Key, dst interface{}) error {
	return ads.Get(c.context, key.getInternal().gae, dst)
}

func (c *gaeClient) Put(key Key, src interface{}) (Key, error) {
	k, err := ads.Put(c.context, key.getInternal().gae, src)
	return &gaeKey{key: k}, err
}

func (c *gaeClient) Delete(key Key) error {
	return ads.Delete(c.context, key.getInternal().gae)
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

func (c *gaeKey) GetId() int64 {
	return c.key.IntID()
}

func (c *gaeKey) GetName() string {
	return c.key.StringID()
}

type gaeQuery struct {
	internal *ads.Query
}

func (c *gaeQuery) Limit(limit int) Query {
	return &gaeQuery{internal: c.internal.Limit(limit)}
}

func (c *gaeQuery) getInternal() internalQuery {
	return internalQuery{gae: c.internal}
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
