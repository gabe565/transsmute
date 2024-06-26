package kemono

import (
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/gorilla/feeds"
)

type Post struct {
	creator     *Creator
	ID          string        `json:"id"`
	User        string        `json:"user"`
	Service     string        `json:"service"`
	Title       string        `json:"title"`
	Content     string        `json:"content"`
	Embed       Embed         `json:"embed"`
	Added       string        `json:"added"`
	Published   string        `json:"published"`
	Edited      string        `json:"edited"`
	Tags        Tags          `json:"tags"`
	Attachments []*Attachment `json:"attachments"`
}

func (p *Post) FeedItem() *feeds.Item {
	item := &feeds.Item{
		Id:    p.ID,
		Link:  &feeds.Link{Href: p.creator.PostURL(p).String()},
		Title: p.Title,
	}
	if parsed, err := time.Parse("2006-01-02T15:04:05", p.Published); err == nil {
		item.Created = parsed
	}
	if parsed, err := time.Parse("2006-01-02T15:04:05", p.Edited); err == nil {
		item.Updated = parsed
	}

	var buf strings.Builder
	if err := postTmpl.Execute(&buf, p); err != nil {
		panic(err)
	}
	item.Content = buf.String()
	return item
}

type Embed struct {
	URL         string `json:"url"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

type Attachment struct {
	post *Post
	Name string `json:"name"`
	Path string `json:"path"`
}

func (a *Attachment) IsImage() bool {
	ext := path.Ext(a.Path)
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}

func (a *Attachment) IsVideo() bool {
	ext := path.Ext(a.Path)
	return ext == ".mp4" || ext == ".webm"
}

func (a *Attachment) ThumbURL() *url.URL {
	u := &url.URL{
		Scheme: "https",
		Host:   "img." + a.post.creator.host,
		Path:   path.Join("thumbnail", "data", a.Path),
	}
	return u
}

func (a *Attachment) URL() *url.URL {
	u := &url.URL{
		Scheme:   "https",
		Host:     a.post.creator.host,
		Path:     path.Join("data", a.Path),
		RawQuery: url.Values{"f": []string{a.Name}}.Encode(),
	}
	u.RawQuery = strings.ReplaceAll(u.RawQuery, "+", "%20")
	return u
}
