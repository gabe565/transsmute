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

func (p Channel) Meta(ctx context.Context) (*youtube.Channel, error) {
	call := p.Service.Channels.List([]string{"snippet", "contentDetails"})
	switch {
	case p.Username != "":
		call.ForUsername(p.Username)
	case p.ID != "":
		call.Id(p.ID)
	}
	call.Context(ctx)
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	if len(resp.Items) < 1 {
		return nil, fmt.Errorf("%s: %w", p.ID, ErrInvalid)
	}

	return resp.Items[0], nil
}

func (p Channel) Feed(ctx context.Context) (*feeds.Feed, error) {
	meta, err := p.Meta(ctx)
	if err != nil {
		return nil, err
	}

	pl := playlist.New(p.Service, meta.ContentDetails.RelatedPlaylists.Uploads)
	pl.Iframe = p.Iframe
	return pl.Feed(ctx)
}
