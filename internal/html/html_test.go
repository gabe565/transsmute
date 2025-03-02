package html

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscape(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"no change", args{"test"}, "test"},
		{"angled brackets", args{"<test>"}, "&lt;test&gt;"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Escape(tt.args.s))
		})
	}
}

func TestFormatURLs(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"http", args{"prefix http://google.com suffix"}, `prefix <a href="http://google.com">http://google.com</a> suffix`},
		{"https", args{"prefix https://google.com suffix"}, `prefix <a href="https://google.com">https://google.com</a> suffix`},
		{"multiple", args{"prefix https://google.com https://google.com suffix"}, `prefix <a href="https://google.com">https://google.com</a> <a href="https://google.com">https://google.com</a> suffix`},
		{"missing scheme", args{"prefix example.com suffix"}, `prefix <a href="https://example.com">example.com</a> suffix`},
		{"email", args{"prefix example@example.com suffix"}, `prefix <a href="mailto:example@example.com">example@example.com</a> suffix`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, FormatURLs(tt.args.s))
		})
	}
}

func TestNL2BR(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"simple", args{"hello\nworld"}, "hello<br>\nworld"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NL2BR(tt.args.s))
		})
	}
}
