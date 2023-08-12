package template_funcs

import (
	"testing"
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
			if got := Escape(tt.args.s); got != tt.want {
				t.Errorf("Escape() = %v, want %v", got, tt.want)
			}
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatHashtags(tt.args.s); got != tt.want {
				t.Errorf("FormatHashtags() = %v, want %v", got, tt.want)
			}
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
		{"seconds", args{"dQw4w9WgXcQ", "0:30"}, `<a href="https://youtube.com/watch?v=dQw4w9WgXcQ&t=30s">0:30</a>`},
		{"minutes", args{"dQw4w9WgXcQ", "2:00"}, `<a href="https://youtube.com/watch?v=dQw4w9WgXcQ&t=120s">2:00</a>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatTimestamps(tt.args.id, tt.args.s); got != tt.want {
				t.Errorf("FormatTimestamps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatUrls(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"simple", args{"https://google.com"}, `<a href="https://google.com">https://google.com</a>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatUrls(tt.args.s); got != tt.want {
				t.Errorf("FormatUrls() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNl2br(t *testing.T) {
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
			if got := Nl2br(tt.args.s); got != tt.want {
				t.Errorf("Nl2br() = %v, want %v", got, tt.want)
			}
		})
	}
}
