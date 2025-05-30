package playlist

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"time"

	"gabe565.com/transsmute/internal/youtube/middleware"
	"github.com/gorilla/feeds"
	"google.golang.org/api/youtube/v3"
)

func New(service *youtube.Service, id string) Playlist {
	return Playlist{
		Service: service,
		ID:      id,
		Embed:   true,
		Limit:   middleware.DefaultLimit,
	}
}

type Playlist struct {
	Service *youtube.Service
	ID      string
	Embed   bool
	Limit   int
}

var ErrInvalid = errors.New("invalid playlist")

func (p Playlist) Meta(ctx context.Context) (*youtube.PlaylistSnippet, error) {
	call := p.Service.Playlists.List([]string{"snippet"})
	call.Id(p.ID)
	call.Context(ctx)
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	if len(resp.Items) < 1 {
		return nil, fmt.Errorf("%s: %w", p.ID, ErrInvalid)
	}

	return resp.Items[0].Snippet, nil
}

func (p Playlist) Feed(ctx context.Context) (*feeds.Feed, error) {
	meta, err := p.Meta(ctx)
	if err != nil {
		return nil, err
	}

	created, err := time.Parse(time.RFC3339, meta.PublishedAt)
	if err != nil {
		return nil, err
	}

	feed := &feeds.Feed{
		Title:       "YouTube - " + meta.Title,
		Link:        &feeds.Link{Href: p.URL().String()},
		Description: meta.Description,
		Author: &feeds.Author{
			Name: meta.ChannelTitle,
		},
		Created: created,
	}

	feed.Items, err = p.FeedItems(ctx)
	if err != nil {
		return feed, err
	}

	if len(feed.Items) != 0 {
		feed.Updated = feed.Items[0].Created
	}

	return feed, nil
}

var ErrLimit = errors.New("exceeded fetch limit")

func (p Playlist) Items(ctx context.Context) ([]*Item, error) {
	call := p.Service.PlaylistItems.List([]string{"snippet", "status"})
	call.MaxResults(int64(min(p.Limit, 50)))
	call.PlaylistId(p.ID)

	var items []*Item
	i := 0
	err := call.Pages(ctx, func(response *youtube.PlaylistItemListResponse) error {
		items = slices.Grow(items, len(response.Items))
		for _, item := range response.Items {
			if item.Status.PrivacyStatus == "private" {
				continue
			}
			items = append(items, (*Item)(item.Snippet))

			i++
			if i >= p.Limit {
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

	slices.SortStableFunc(items, func(a, b *Item) int {
		return cmp.Compare(b.PublishedAt, a.PublishedAt)
	})

	return items, nil
}

func (p Playlist) FeedItems(ctx context.Context) ([]*feeds.Item, error) {
	items, err := p.Items(ctx)
	if err != nil {
		return nil, err
	}

	feedItems := make([]*feeds.Item, 0, len(items))

	for _, item := range items {
		feedItem, err := item.FeedItem(p.Embed)
		if err != nil {
			return feedItems, err
		}

		feedItems = append(feedItems, feedItem)
	}

	return feedItems, nil
}

func (p Playlist) URL() *url.URL {
	return &url.URL{
		Scheme:   "https",
		Host:     "youtube.com",
		Path:     "/playlist",
		RawQuery: url.Values{"list": []string{p.ID}}.Encode(),
	}
}
