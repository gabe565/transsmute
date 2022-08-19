package feed

import (
	"context"
	"github.com/gorilla/feeds"
	"net/http"
	"path"
	"strings"
)

type OutputFormat uint8

const (
	OutputUnknown OutputFormat = iota
	OutputAtom
	OutputRSS
	OutputJSON
)

type CtxKey uint8

const (
	TypeKey CtxKey = iota
	FeedKey
)

func SetType(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ext := path.Ext(r.URL.Path)
		output := OutputUnknown
		switch ext {
		case ".json":
			output = OutputJSON
		case ".atom", "":
			output = OutputAtom
		case ".rss":
			output = OutputRSS
		}
		r = r.WithContext(context.WithValue(r.Context(), TypeKey, output))
		r.URL.Path = strings.TrimSuffix(r.URL.Path, ext)

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func WriteFeed(w http.ResponseWriter, r *http.Request) (err error) {
	format := r.Context().Value(TypeKey).(OutputFormat)
	feed := r.Context().Value(FeedKey).(*feeds.Feed)

	switch format {
	case OutputAtom:
		if err := feed.WriteAtom(w); err != nil {
			return err
		}
	case OutputJSON:
		if err := feed.WriteJSON(w); err != nil {
			return err
		}
	case OutputRSS:
		if err := feed.WriteRss(w); err != nil {
			return err
		}
	default:
		http.Error(w, "400 invalid format", http.StatusBadRequest)
		return nil
	}

	return nil
}
