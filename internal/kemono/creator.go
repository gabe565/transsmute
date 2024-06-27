package kemono

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"slices"
	"strconv"
	"time"

	"github.com/eduncan911/podcast"
	"github.com/gabe565/transsmute/internal/util"
	"github.com/gorilla/feeds"
	"github.com/jellydator/ttlcache/v3"
)

type creatorCacheKey struct {
	host    string
	service string
	creator string
}

type Creator struct {
	host    string
	ID      string `json:"id"`
	Name    string `json:"name"`
	Service string `json:"service"`
	Indexed uint   `json:"indexed"`
	Updated uint   `json:"updated"`
}

func (c *Creator) ImageURL() *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   "img." + c.host,
		Path:   path.Join("icons", c.Service, c.ID),
	}
}

func (c *Creator) PublicURL() *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.Service, "user", c.ID),
	}
}

func (c *Creator) PostAPIURL(page uint64, query string) *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join("api", "v1", c.Service, "user", c.ID),
		RawQuery: url.Values{
			"o": []string{strconv.FormatUint(page*50, 10)},
			"q": []string{query},
		}.Encode(),
	}
}

func (c *Creator) FetchPostPage(ctx context.Context, page uint64, query string) ([]*Post, error) {
	u := c.PostAPIURL(page, query).String()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %s", util.ErrUpstreamRequest, resp.Status)
	}

	var posts []*Post
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		return nil, err
	}

	for _, post := range posts {
		post.creator = c
		seen := make([]string, 0, len(post.Attachments))
		post.Attachments = slices.DeleteFunc(post.Attachments, func(attachment *Attachment) bool {
			if slices.Contains(seen, attachment.Path) {
				return true
			}
			seen = append(seen, attachment.Path)
			return false
		})
		for _, attachment := range post.Attachments {
			attachment.post = post
		}
	}
	return posts, nil
}

var ErrCreatorNotFound = errors.New("creator not found")

func GetCreatorInfo(ctx context.Context, host, service, name string) (*Creator, error) {
	cacheKey := creatorCacheKey{
		host:    host,
		service: service,
		creator: name,
	}
	if cached := creatorCache.Get(cacheKey); cached != nil {
		return cached.Value(), nil
	}

	creator := &Creator{host: host}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	u := url.URL{Scheme: "https", Host: host, Path: "/api/v1/creators"}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		go func() {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %s", util.ErrUpstreamRequest, resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)

	if t, err := decoder.Token(); err != nil {
		return nil, err
	} else if t != json.Delim('[') {
		return nil, &json.UnmarshalTypeError{Value: "object", Type: reflect.TypeOf([]Creator{})}
	}

	for decoder.More() {
		if err := decoder.Decode(&creator); err != nil {
			return nil, err
		}

		if creator.Name == name && creator.Service == service {
			cancel()
			creatorCache.Set(cacheKey, creator, ttlcache.DefaultTTL)
			return creator, nil
		}
	}

	return nil, ErrCreatorNotFound
}

func (c *Creator) Feed(ctx context.Context, pages uint64, tag, query string) (*feeds.Feed, error) {
	f := &feeds.Feed{
		Title:   formatServiceName(c.Service) + " - " + c.Name,
		Link:    &feeds.Link{Href: c.PublicURL().String()},
		Updated: time.Now(),
		Items:   make([]*feeds.Item, 0, 50),
		Image: &feeds.Image{
			Url:   c.ImageURL().String(),
			Title: c.Name,
			Link:  c.PublicURL().String(),
		},
	}
	if c.Indexed != 0 {
		f.Created = time.Unix(int64(c.Updated), 0)
	}

	for page := range pages {
		posts, err := c.FetchPostPage(ctx, page, query)
		if err != nil {
			return nil, err
		}
		f.Items = slices.Grow(f.Items, len(posts))

		for _, post := range posts {
			if tag != "" && !slices.Contains(post.Tags, tag) {
				continue
			}
			f.Items = append(f.Items, post.FeedItem())
		}

		if len(posts) < 50 {
			break
		}
	}

	return f, nil
}

func (c *Creator) Podcast(ctx context.Context, pages uint64, tag, query string) (*podcast.Podcast, error) {
	var pubdate *time.Time
	if c.Indexed != 0 {
		d := time.Unix(int64(c.Updated), 0)
		pubdate = &d
	}
	f := podcast.New(c.Name, c.PublicURL().String(), "", pubdate, nil)
	f.IBlock = "Yes"
	f.TTL = int(24 * time.Hour / time.Second)

	for page := range pages {
		posts, err := c.FetchPostPage(ctx, page, query)
		if err != nil {
			return nil, err
		}
		f.Items = slices.Grow(f.Items, len(posts))

		for _, post := range posts {
			if tag != "" && !slices.Contains(post.Tags, tag) {
				continue
			}

			item, image, err := post.PodcastItem(ctx)
			if err != nil {
				if errors.Is(err, ErrNoAudio) {
					continue
				}
				return nil, err
			}
			f.Items = append(f.Items, item)
			if image != nil && f.Image == nil {
				f.AddImage(image.ThumbURL().String())
			}
		}

		if len(posts) < 50 {
			break
		}
	}

	if f.Image == nil {
		f.AddImage(c.ImageURL().String())
	}

	return &f, nil
}
