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
	"slices"
	"strconv"
	"time"

	"github.com/eduncan911/podcast"
	"github.com/gabe565/transsmute/internal/util"
	"github.com/gorilla/feeds"
)

type Creator struct {
	host    string
	ID      string `json:"id"`
	Name    string `json:"name"`
	Service string `json:"service"`
	Indexed Time   `json:"indexed"`
	Updated Time   `json:"updated"`
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

func (c *Creator) TagURL(t string) *url.URL {
	u := c.PublicURL()
	q := u.Query()
	q.Set("tag", t)
	u.RawQuery = q.Encode()
	return u
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

func (c *Creator) ProfileAPIURL() *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join("api", "v1", c.Service, "user", c.ID, "profile"),
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
		post.Creator = c
		if len(post.Attachments) == 0 {
			if post.File != nil && post.File.Path != "" {
				post.Attachments = append(post.Attachments, post.File)
			}
		} else {
			seen := make([]string, 0, len(post.Attachments))
			post.Attachments = slices.DeleteFunc(post.Attachments, func(attachment *Attachment) bool {
				if slices.Contains(seen, attachment.Path) {
					return true
				}
				seen = append(seen, attachment.Path)
				return false
			})
		}
		for _, attachment := range post.Attachments {
			attachment.post = post
		}
	}
	return posts, nil
}

var ErrCreatorNotFound = errors.New("creator not found")

func GetCreatorByID(ctx context.Context, host, service, id string) (*Creator, error) {
	creator := &Creator{
		host:    host,
		Service: service,
		ID:      id,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, creator.ProfileAPIURL().String(), nil)
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

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, ErrCreatorNotFound
	default:
		return nil, fmt.Errorf("%w: %s", util.ErrUpstreamRequest, resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&creator); err != nil {
		return nil, err
	}
	return creator, nil
}

func (c *Creator) Feed(ctx context.Context, pages uint64, tag, query string) (*feeds.Feed, error) {
	f := &feeds.Feed{
		Title: formatServiceName(c.Service) + " - " + c.Name,
		Link:  &feeds.Link{Href: c.PublicURL().String()},
		Items: make([]*feeds.Item, 0, 50),
		Image: &feeds.Image{
			Url:   c.ImageURL().String(),
			Title: c.Name,
			Link:  c.PublicURL().String(),
		},
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
			item := post.FeedItem()
			f.Items = append(f.Items, item)
			if f.Updated.IsZero() {
				if !item.Updated.IsZero() {
					f.Updated = item.Updated
				} else if !item.Created.IsZero() {
					f.Updated = item.Created
				}
			}
		}

		if len(posts) < 50 {
			break
		}
	}

	return f, nil
}

func (c *Creator) Podcast(ctx context.Context, pages uint64, tag, query string) (*podcast.Podcast, error) {
	f := podcast.New(c.Name, c.PublicURL().String(), "", nil, nil)
	f.IBlock = "Yes"
	f.TTL = int(24 * time.Hour / time.Second)
	var setPubDate bool

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
			if !setPubDate {
				f.AddPubDate(item.PubDate)
				setPubDate = true
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
