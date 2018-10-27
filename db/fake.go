package db

import (
	"context"
	"encoding/json"
	pkgpath "path"
	"strings"
)

type fakeClient struct {
	db map[string]interface{}
}

// NewFake returns an in-memory implementation of Client.
func NewFake() Client {
	return &fakeClient{
		db: make(map[string]interface{}),
	}
}

func (f *fakeClient) NewRef(path string) Ref {
	return &fakeRef{
		client: f,
		path:   pkgpath.Join("root/", path),
	}
}

type fakeRef struct {
	client *fakeClient
	path   string
}

func (f *fakeRef) Child(path string) Ref {
	return &fakeRef{
		client: f.client,
		path:   pkgpath.Join(f.path, path),
	}
}

func (f *fakeRef) getParent(createMissing bool) (parent map[string]interface{}, childKey string) {
	keys := strings.Split(f.path, "/")
	parent = f.client.db
	for _, key := range keys[:len(keys)-1] {
		child, ok := parent[key].(map[string]interface{})
		if !ok {
			// No such key.
			if createMissing {
				child = make(map[string]interface{})
				parent[key] = child
			} else {
				return nil, ""
			}
		}
		parent = child
	}
	return parent, pkgpath.Base(f.path)
}

func (f *fakeRef) Delete(context.Context) error {
	parent, key := f.getParent(false)
	delete(parent, key)
	return nil
}

func convert(a, b interface{}) error {
	serialized, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return json.Unmarshal(serialized, b)
}

func (f *fakeRef) Get(_ context.Context, v interface{}) error {
	parent, key := f.getParent(false)
	if parent == nil {
		return nil
	}
	value, ok := parent[key]
	if !ok {
		return nil
	}
	return convert(value, v)
}

func (f *fakeRef) Set(_ context.Context, v interface{}) error {
	var generic interface{}
	if err := convert(v, &generic); err != nil {
		return err
	}
	parent, key := f.getParent(true)
	parent[key] = generic
	return nil
}
