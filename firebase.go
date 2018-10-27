// Package firebase wraps the official Firebase client with an interface that can be faked in tests.
//
// This package is designed to be close to a drop-in replacement for firebase.google.org/go, though it
// only implements a very small subset of that interface.
package firebase

import (
	"context"

	"github.com/jeremyschlatter/firebase/db"
)

// App mirrors https://godoc.org/firebase.google.com/go#App
type App interface {
	Database(context.Context) (db.Client, error)
}
