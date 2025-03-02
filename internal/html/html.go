package html

import (
	"html/template"
	"net/mail"
	"net/url"
	"strings"

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
