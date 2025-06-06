package docker

import (
	"errors"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"gabe565.com/transsmute/internal/feed"
	"github.com/go-chi/chi/v5"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
	"github.com/gorilla/feeds"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func Handler(registries Registries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := chi.URLParam(r, "*")

		var filter *regexp.Regexp
		if v := r.URL.Query().Get("filter"); v != "" {
			v = strings.ReplaceAll(v, " ", "+")
			var err error
			if filter, err = regexp.Compile("^" + v + "$"); err != nil {
				http.Error(w, "Filter regex invalid", http.StatusBadRequest)
				return
			}
		}

		reg, err := registries.Find(repo)
		if err != nil {
			if errors.Is(err, ErrInvalidRegistry) {
				http.Error(w, "404 "+ErrInvalidRegistry.Error(), http.StatusNotFound)
				return
			}
			panic(err)
		}

		auth, err := reg.Authenticator(r.Context(), repo)
		if err != nil {
			panic(err)
		}

		tags, err := crane.ListTags(repo, crane.WithContext(r.Context()), crane.WithAuth(auth))
		if err != nil {
			var transportErr *transport.Error
			switch {
			case errors.As(err, &transportErr):
				msg := http.StatusText(transportErr.StatusCode)
				if len(transportErr.Errors) != 0 {
					msg = transportErr.Errors[0].Message
				}
				http.Error(w, msg, transportErr.StatusCode)
				if transportErr.StatusCode == http.StatusNotFound {
					return
				}
			case errors.Is(err, &name.ErrBadName{}):
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			panic(err)
		}

		slices.Reverse(tags)

		f := &feeds.Feed{
			Title:  "Docker - " + repo,
			Link:   &feeds.Link{Href: reg.GetRepoURL(repo).String()},
			Author: &feeds.Author{Name: reg.GetOwner(repo)},
			Items:  make([]*feeds.Item, 0, len(tags)),
		}

		for _, tag := range tags {
			if filter != nil && !filter.MatchString(tag) {
				continue
			}

			item := &feeds.Item{
				Title: tag,
				Link:  &feeds.Link{Href: reg.GetTagURL(repo, tag).String()},
				Id:    tag,
				Description: g.Group{
					g.Text("Docker tag: "),
					html.Code(g.Text(repo + ":" + tag)),
				}.String(),
			}
			f.Items = append(f.Items, item)
		}

		if err := feed.WriteFeed(w, r, f); err != nil {
			panic(err)
		}
	}
}
