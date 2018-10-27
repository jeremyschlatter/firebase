package db

import (
	"context"

	"firebase.google.com/go/db"
)

type realClient struct {
	client *db.Client
}

// WrapReal returns a Client that delegates all calls to a real *db.Client.
func WrapReal(client *db.Client) Client {
	return realClient{client}
}

func (r realClient) NewRef(path string) Ref {
	return realRef{r.client.NewRef(path)}
}

type realRef struct {
	ref *db.Ref
}

func (r realRef) Child(path string) Ref {
	return realRef{r.ref.Child(path)}
}

func (r realRef) Delete(ctx context.Context) error {
	return r.ref.Delete(ctx)
}

func (r realRef) Get(ctx context.Context, v interface{}) error {
	return r.ref.Get(ctx, v)
}

func (r realRef) Set(ctx context.Context, v interface{}) error {
	return r.ref.Set(ctx, v)
}
