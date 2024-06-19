package playlist

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"time"

	"github.com/gorilla/feeds"
	"google.golang.org/api/youtube/v3"
)

func New(ctx context.Context, service *youtube.Service, id string) Playlist {
	return Playlist{
		Service: service,
		Context: ctx,
		ID:      id,
	}
}

type Playlist struct {
	Service *youtube.Service
	Context context.Context
	ID      string
}

var ErrInvalid = errors.New("invalid playlist")

func (p Playlist) Meta() (*youtube.PlaylistSnippet, error) {
	call := p.Service.Playlists.List([]string{"snippet"})
	call.Id(p.ID)
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	if len(resp.Items) < 1 {
		return nil, fmt.Errorf("%s: %w", p.ID, ErrInvalid)
	}

	return resp.Items[0].Snippet, nil
}

func (p Playlist) Feed(disableIframe bool) (*feeds.Feed, error) {
	meta, err := p.Meta()
	if err != nil {
		return nil, err
	}

	u := url.URL{
		Scheme:   "https",
		Host:     "youtube.com",
		Path:     "/playlist",
		RawQuery: url.Values{"list": []string{p.ID}}.Encode(),
	}

	feed := &feeds.Feed{
		Title:       "YouTube - " + meta.Title,
		Link:        &feeds.Link{Href: u.String()},
		Description: meta.Description,
		Created:     time.Now(),
	}

	feed.Items, err = p.FeedItems(disableIframe)
	if err != nil {
		return feed, err
	}

	return feed, nil
}

var ErrLimit = errors.New("exceeded fetch limit")

func (p Playlist) Items() ([]*Item, error) {
	call := p.Service.PlaylistItems.List([]string{"snippet", "status"})
	call.MaxResults(50)
	call.PlaylistId(p.ID)
	limit := 200

	items := make([]*Item, 0)
	i := 0
	err := call.Pages(p.Context, func(response *youtube.PlaylistItemListResponse) error {
		for _, item := range response.Items {
			if item.Status.PrivacyStatus == "private" {
				continue
			}
			items = append(items, (*Item)(item.Snippet))

			i++
			if i >= limit {
				return ErrLimit
			}
		}
		return nil
	})
	if err != nil {
		if !errors.Is(err, ErrLimit) {
			return items, err
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].PublishedAt > items[j].PublishedAt
	})

	return items, nil
}

func (p Playlist) FeedItems(disableIframe bool) ([]*feeds.Item, error) {
	items, err := p.Items()
	if err != nil {
		return nil, err
	}

	feedItems := make([]*feeds.Item, 0, len(items))

	for _, item := range items {
		feedItem, err := item.FeedItem(disableIframe)
		if err != nil {
			return feedItems, err
		}

		feedItems = append(feedItems, feedItem)
	}

	return feedItems, nil
}
