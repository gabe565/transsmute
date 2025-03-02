package playlist

import (
	"net/url"
	"time"

	"github.com/gorilla/feeds"
	"google.golang.org/api/youtube/v3"
)

type Item youtube.PlaylistItemSnippet

func (i *Item) FeedItem(embed bool) (*feeds.Item, error) {
	published, err := time.Parse(time.RFC3339, i.PublishedAt)
	if err != nil {
		return nil, err
	}

	u := url.URL{
		Scheme:   "https",
		Host:     "youtube.com",
		Path:     "/watch",
		RawQuery: url.Values{"v": []string{i.ResourceId.VideoId}}.Encode(),
	}

	return &feeds.Item{
		Title:   i.Title,
		Link:    &feeds.Link{Href: u.String()},
		Author:  &feeds.Author{Name: i.ChannelTitle},
		Id:      i.ResourceId.VideoId,
		Created: published,
		Content: i.TemplateDescription(embed).String(),
	}, nil
}
