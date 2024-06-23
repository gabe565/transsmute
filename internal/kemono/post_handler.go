package kemono

import (
	"context"
	_ "embed"
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gabe565/transsmute/internal/feed"
	"github.com/gabe565/transsmute/internal/templatefuncs"
	"github.com/gabe565/transsmute/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/jellydator/ttlcache/v3"
)

//go:embed post.html.gotmpl
var postTmplStr string

//nolint:gochecknoglobals
var postTmpl = template.Must(
	template.New("").Funcs(templatefuncs.FuncMap()).Parse(postTmplStr),
)

func postHandler(host string) http.HandlerFunc {
	type creatorCacheKey struct {
		service string
		creator string
	}

	creatorCache := ttlcache.New[creatorCacheKey, *Creator](
		ttlcache.WithTTL[creatorCacheKey, *Creator](24*time.Hour),
		ttlcache.WithDisableTouchOnHit[creatorCacheKey, *Creator](),
	)
	go creatorCache.Start()

	return func(w http.ResponseWriter, r *http.Request) {
		cacheKey := creatorCacheKey{
			service: chi.URLParam(r, "service"),
			creator: chi.URLParam(r, "creator"),
		}
		var creator *Creator
		if cached := creatorCache.Get(cacheKey); cached != nil {
			creator = cached.Value()
		} else {
			var err error
			if creator, err = GetCreatorInfo(r.Context(), host, cacheKey.creator, cacheKey.service); err != nil {
				if errors.Is(err, ErrCreatorNotFound) {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				} else if errors.Is(err, util.ErrUpstreamRequest) {
					http.Error(w, err.Error(), http.StatusServiceUnavailable)
					return
				}
				panic(err)
			}
			creatorCache.Set(cacheKey, creator, ttlcache.DefaultTTL)
		}

		pagesRaw := r.URL.Query().Get("pages")
		pages := uint64(1)
		if pagesRaw != "" {
			var err error
			if pages, err = strconv.ParseUint(pagesRaw, 10, 64); err != nil || pages == 0 {
				http.Error(w, "pages must be a positive integer", http.StatusBadRequest)
				return
			}
		}

		query := r.URL.Query().Get("q")

		f, err := creator.Feed(r.Context(), pages, query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			panic(err)
		}

		r = r.WithContext(context.WithValue(r.Context(), feed.FeedKey, f))
		if err := feed.WriteFeed(w, r); err != nil {
			panic(err)
		}
	}
}
