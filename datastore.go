// Package datastore provides a client for Google Cloud Datastore that
// works both on App Engine and in stand alone applications.
//
// This package is a thin wrapper or compatibility layer for
// cloud.google.com/go/datastore and google.golang.org/appengine/datastore.
// The documentation for this project is purposely very light, with a
// recommendation that you normally should just follow the documentation for
// https://godoc.org/google.golang.org/appengine/datastore.
//
// I have only really created wrappers around the functionality I use, but
// overtime I expect this to be a complete wrapper. Pull requests are
// welcomed.
package datastore
