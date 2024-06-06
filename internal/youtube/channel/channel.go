package channel

import (
	"context"
	"errors"
	"fmt"

	"github.com/gabe565/transsmute/internal/youtube/playlist"
	"github.com/gorilla/feeds"
	"google.golang.org/api/youtube/v3"
)

func New(ctx context.Context, service *youtube.Service, id string) Channel {
	return Channel{
		Service: service,
		Context: ctx,
		ID:      id,
	}
}

type Channel struct {
	Service *youtube.Service
	Context context.Context
	ID      string
}

var ErrInvalid = errors.New("invalid channel")

func (p Channel) Meta() (*youtube.Channel, error) {
	call := p.Service.Channels.List([]string{"snippet", "contentDetails"})
	call.Id(p.ID)
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	if len(resp.Items) < 1 {
		return nil, fmt.Errorf("%s: %w", p.ID, ErrInvalid)
	}

	return resp.Items[0], nil
}

func (p Channel) Feed(disableIframe bool) (*feeds.Feed, error) {
	meta, err := p.Meta()
	if err != nil {
		return nil, err
	}

	pl := playlist.New(
		p.Context,
		p.Service,
		meta.ContentDetails.RelatedPlaylists.Uploads,
	)
	return pl.Feed(disableIframe)
}
