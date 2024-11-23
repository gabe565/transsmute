package channel

import (
	"context"
	"errors"
	"fmt"

	"gabe565.com/transsmute/internal/youtube/playlist"
	"github.com/gorilla/feeds"
	"google.golang.org/api/youtube/v3"
)

func New(service *youtube.Service, id, username string) Channel {
	return Channel{
		Service:  service,
		ID:       id,
		Username: username,
	}
}

type Channel struct {
	Service  *youtube.Service
	Username string
	ID       string
	Iframe   bool
}

var ErrInvalid = errors.New("invalid channel")

func (c Channel) Meta(ctx context.Context) (*youtube.Channel, error) {
	call := c.Service.Channels.List([]string{"snippet", "contentDetails"})
	switch {
	case c.Username != "":
		call.ForUsername(c.Username)
	case c.ID != "":
		call.Id(c.ID)
	}
	call.Context(ctx)
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	if len(resp.Items) < 1 {
		return nil, fmt.Errorf("%s: %w", c.ID, ErrInvalid)
	}

	return resp.Items[0], nil
}

func (c Channel) Feed(ctx context.Context) (*feeds.Feed, error) {
	meta, err := c.Meta(ctx)
	if err != nil {
		return nil, err
	}

	p := playlist.New(c.Service, meta.ContentDetails.RelatedPlaylists.Uploads)
	p.Iframe = c.Iframe
	return p.Feed(ctx)
}
