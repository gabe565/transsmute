package kemono

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gabe565/transsmute/internal/feed"
	"github.com/gabe565/transsmute/internal/templatefuncs"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/feeds"
)

type Post struct {
	ID          string       `json:"id"`
	User        string       `json:"user"`
	Service     string       `json:"service"`
	Title       string       `json:"title"`
	Content     string       `json:"content"`
	Embed       Embed        `json:"embed"`
	Added       string       `json:"added"`
	Published   string       `json:"published"`
	Edited      string       `json:"edited"`
	Attachments []Attachment `json:"attachments"`
}

type Embed struct {
	URL         string `json:"url"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

type Attachment struct {
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

func (a *Attachment) ThumbURL(host string) *url.URL {
	u := a.URL("img." + host)
	u.Path = path.Join("thumbnail", a.Path)
	u.RawQuery = ""
	return u
}

func (a *Attachment) URL(host string) *url.URL {
	u := &url.URL{
		Scheme:   "https",
		Host:     host,
		Path:     path.Join("data", a.Path),
		RawQuery: url.Values{"f": []string{a.Name}}.Encode(),
	}
	u.RawQuery = strings.ReplaceAll(u.RawQuery, "+", "%20")
	return u
}

//go:embed post.html.gotmpl
var postTmplStr string

//nolint:gochecknoglobals
var postTmpl = template.Must(
	template.New("").Funcs(templatefuncs.FuncMap()).Parse(postTmplStr),
)

func postHandler(host string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creator, err := getCreatorInfo(r.Context(), host, chi.URLParam(r, "creator"), chi.URLParam(r, "service"))
		if err != nil {
			if errors.Is(err, ErrCreatorNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			panic(err)
		}

		f := &feeds.Feed{
			Title:   formatServiceName(creator.Service) + " - " + creator.Name,
			Link:    &feeds.Link{Href: publicURL(host, creator).String()},
			Updated: time.Now(),
			Items:   make([]*feeds.Item, 0, 50),
			Image: &feeds.Image{
				Url:   creator.ImageURL(host).String(),
				Title: creator.Name,
				Link:  publicURL(host, creator).String(),
			},
		}
		if creator.Indexed != 0 {
			f.Created = time.Unix(int64(creator.Updated), 0)
		}

		pagesRaw := r.URL.Query().Get("pages")
		pages := uint64(1)
		if pagesRaw != "" {
			if pages, err = strconv.ParseUint(pagesRaw, 10, 64); err != nil || pages == 0 {
				http.Error(w, "pages must be a positive integer", http.StatusBadRequest)
				return
			}
		}

		query := r.URL.Query().Get("q")

		for page := range pages {
			posts, err := fetchPostPage(r.Context(), postAPIURL(host, creator, page, query).String())
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				panic(err)
			}
			f.Items = slices.Grow(f.Items, len(posts))

			for _, post := range posts {
				item := feeds.Item{
					Id:    post.ID,
					Link:  &feeds.Link{Href: postURL(host, creator, post).String()},
					Title: post.Title,
				}
				if parsed, err := time.Parse("2006-01-02T15:04:05", post.Published); err == nil {
					item.Created = parsed
				}
				if parsed, err := time.Parse("2006-01-02T15:04:05", post.Edited); err == nil {
					item.Updated = parsed
				}

				var buf strings.Builder
				if err := postTmpl.Execute(&buf, map[string]any{
					"Host": host,
					"Post": post,
				}); err != nil {
					panic(err)
				}
				item.Content = buf.String()

				f.Items = append(f.Items, &item)
			}

			if len(posts) < 50 {
				break
			}
		}

		r = r.WithContext(context.WithValue(r.Context(), feed.FeedKey, f))
		if err := feed.WriteFeed(w, r); err != nil {
			panic(err)
		}
	}
}

func fetchPostPage(ctx context.Context, url string) ([]Post, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

	var posts []Post
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		return nil, err
	}

	return posts, nil
}
