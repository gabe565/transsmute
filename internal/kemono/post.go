package kemono

import (
	"cmp"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gabe565/transsmute/internal/feed"
	"github.com/gabe565/transsmute/internal/template_funcs"
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

//go:embed post.html.gotmpl
var postTmplStr string

var postTmpl = template.Must(
	template.New("").Funcs(template_funcs.FuncMap()).Parse(postTmplStr),
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
			Title: formatServiceName(creator.Service) + " - " + creator.Name,
			Link:  &feeds.Link{Href: publicURL(host, creator).String()},
			Items: make([]*feeds.Item, 0, 50),
		}
		if creator.Indexed != 0 {
			f.Created = time.Unix(int64(creator.Indexed), 0)
		}
		if creator.Updated != 0 {
			f.Updated = time.Unix(int64(creator.Updated), 0)
		}

		pagesRaw := r.URL.Query().Get("pages")
		pages := uint64(1)
		if pagesRaw != "" {
			if pages, err = strconv.ParseUint(pagesRaw, 10, 64); err != nil || pages == 0 {
				http.Error(w, "pages must be a positive integer", http.StatusBadRequest)
				return
			}
		}

		for page := range pages {
			posts, err := fetchPostPage(r.Context(), postAPIURL(host, creator, page).String())
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

				slices.SortStableFunc(post.Attachments, func(a, b Attachment) int {
					return cmp.Compare(a.Name, b.Name)
				})
				post.Attachments = slices.Compact(post.Attachments)

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
