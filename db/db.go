// Package db wraps https://godoc.org/firebase.google.com/go/db with an interface that can be faked in tests.
package db

import "context"

// Client mirrors https://godoc.org/firebase.google.com/go/db#Client
type Client interface {
	NewRef(path string) Ref
}

// Ref mirrors https://godoc.org/firebase.google.com/go/db#Ref
type Ref interface {
	Child(path string) Ref
	Delete(ctx context.Context) error
	Get(ctx context.Context, v interface{}) error
	Set(ctx context.Context, v interface{}) error
}
