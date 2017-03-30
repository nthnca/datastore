package datastore

import (
	cds "cloud.google.com/go/datastore"
	ads "google.golang.org/appengine/datastore"
)

// Client is a client for reading and writing data in a datastore dataset.
type Client interface {
	NameKey(kind, name string) Key

	IncompleteKey(kind string) Key
	NewQuery(kind string) Query

	Close() error

	Delete(key Key) error
	DeleteMulti(keys []Key) error
	Get(key Key, dst interface{}) error
	GetAll(q Query, dst interface{}) ([]Key, error)
	GetMulti(keys []Key, dst interface{}) error
	Put(key Key, src interface{}) (Key, error)
	PutMulti(keys []Key, src interface{}) ([]Key, error)
	Run(query Query) Iterator
}

// Key represents the datastore key for a stored entity.
type Key interface {
	GetID() int64
	GetName() string

	getInternal() internalKey
}

// Query represents a datastore query.
type Query interface {
	Limit(limit int) Query

	getInternal() internalQuery
}

// Iterator is the result of running a query.
type Iterator interface {
	Next(dst interface{}) (Key, error)
}

type internalQuery struct {
	cloud *cds.Query
	gae   *ads.Query
}

type internalKey struct {
	cloud *cds.Key
	gae   *ads.Key
}
