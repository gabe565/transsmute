package channel

import (
	"context"
	"errors"
	"fmt"
	"github.com/gabe565/transsmute/internal/youtube/playlist"
	"github.com/gorilla/feeds"
	"google.golang.org/api/youtube/v3"
)

func New(service *youtube.Service, ctx context.Context, id string) Channel {
	return Channel{
		Service: service,
		Context: ctx,
		Id:      id,
	}
}

type Channel struct {
	Service *youtube.Service
	Context context.Context
	Id      string
}

var ErrInvalid = errors.New("invalid channel")

func (p Channel) Meta() (*youtube.Channel, error) {
	call := p.Service.Channels.List([]string{"snippet", "contentDetails"})
	call.Id(p.Id)
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	if len(resp.Items) < 1 {
		return nil, fmt.Errorf("%s: %w", p.Id, ErrInvalid)
	}

	return resp.Items[0], nil
}

func (p Channel) Feed(disableIframe bool) (*feeds.Feed, error) {
	meta, err := p.Meta()
	if err != nil {
		return nil, err
	}

	pl := playlist.New(
		p.Service,
		p.Context,
		meta.ContentDetails.RelatedPlaylists.Uploads,
	)
	return pl.Feed(disableIframe)
}
