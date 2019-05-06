package main

import (
	"reflect"
	"strings"
	"testing"
)

func Test_parseEnv(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    []string
		wantErr bool
	}{
		{
			name:    "allEmpty",
			args:    ``,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "simple",
			args:    `hi=1`,
			want:    []string{"hi=1"},
			wantErr: false,
		},
		{
			name: "multi",
			args: `hi=1
other=2
why=3`,
			want:    []string{"hi=1", "other=2", "why=3"},
			wantErr: false,
		},
		{
			name: "left trims",
			args: `   hi=1
  why=3`,
			want:    []string{"hi=1", "why=3"},
			wantErr: false,
		},
		{
			name: "trailing comments",
			args: `hi=1 # hi
  why=3 		# ok`,
			want:    []string{"hi=1", "why=3"},
			wantErr: false,
		},
		{
			name: "blank lines",
			args: `hi=1

why=3`,
			want:    []string{"hi=1", "why=3"},
			wantErr: false,
		},
		{
			name: "line without equals",
			args: `hi=1
linenoequals
why=3`,
			want:    []string{"hi=1", "why=3"},
			wantErr: false,
		},
		{
			name: "comment line",
			args: `hi=1
# this is a comment with an = sign in it
why=3`,
			want:    []string{"hi=1", "why=3"},
			wantErr: false,
		},
		{
			name: "comment line whitespace",
			args: `hi=1
   # this is a comment with an = sign in it
why=3`,
			want:    []string{"hi=1", "why=3"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(strings.NewReader(tt.args))
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_takeLast(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		out  []string
	}{
		{
			"empty",
			nil,
			nil,
		},
		{
			"no overwrite",
			[]string{"a=1", "b=2"},
			[]string{"a=1", "b=2"},
		},
		{
			"overwrites",
			[]string{"a=1", "b=2", "a=3"},
			[]string{"a=3", "b=2"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if out := takeLast(test.in); !reflect.DeepEqual(out, test.out) {
				t.Fatalf("Got %v but wanted %v", out, test.out)
			}
		})
	}
}
