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

func TestFormatHashtags(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"no change", args{"test"}, "test"},
		{"hashtag", args{" #test"}, ` <a href="https://youtube.com/hashtag/test">#test</a>`},
		{"no prefix", args{"#test"}, `<a href="https://youtube.com/hashtag/test">#test</a>`},
		{"multiple", args{"#test #test"}, `<a href="https://youtube.com/hashtag/test">#test</a> <a href="https://youtube.com/hashtag/test">#test</a>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, FormatHashtags(tt.args.s))
		})
	}
}

func TestFormatTimestamps(t *testing.T) {
	type args struct {
		id string
		s  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"seconds", args{"dQw4w9WgXcQ", "0:30"}, `<a href="https://youtube.com/watch?t=30s&v=dQw4w9WgXcQ">0:30</a>`},
		{"minutes", args{"dQw4w9WgXcQ", "2:00"}, `<a href="https://youtube.com/watch?t=120s&v=dQw4w9WgXcQ">2:00</a>`},
		{"hours", args{"dQw4w9WgXcQ", "2:2:00"}, `<a href="https://youtube.com/watch?t=7320s&v=dQw4w9WgXcQ">2:2:00</a>`},
		{"multiple", args{"dQw4w9WgXcQ", "2:00 2:00"}, `<a href="https://youtube.com/watch?t=120s&v=dQw4w9WgXcQ">2:00</a> <a href="https://youtube.com/watch?t=120s&v=dQw4w9WgXcQ">2:00</a>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, FormatTimestamps(tt.args.id, tt.args.s))
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
		{"simple", args{"https://google.com"}, `<a href="https://google.com">https://google.com</a>`},
		{"multiple", args{"https://google.com https://google.com"}, `<a href="https://google.com">https://google.com</a> <a href="https://google.com">https://google.com</a>`},
		{"missing host", args{"example.com"}, `<a href="https://example.com">example.com</a>`},
		{"email", args{"example@example.com"}, `<a href="mailto:example@example.com">example@example.com</a>`},
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
