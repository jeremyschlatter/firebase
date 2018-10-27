package firebase

import (
	"context"

	"firebase.google.com/go"
	"google.golang.org/api/option"

	"github.com/jeremyschlatter/firebase/db"
)

type Config = firebase.Config

type real struct {
	app *firebase.App
}

// NewApp mirrors https://godoc.org/firebase.google.com/go#NewApp
//
// It returns a real Firebase client that implements App.
func NewApp(ctx context.Context, config *Config, opts ...option.ClientOption) (App, error) {
	app, err := firebase.NewApp(ctx, config, opts...)
	return real{app}, err
}

func (r real) Database(ctx context.Context) (db.Client, error) {
	client, err := r.app.Database(ctx)
	return db.WrapReal(client), err
}
