package html

import (
	"html/template"
	"net/mail"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

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
			if _, err := mail.ParseAddress(v); err == nil && !strings.Contains(v, "/") {
				u.Scheme = "mailto"
				u.OmitHost = true
			} else {
				u.Scheme = "https"
			}

			v = `<a href="` + u.String() + `">` + template.HTMLEscapeString(v) + `</a>`
		}
		buf.WriteString(v)
	}
	buf.WriteString(s[offset:])

	return buf.String()
}

var hashtagRe = regexp.MustCompile("(^|\n| )#([A-Za-z0-9]+)")

func FormatHashtags(s string) string {
	matches := hashtagRe.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return s
	}

	var buf strings.Builder
	buf.Grow(len(s))
	var offset int
	for _, match := range matches {
		v := match[0]
		prefix := match[1]
		slug := match[2]

		idx := strings.Index(s[offset:], v)
		if idx == -1 {
			continue
		}
		buf.WriteString(s[offset : offset+idx])
		offset += idx + len(v)

		u := url.URL{
			Scheme: "https",
			Host:   "youtube.com",
			Path:   path.Join("hashtag", slug),
		}
		v = template.HTMLEscapeString(prefix) + `<a href="` + u.String() + `">#` + template.HTMLEscapeString(slug) + `</a>`
		buf.WriteString(v)
	}
	buf.WriteString(s[offset:])

	return buf.String()
}

var timestampRe = regexp.MustCompile("([0-9]:)?[0-9]+:[0-9]+")

func FormatTimestamps(id, s string) string {
	times := timestampRe.FindAllString(s, -1)
	if len(times) == 0 {
		return s
	}

	var buf strings.Builder
	buf.Grow(len(s))
	var offset int
	for _, v := range times {
		idx := strings.Index(s[offset:], v)
		if idx == -1 {
			continue
		}
		buf.WriteString(s[offset : offset+idx])
		offset += idx + len(v)

		cleaned := v
		if strings.Count(v, ":") == 2 {
			cleaned = strings.Replace(cleaned, ":", "h", 1)
		}
		cleaned = strings.Replace(cleaned, ":", "m", 1)
		cleaned += "s"

		if d, err := time.ParseDuration(cleaned); err == nil {
			u := url.URL{
				Scheme: "https",
				Host:   "youtube.com",
				Path:   "/watch",
				RawQuery: url.Values{
					"v": []string{id},
					"t": []string{strconv.Itoa(int(d.Seconds())) + "s"},
				}.Encode(),
			}

			v = `<a href="` + u.String() + `">` + template.HTMLEscapeString(v) + `</a>`
		}

		buf.WriteString(v)
	}
	buf.WriteString(s[offset:])

	return buf.String()
}
