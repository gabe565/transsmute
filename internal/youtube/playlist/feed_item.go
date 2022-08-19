package playlist

import (
	"github.com/gabe565/transsmute/internal/youtube/tmpl"
	"github.com/gorilla/feeds"
	"google.golang.org/api/youtube/v3"
	"strings"
	"time"
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

	return &feeds.Item{
		Title:       i.Title,
		Link:        &feeds.Link{Href: "https://youtube.com/watch?v=" + i.ResourceId.VideoId},
		Author:      &feeds.Author{Name: i.ChannelTitle},
		Description: description.String(),
		Id:          i.ResourceId.VideoId,
		Created:     published,
	}, nil
}
