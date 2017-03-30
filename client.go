package datastore

import (
	cds "cloud.google.com/go/datastore"
	ads "google.golang.org/appengine/datastore"
)

type Client interface {
	NameKey(kind, name string) Key
	IncompleteKey(kind string) Key
	NewQuery(kind string) Query
	Delete(key Key) error
	Get(key Key, dst interface{}) error
	Put(key Key, src interface{}) (Key, error)
	Run(query Query) Iterator
}

type Key interface {
	GetId() int64
	GetName() string

	getInternal() internalKey
}

type Query interface {
	Limit(limit int) Query

	getInternal() internalQuery
}

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
