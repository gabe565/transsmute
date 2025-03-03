package html

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"net/mail"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"mvdan.cc/xurls/v2"
)

func Escape(s string) string {
	return template.HTMLEscapeString(s)
}

func NL2BR(s string) string {
	return strings.ReplaceAll(s, "\n", "<br>\n")
}

func FormatURLs(s string) string {
	urls := xurls.Relaxed().FindAllString(s, -1)
	if len(urls) == 0 {
		return s
	}

	var buf strings.Builder
	buf.Grow(len(s))
	var offset int
	for _, v := range urls {
		idx := strings.Index(s[offset:], v)
		if idx == -1 {
			continue
		}
		buf.WriteString(s[offset : offset+idx])
		offset += idx + len(v)

		if u, err := url.Parse(v); err == nil {
			if u.Scheme == "" {
				if _, err := mail.ParseAddress(v); err == nil && !strings.Contains(v, "/") {
					u.Scheme = "mailto"
					u.OmitHost = true
				} else {
					u.Scheme = "https"
				}
			}

			v = `<a href="` + u.String() + `">` + template.HTMLEscapeString(v) + `</a>`
		}
		buf.WriteString(v)
	}
	buf.WriteString(s[offset:])

	return buf.String()
}

var hrRe = regexp.MustCompile(`(^|\n)(?:---+|___+)(\n|$)`)

func FormatHR(s string, parseHTML bool) string {
	const replace = "$1<hr>"
	switch {
	case !strings.Contains(s, "---") && !strings.Contains(s, "___"):
		return s
	case !parseHTML:
		return hrRe.ReplaceAllString(s, replace)
	default:
		var buf strings.Builder
		buf.Grow(len(s))
		z := html.NewTokenizer(strings.NewReader(s))
		for {
			switch z.Next() {
			case html.ErrorToken:
				if errors.Is(z.Err(), io.EOF) {
					return buf.String()
				}
				return s
			case html.TextToken:
				text := z.Text()
				if bytes.Contains(text, []byte("---")) || bytes.Contains(text, []byte("___")) {
					buf.Write(hrRe.ReplaceAll(text, []byte(replace)))
				} else {
					buf.Write(text)
				}
			default:
				buf.Write(z.Raw())
			}
		}
	}
}
