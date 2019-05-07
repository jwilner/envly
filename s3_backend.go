// +build s3

package main

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	errors "golang.org/x/xerrors"
	"net/url"
	"strings"
)

func init() {
	loaders["s3"] = loadS3
}

func loadS3(s3URI string) ([]string, error) {
	u, err := url.Parse(s3URI)
	if err != nil || u.Scheme != "s3" {
		return nil, errNoMatch
	}

	sess, err := session.NewSession()
	if err != nil {
		return nil, errors.Errorf("aws conn: %w", err)
	}

	obj, err := s3.
		New(sess).
		GetObjectWithContext(
			context.Background(),
			&s3.GetObjectInput{
				Bucket: aws.String(u.Host),
				Key:    aws.String(strings.TrimPrefix(u.Path, "/")),
			},
		)
	if err != nil {
		return nil, errors.Errorf("aws GetObject: %v %v %w", u.Host, u.Path, err)
	}
	defer func() {
		_ = obj.Body.Close()
	}()

	return parse(obj.Body)
}
