package docker

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gabe565/transsmute/internal/feed"
	"github.com/go-chi/chi/v5"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
	"github.com/gorilla/feeds"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	repo := chi.URLParam(r, "*")

	reg, err := FindRegistry(repo)
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

	f := &feeds.Feed{
		Title:   "Docker - " + repo,
		Link:    &feeds.Link{Href: reg.GetRepoURL(repo).String()},
		Author:  owner,
		Created: time.Now(),
		Items:   make([]*feeds.Item, 0, len(tags)),
	}

	for _, tag := range tags {
		if err := r.Context().Err(); err != nil {
			panic(err)
		}

		var description strings.Builder
		if err := descriptionTmpl.Execute(&description, DescriptionValues{
			Repo: repo,
			Tag:  tag,
		}); err != nil {
			panic(err)
		}

		item := &feeds.Item{
			Title:       tag,
			Link:        &feeds.Link{Href: reg.GetTagURL(repo, tag).String()},
			Author:      owner,
			Id:          tag,
			Description: description.String(),
		}
		f.Items = append(f.Items, item)
	}

	r = r.WithContext(context.WithValue(r.Context(), feed.FeedKey, f))

	if err := feed.WriteFeed(w, r); err != nil {
		panic(err)
	}
}
