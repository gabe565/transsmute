package template_funcs

import (
	"bytes"
	"html/template"
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

var linkTmpl = template.Must(
	template.New("").Parse(`<a href="{{ . }}">{{ . }}</a>`),
)

func FormatUrls(s string) string {
	urls := xurls.Relaxed().FindAllString(s, -1)
	if urls == nil {
		return s
	}

	var buf bytes.Buffer
	for _, url := range urls {
		if strings.Contains(url, "@") {
			continue
		}

		if err := linkTmpl.Execute(&buf, url); err != nil {
			continue
		}

		s = strings.Replace(
			s,
			url,
			buf.String(),
			1,
		)
		buf.Reset()
	}

	return s
}

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

	var buf bytes.Buffer
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
			"slug":   slug,
			"text":   text,
		}); err != nil {
			continue
		}

		s = strings.Replace(
			s,
			match,
			buf.String(),
			1,
		)
		buf.Reset()
	}

	return s
}

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

	var buf bytes.Buffer

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

		s = strings.Replace(
			s,
			time,
			buf.String(),
			1,
		)
		buf.Reset()
	}

	return s
}

func Html(s string) template.HTML {
	return template.HTML(s)
}
