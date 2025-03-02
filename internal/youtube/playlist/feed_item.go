package playlist

import (
	"net/url"
	"path"
	"time"

	"github.com/gorilla/feeds"
	"google.golang.org/api/youtube/v3"
)

type Item youtube.PlaylistItemSnippet

func (i *Item) FeedItem(embed bool) (*feeds.Item, error) {
	published, err := i.ParsePublishedAt()
	if err != nil {
		return nil, err
	}

	return &feeds.Item{
		Title:   i.Title,
		Link:    &feeds.Link{Href: i.URL().String()},
		Author:  &feeds.Author{Name: i.ChannelTitle},
		Id:      i.ResourceId.VideoId,
		Created: published,
		Content: i.TemplateDescription(embed).String(),
	}, nil
}

func (i *Item) ParsePublishedAt() (time.Time, error) {
	return time.Parse(time.RFC3339, i.PublishedAt)
}

func (i *Item) URL() *url.URL {
	return &url.URL{
		Scheme:   "https",
		Host:     "youtube.com",
		Path:     "/watch",
		RawQuery: url.Values{"v": []string{i.ResourceId.VideoId}}.Encode(),
	}
}

func (i *Item) EmbedURL() *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   "youtube.com",
		Path:   path.Join("embed", i.ResourceId.VideoId),
	}
}
