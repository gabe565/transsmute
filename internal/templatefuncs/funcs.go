package templatefuncs

import (
	"html/template"
	"net/mail"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"mvdan.cc/xurls/v2"
)

func Escape(s string) string {
	return template.HTMLEscapeString(s)
}

func Nl2br(s string) string {
	s = strings.ReplaceAll(s, "\n", "<br>\n")
	return s
}

//nolint:gochecknoglobals
var linkTmpl = template.Must(
	template.New("").Parse(`<a href="{{ .url }}">{{ .text }}</a>`),
)

func FormatUrls(s string) string {
	urls := xurls.Relaxed().FindAllString(s, -1)
	if urls == nil {
		return s
	}

	var offset int
	var buf strings.Builder
	for _, match := range urls {
		if _, err := mail.ParseAddress(match); err == nil && !strings.Contains(match, "/") {
			continue
		}

		u, err := url.Parse(match)
		if err != nil {
			continue
		}
		u.Scheme = "https"

		if err := linkTmpl.Execute(&buf, map[string]string{
			"url":  u.String(),
			"text": match,
		}); err != nil {
			continue
		}

		s, offset = stringReplaceOffset(s, offset, match, buf.String())
		buf.Reset()
	}

	return s
}

//nolint:gochecknoglobals
var (
	hashtagRe   = regexp.MustCompile("(^|\n| )#[A-Za-z0-9]+")
	hashtagTmpl = template.Must(
		template.New("").Parse(`{{ .prefix }}<a href="https://youtube.com/hashtag/{{ .slug }}">{{ .text }}</a>`),
	)
)

func FormatHashtags(s string) string {
	matches := hashtagRe.FindAllString(s, -1)
	if matches == nil {
		return s
	}

	var offset int
	var buf strings.Builder
	for _, match := range matches {
		prefix := string(match[0])
		slug := match[2:]
		text := match[1:]

		if prefix == "#" {
			prefix = ""
			slug = text
			text = match
		}

		if err := hashtagTmpl.Execute(&buf, map[string]string{
			"prefix": prefix,
			"slug":   path.Clean(slug),
			"text":   text,
		}); err != nil {
			continue
		}

		s, offset = stringReplaceOffset(s, offset, match, buf.String())
		buf.Reset()
	}

	return s
}

//nolint:gochecknoglobals
var (
	timestampRe   = regexp.MustCompile("([0-9]:)?[0-9]+:[0-9]+")
	timestampTmpl = template.Must(
		template.New("").Parse(`<a href="https://youtube.com/watch?v={{ .id }}&t={{ .seconds }}s">{{ .time }}</a>`),
	)
)

func FormatTimestamps(id, s string) string {
	times := timestampRe.FindAllString(s, -1)
	if times == nil {
		return s
	}

	var offset int
	var buf strings.Builder
	for _, time := range times {
		segments := strings.Split(time, ":")

		seconds, err := strconv.Atoi(segments[len(segments)-1])
		if err != nil {
			continue
		}

		min, err := strconv.Atoi(segments[len(segments)-2])
		if err != nil {
			continue
		}
		seconds += min * 60

		if len(segments) == 3 {
			hour, err := strconv.Atoi(segments[len(segments)-3])
			if err != nil {
				continue
			}
			seconds += hour * 60 * 60
		}

		if err := timestampTmpl.Execute(&buf, map[string]any{
			"id":      id,
			"seconds": seconds,
			"time":    time,
		}); err != nil {
			continue
		}

		s, offset = stringReplaceOffset(s, offset, time, buf.String())
		buf.Reset()
	}

	return s
}

//nolint:gosec
func HTML(s string) template.HTML {
	return template.HTML(s)
}

func stringReplaceOffset(s string, offset int, old, new string) (string, int) { //nolint:predeclared
	idx := strings.Index(s[offset:], old)
	if idx == -1 {
		return s, offset
	}

	offset += idx
	s = s[:offset] + new + s[offset+len(old):]
	return s, offset + len(new)
}
