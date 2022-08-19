package template_funcs

import (
	"fmt"
	"html/template"
	"mvdan.cc/xurls/v2"
	"regexp"
	"strconv"
	"strings"
)

func Escape(s string) string {
	return template.HTMLEscapeString(s)
}

func Nl2br(s string) string {
	s = strings.ReplaceAll(s, "\n", "<br>\n")
	return s
}

func FormatUrls(s string) string {
	urls := xurls.Relaxed().FindAllString(s, -1)
	if urls == nil {
		return s
	}

	for _, url := range urls {
		if strings.Contains(url, "@") {
			continue
		}
		s = strings.Replace(
			s,
			url,
			fmt.Sprintf(`<a href="%s">%s</a>`, url, url),
			1,
		)
	}

	return s
}

func FormatHashtags(s string) string {
	re := regexp.MustCompile("(^|\n| )#[A-Za-z0-9]+")
	hashtags := re.FindAllString(s, -1)
	if hashtags == nil {
		return s
	}

	for _, hashtag := range hashtags {
		s = strings.Replace(
			s,
			hashtag,
			fmt.Sprintf(`%c<a href="https://youtube.com/hashtag/%s">%s</a>`, hashtag[0], hashtag[2:], hashtag[1:]),
			1,
		)
	}

	return s
}

func FormatTimestamps(id, s string) string {
	re := regexp.MustCompile("([0-9]:)?[0-9]+:[0-9]+")
	times := re.FindAllString(s, -1)
	if times == nil {
		return s
	}

	for _, time := range times {
		segments := strings.Split(time, ":")

		seconds, err := strconv.Atoi(segments[len(segments)-1])
		if err != nil {
			return s
		}

		min, err := strconv.Atoi(segments[len(segments)-2])
		if err != nil {
			return s
		}
		seconds += min * 60

		if len(segments) == 3 {
			hour, err := strconv.Atoi(segments[len(segments)-3])
			if err != nil {
				panic(err)
			}
			seconds += hour * 60 * 60
		}

		s = strings.Replace(
			s,
			time,
			fmt.Sprintf(`<a href="https://youtube.com/watch?v=%s&t=%ds">%s</a>`, id, seconds, time),
			1,
		)
	}

	return s
}

func Html(s string) template.HTML {
	return template.HTML(s)
}
