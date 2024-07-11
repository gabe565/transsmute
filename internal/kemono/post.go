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
	"strconv"
	"strings"

	"github.com/eduncan911/podcast"
	"github.com/gabe565/transsmute/internal/util"
	"github.com/gorilla/feeds"
)

type Post struct {
	Creator     *Creator      `json:"-"`
	ID          string        `json:"id"`
	User        string        `json:"user"`
	Service     string        `json:"service"`
	Title       string        `json:"title"`
	Content     string        `json:"content"`
	Embed       Embed         `json:"embed"`
	Added       Time          `json:"added"`
	Published   Time          `json:"published"`
	Edited      Time          `json:"edited"`
	Tags        Tags          `json:"tags"`
	File        *Attachment   `json:"file"`
	Attachments []*Attachment `json:"attachments"`
}

func (p *Post) URL() *url.URL {
	u := p.Creator.PublicURL()
	u.Path = path.Join(u.Path, "post", p.ID)
	return u
}

func (p *Post) FeedItem() *feeds.Item {
	item := &feeds.Item{
		Id:      p.ID,
		Link:    &feeds.Link{Href: p.URL().String()},
		Title:   p.Title,
		Created: p.Published.Time,
		Updated: p.Edited.Time,
	}

	var buf strings.Builder
	if err := postTmpl.Execute(&buf, p); err != nil {
		panic(err)
	}
	item.Content = buf.String()
	return item
}

var ErrNoAudio = errors.New("no audio file")

func (p *Post) PodcastItem(ctx context.Context) (*podcast.Item, *Attachment, Time, Time, error) {
	var audio, image *Attachment
	for _, attachment := range p.Attachments {
		switch {
		case audio == nil && attachment.IsAudio():
			audio = attachment
		case image == nil && attachment.IsImage():
			image = attachment
		}
	}
	if audio == nil {
		return nil, nil, Time{}, Time{}, ErrNoAudio
	}

	audioInfo, err := audio.Info(ctx)
	if err != nil {
		return nil, nil, Time{}, Time{}, err
	}

	item := &podcast.Item{
		Title:       p.Title,
		Link:        p.URL().String(),
		Description: p.Content,
		GUID:        p.ID,
		Enclosure: &podcast.Enclosure{
			URL:             audio.URL().String(),
			LengthFormatted: strconv.Itoa(audioInfo.Size),
			TypeFormatted:   audioInfo.MIMEType,
		},
	}
	item.AddPubDate(&p.Published.Time)
	if image != nil {
		item.AddImage(image.ThumbURL().String())
	}
	return item, image, p.Published, p.Edited, nil
}

type Embed struct {
	URL         string `json:"url"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

type Attachment struct {
	post *Post
	Name string `json:"name"`
	Path string `json:"path"`
}

func (a *Attachment) IsImage() bool {
	ext := path.Ext(a.Path)
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}

func (a *Attachment) IsVideo() bool {
	ext := path.Ext(a.Path)
	return ext == ".mp4" || ext == ".webm"
}

func (a *Attachment) IsAudio() bool {
	ext := path.Ext(a.Path)
	return ext == ".mp3" || ext == ".m4a"
}

func (a *Attachment) ThumbURL() *url.URL {
	u := &url.URL{
		Scheme: "https",
		Host:   "img." + a.post.Creator.host,
		Path:   path.Join("thumbnail", "data", a.Path),
	}
	return u
}

func (a *Attachment) URL() *url.URL {
	u := &url.URL{
		Scheme:   "https",
		Host:     a.post.Creator.host,
		Path:     path.Join("data", a.Path),
		RawQuery: url.Values{"f": []string{a.Name}}.Encode(),
	}
	u.RawQuery = strings.ReplaceAll(u.RawQuery, "+", "%20")
	return u
}

type AttachmentInfo struct {
	ID       int    `json:"id"`
	Hash     string `json:"hash"`
	Created  string `json:"ctime"`
	Modified string `json:"mtime"`
	MIMEType string `json:"mime"`
	Ext      string `json:"ext"`
	Added    string `json:"added"`
	Size     int    `json:"size"`
}

func (a *Attachment) InfoAPIURL() *url.URL {
	hash := strings.TrimSuffix(path.Base(a.Path), path.Ext(a.Path))
	return &url.URL{
		Scheme: "https",
		Host:   a.post.Creator.host,
		Path:   path.Join("/api/v1/search_hash", hash),
	}
}

func (a *Attachment) Info(ctx context.Context) (*AttachmentInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.InfoAPIURL().String(), nil)
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

	info := &AttachmentInfo{}
	if err := json.NewDecoder(resp.Body).Decode(info); err != nil {
		return nil, err
	}

	return info, nil
}
