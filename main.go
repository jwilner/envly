package main

import (
	"bufio"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	errors "golang.org/x/xerrors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"unicode"
)

var (
	commentRegex = regexp.MustCompile(`(?:^\s*|\s+)#.*$`)
	errNoMatch   = errors.New("no match")
)

func main() {
	environ := os.Environ()
	if uri, ok := os.LookupEnv("ENVLY_URI"); ok {
		vals, err := load(uri)
		if err != nil {
			log.Fatal(err)
		}
		environ = takeLast(append(environ, vals...))
	}
	if err := run(os.Args[1:], environ); err != nil {
		log.Fatal(err)
	}
}

func load(uri string) ([]string, error) {
	for _, f := range []func(string) ([]string, error){
		loadFile,
		loadS3,
		loadHTTP,
	} {
		env, err := f(uri)
		if err == errNoMatch {
			continue
		}
		if err != nil {
			return nil, err
		}
		return env, nil
	}
	return nil, errNoMatch
}

// takeLast takes the last value added to the provided list according to the name.
// order is preserved.
func takeLast(environ []string) []string {
	positions := make(map[string]int)

	var uniq []string
	for _, e := range environ {
		parts := strings.SplitN(e, "=", 2)
		if i, ok := positions[parts[0]]; ok {
			uniq[i] = e // uniq can never be nil here
		} else {
			positions[parts[0]] = len(uniq)
			uniq = append(uniq, e)
		}
	}

	return uniq
}

func run(argv, env []string) interface{} {
	path, err := exec.LookPath(argv[0])
	if err != nil {
		return errors.Errorf("error finding path: %w", err)
	}
	return syscall.Exec(path, argv, env)
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

func parse(r io.Reader) ([]string, error) {
	var env []string

	s := bufio.NewScanner(r)

	for s.Scan() {
		line := commentRegex.ReplaceAllString(strings.TrimFunc(s.Text(), unicode.IsSpace), "")
		if strings.IndexByte(line, '=') < 0 {
			continue
		}
		env = append(env, line)
	}
	if err := s.Err(); err != nil {
		return nil, errors.Errorf("scanner.Error: %w", err)
	}

	return env, nil
}
