package html

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatHashtags(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"no change", args{"prefix test suffix"}, "prefix test suffix"},
		{"hashtag", args{"prefix #test suffix"}, `prefix <a href="https://youtube.com/hashtag/test">#test</a> suffix`},
		{"no prefix", args{"prefix #test suffix"}, `prefix <a href="https://youtube.com/hashtag/test">#test</a> suffix`},
		{"multiple", args{"prefix #test #test suffix"}, `prefix <a href="https://youtube.com/hashtag/test">#test</a> <a href="https://youtube.com/hashtag/test">#test</a> suffix`},
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
		{"seconds", args{"dQw4w9WgXcQ", "prefix 0:30 suffix"}, `prefix <a href="https://youtube.com/watch?t=30s&v=dQw4w9WgXcQ">0:30</a> suffix`},
		{"minutes", args{"dQw4w9WgXcQ", "prefix 2:00 suffix"}, `prefix <a href="https://youtube.com/watch?t=120s&v=dQw4w9WgXcQ">2:00</a> suffix`},
		{"hours", args{"dQw4w9WgXcQ", "prefix 2:2:00 suffix"}, `prefix <a href="https://youtube.com/watch?t=7320s&v=dQw4w9WgXcQ">2:2:00</a> suffix`},
		{"multiple", args{"dQw4w9WgXcQ", "prefix 2:00 2:00 suffix"}, `prefix <a href="https://youtube.com/watch?t=120s&v=dQw4w9WgXcQ">2:00</a> <a href="https://youtube.com/watch?t=120s&v=dQw4w9WgXcQ">2:00</a> suffix`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, FormatTimestamps(tt.args.id, tt.args.s))
		})
	}
}
