package main

import (
	"bufio"
	errors "golang.org/x/xerrors"
	"io"
	"log"
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
	loaders      = make(map[string]func(string) ([]string, error))
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
	s, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	f, ok := loaders[s.Scheme]
	if !ok {
		return nil, errNoMatch
	}
	return f(uri)
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
