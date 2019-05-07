package main

import (
	errors "golang.org/x/xerrors"
	"os"
	"strings"
)

func init() {
	loaders["file"] = loadFile
}

func loadFile(fileURI string) ([]string, error) {
	if !strings.HasPrefix(fileURI, "file://") {
		return nil, errNoMatch
	}
	f, err := os.Open(fileURI[len("file://"):])
	if err != nil {
		return nil, errors.Errorf("os.Open: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()
	return parse(f)
}
