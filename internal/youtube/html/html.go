package html

import (
	"html/template"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var hashtagRe = regexp.MustCompile(`(^|\s)#([A-Za-z0-9_]+)`)

func FormatHashtags(s string) string {
	return hashtagRe.ReplaceAllString(s, `$1<a href="https://youtube.com/hashtag/$2">#$2</a>`)
}

var timestampRe = regexp.MustCompile("([0-9]:)?[0-9]+:[0-9]+")

func FormatTimestamps(id, s string) string {
	return timestampRe.ReplaceAllStringFunc(s, func(s string) string {
		cleaned := s
		if strings.Count(s, ":") == 2 {
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

			s = `<a href="` + u.String() + `">` + template.HTMLEscapeString(s) + `</a>`
		}
		return s
	})
}
