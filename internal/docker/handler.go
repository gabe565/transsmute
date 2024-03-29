package docker

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gabe565/transsmute/internal/feed"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/feeds"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/sirupsen/logrus"
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

	tr, err := reg.Transport(r.Context(), repo)
	if err != nil {
		panic(err)
	}

	hub := registry.Registry{
		URL:    reg.ApiUrl(),
		Client: &http.Client{Transport: tr},
		Logf:   logrus.StandardLogger().Debugf,
	}

	f := &feeds.Feed{
		Title:   "Docker - " + repo,
		Link:    &feeds.Link{Href: reg.GetRepoUrl(repo)},
		Author:  owner,
		Created: time.Now(),
	}

	tags, err := hub.Tags(repo)
	if err != nil {
		var errStatus *registry.HTTPStatusError
		if ok := errors.As(err, &errStatus); ok && errStatus.Response != nil {
			http.Error(w, http.StatusText(errStatus.Response.StatusCode), errStatus.Response.StatusCode)
			return
		}
		panic(err)
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
			Link:        &feeds.Link{Href: reg.GetTagUrl(repo, tag)},
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
