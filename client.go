package datastore

import (
	cds "cloud.google.com/go/datastore"
	ads "google.golang.org/appengine/datastore"
)

// Client is a client for reading and writing data in a datastore dataset.
type Client interface {
	IDKey(kind string, id int64, parent Key) Key
	IncompleteKey(kind string, parent Key) Key
	NameKey(kind, name string, parent Key) Key

	Delete(key Key) error
	DeleteMulti(keys []Key) error
	Get(key Key, dst interface{}) error
	GetAll(q Query, dst interface{}) ([]Key, error)
	GetMulti(keys []Key, dst interface{}) error
	Put(key Key, src interface{}) (Key, error)
	PutMulti(keys []Key, src interface{}) ([]Key, error)
	Run(query Query) Iterator

	NewQuery(kind string) Query

	Close() error
}

// Key represents the datastore key for a stored entity.
type Key interface {
	GetID() int64
	GetName() string

	getInternal() internalKey
}

// Query represents a datastore query.
type Query interface {
	Ancestor(ancestor Key) Query
	Distinct() Query
	EventualConsistency() Query
	Filter(filterStr string, value interface{}) Query
	Limit(limit int) Query
	Offset(offset int) Query
	Order(fieldName string) Query

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
