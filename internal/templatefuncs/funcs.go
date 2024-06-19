package templatefuncs

import (
	"html/template"
	"net/mail"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"

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
		u, err := url.Parse(match)
		if err != nil {
			continue
		}

		if _, err := mail.ParseAddress(match); err == nil && !strings.Contains(match, "/") {
			u.Scheme = "mailto"
			u.OmitHost = true
		} else {
			u.Scheme = "https"
		}

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
	for _, match := range times {
		replaced := match
		if strings.Count(match, ":") == 2 {
			replaced = strings.Replace(replaced, ":", "h", 1)
		}
		replaced = strings.Replace(replaced, ":", "m", 1)
		replaced += "s"

		d, err := time.ParseDuration(replaced)
		if err != nil {
			continue
		}

		if err := timestampTmpl.Execute(&buf, map[string]any{
			"id":      id,
			"seconds": int(d.Seconds()),
			"time":    match,
		}); err != nil {
			continue
		}

		s, offset = stringReplaceOffset(s, offset, match, buf.String())
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
