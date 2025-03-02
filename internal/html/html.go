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

	var offset int
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

		newVal := `<a href="` + u.String() + `">` + template.HTMLEscapeString(match) + `</a>`
		s, offset = stringReplaceOffset(s, offset, match, newVal)
	}

	return s
}

var hashtagRe = regexp.MustCompile("(^|\n| )#([A-Za-z0-9]+)")

func FormatHashtags(s string) string {
	matches := hashtagRe.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return s
	}

	var offset int
	for _, match := range matches {
		prefix := match[1]
		slug := match[2]

		u := url.URL{
			Scheme: "https",
			Host:   "youtube.com",
			Path:   path.Join("hashtag", slug),
		}
		newVal := template.HTMLEscapeString(prefix) + `<a href="` + u.String() + `">#` + template.HTMLEscapeString(slug) + `</a>`
		s, offset = stringReplaceOffset(s, offset, match[0], newVal)
	}

	return s
}

var timestampRe = regexp.MustCompile("([0-9]:)?[0-9]+:[0-9]+")

func FormatTimestamps(id, s string) string {
	times := timestampRe.FindAllString(s, -1)
	if len(times) == 0 {
		return s
	}

	var offset int
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

		u := url.URL{
			Scheme: "https",
			Host:   "youtube.com",
			Path:   "/watch",
			RawQuery: url.Values{
				"v": []string{id},
				"t": []string{strconv.Itoa(int(d.Seconds())) + "s"},
			}.Encode(),
		}
		newVal := `<a href="` + u.String() + `">` + template.HTMLEscapeString(match) + `</a>`
		s, offset = stringReplaceOffset(s, offset, match, newVal)
	}

	return s
}

func stringReplaceOffset(s string, offset int, oldVal, newVal string) (string, int) {
	idx := strings.Index(s[offset:], oldVal)
	if idx == -1 {
		return s, offset
	}

	offset += idx
	s = s[:offset] + newVal + s[offset+len(oldVal):]
	return s, offset + len(newVal)
}
