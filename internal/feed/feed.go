package feed

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/feeds"
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
		var output OutputFormat
		switch ext {
		case ".json":
			output = OutputJSON
		case ".atom":
			output = OutputAtom
		case ".rss":
			output = OutputRSS
		}
		r = r.WithContext(context.WithValue(r.Context(), TypeKey, output))
		if output != OutputUnknown {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, ext)
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func WriteFeed(w http.ResponseWriter, r *http.Request) error {
	format := r.Context().Value(TypeKey).(OutputFormat)
	feed := r.Context().Value(FeedKey).(*feeds.Feed)

	var buf bytes.Buffer

	switch format {
	case OutputAtom, OutputUnknown:
		atomFeed := (&feeds.Atom{Feed: feed}).AtomFeed()
		if feed.Image != nil {
			atomFeed.Icon = feed.Image.Url
		}
		if err := feeds.WriteXML(atomFeed, &buf); err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/rss+xml")
	case OutputJSON:
		jsonFeed := (&feeds.JSON{Feed: feed}).JSONFeed()
		if feed.Image != nil {
			jsonFeed.Icon = feed.Image.Url
		}
		e := json.NewEncoder(&buf)
		e.SetIndent("", "  ")
		if err := e.Encode(jsonFeed); err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
	case OutputRSS:
		if err := feed.WriteRss(&buf); err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/rss+xml")
	default:
		http.Error(w, "400 invalid format", http.StatusBadRequest)
		return nil
	}

	if _, err := io.Copy(w, &buf); err != nil {
		return err
	}

	return nil
}
