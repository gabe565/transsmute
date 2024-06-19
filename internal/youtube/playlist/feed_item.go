package playlist

import (
	"net/url"
	"strings"
	"time"

	"github.com/gabe565/transsmute/internal/youtube/tmpl"
	"github.com/gorilla/feeds"
	"google.golang.org/api/youtube/v3"
)

type Item youtube.PlaylistItemSnippet

func (i Item) FeedItem(disableIframe bool) (*feeds.Item, error) {
	published, err := time.Parse(time.RFC3339, i.PublishedAt)
	if err != nil {
		return nil, err
	}

	var description strings.Builder
	if err := tmpl.DescriptionTmpl.Execute(&description, map[string]any{
		"Item":          i,
		"DisableIframe": disableIframe,
	}); err != nil {
		return nil, err
	}

	u := url.URL{
		Scheme:   "https",
		Host:     "youtube.com",
		Path:     "/watch",
		RawQuery: url.Values{"v": []string{i.ResourceId.VideoId}}.Encode(),
	}

	return &feeds.Item{
		Title:       i.Title,
		Link:        &feeds.Link{Href: u.String()},
		Author:      &feeds.Author{Name: i.ChannelTitle},
		Description: description.String(),
		Id:          i.ResourceId.VideoId,
		Created:     published,
	}, nil
}
