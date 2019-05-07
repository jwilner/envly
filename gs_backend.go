// +build gs

package main

import (
	"cloud.google.com/go/storage"
	"context"
	errors "golang.org/x/xerrors"
	"net/url"
	"strings"
)

func init() {
	loaders["gs"] = loadGS
}

func loadGS(gsURI string) ([]string, error) {
	u, err := url.Parse(gsURI)
	if err != nil || u.Scheme != "gs" {
		return nil, errNoMatch
	}

	client, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, errors.Errorf("storage.NewClient: %w", err)
	}

	r, err := client.
		Bucket(u.Host).
		Object(strings.TrimPrefix(u.Path, "/")).
		NewReader(context.Background())
	if err != nil {
		return nil, errors.Errorf("storage.NewReader for %v %v: %w", u.Host, u.Path, err)
	}
	defer func() {
		_ = r.Close()
	}()

	return parse(r)
}
