package docker

import (
	"context"
	"errors"
	"html"
	"net/http"
	"slices"
	"time"

	"github.com/gabe565/transsmute/internal/feed"
	"github.com/go-chi/chi/v5"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
	"github.com/gorilla/feeds"
)

func Handler(registries Registries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := chi.URLParam(r, "*")

		reg, err := registries.Find(repo)
		if err != nil {
			if errors.Is(err, ErrInvalidRegistry) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			panic(err)
		}

		repo = reg.NormalizeRepo(repo)
		owner := &feeds.Author{Name: reg.GetOwner(repo)}

		auth, err := reg.Authenticator(r.Context(), repo)
		if err != nil {
			panic(err)
		}

		tags, err := crane.ListTags(repo, crane.WithContext(r.Context()), crane.WithAuth(auth))
		if err != nil {
			var transportErr *transport.Error
			if errors.As(err, &transportErr) {
				http.Error(w, http.StatusText(transportErr.StatusCode), transportErr.StatusCode)
				if transportErr.StatusCode == http.StatusNotFound {
					return
				}
			}
			panic(err)
		}

		slices.Reverse(tags)

		f := &feeds.Feed{
			Title:   "Docker - " + repo,
			Link:    &feeds.Link{Href: reg.GetRepoURL(repo).String()},
			Author:  owner,
			Created: time.Now(),
			Items:   make([]*feeds.Item, 0, len(tags)),
		}

		for _, tag := range tags {
			item := &feeds.Item{
				Title:       tag,
				Link:        &feeds.Link{Href: reg.GetTagURL(repo, tag).String()},
				Id:          tag,
				Description: "<p>Docker tag: <code>" + html.EscapeString(repo+":"+tag) + "</code></p>",
			}
			f.Items = append(f.Items, item)
		}

		r = r.WithContext(context.WithValue(r.Context(), feed.FeedKey, f))

		if err := feed.WriteFeed(w, r); err != nil {
			panic(err)
		}
	}
}
