// +build http

package main

import (
	errors "golang.org/x/xerrors"
	"net/http"
	"net/url"
)

func init() {
	loaders["http"] = loadHTTP
	loaders["https"] = loadHTTP
}

func loadHTTP(httpURI string) ([]string, error) {
	if u, err := url.Parse(httpURI); err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return nil, errNoMatch
	}
	r, err := http.Get(httpURI)
	if err != nil {
		return nil, errors.Errorf("error retrieving URI %v: %w", httpURI, err)
	}
	defer func() {
		_ = r.Body.Close()
	}()
	return parse(r.Body)
}
