package firebase

import (
	"context"

	"github.com/jeremyschlatter/firebase/db"
)

type fake struct{}

// NewFake returns a fake implementation of App. The fake implementation provides an
// in-memory database.
func NewFake() App {
	return fake{}
}

func (fake) Database(context.Context) (db.Client, error) {
	return db.NewFake(), nil
}
