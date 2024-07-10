package feed

import (
	"bytes"
	"crypto/sha1" //nolint:gosec
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/eduncan911/podcast"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/feeds"
)

type Format uint8

const (
	FormatUnknown Format = iota
	FormatAtom
	FormatRSS
	FormatJSON
)

func DetectFormat(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ext := path.Ext(r.URL.Path)
		var output Format
		switch ext {
		case ".json":
			output = FormatJSON
		case ".atom":
			output = FormatAtom
		case ".rss":
			output = FormatRSS
		}
		if output != FormatUnknown && ext != "" {
			if ctx := chi.RouteContext(r.Context()); len(ctx.URLParams.Values) != 0 {
				last := len(ctx.URLParams.Values) - 1
				ctx.URLParams.Values[last] = strings.TrimSuffix(ctx.URLParams.Values[last], ext)
			}
		}
		r = r.WithContext(NewFormatContext(r.Context(), output))
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

var (
	ErrContextFormat = errors.New("context format is invalid")
	ErrContextFeed   = errors.New("context feed is invalid")
)

func WriteFeed(w http.ResponseWriter, r *http.Request) error {
	format, ok := FormatFromContext(r.Context())
	if !ok {
		return ErrContextFormat
	}

	feed, ok := FromContext[any](r.Context())
	if !ok {
		return ErrContextFeed
	}

	var buf bytes.Buffer
	hasher := sha1.New() //nolint:gosec
	bufWriter := io.MultiWriter(&buf, hasher)
	switch feed := feed.(type) {
	case *feeds.Feed:
		switch format {
		case FormatAtom, FormatUnknown:
			atomFeed := (&feeds.Atom{Feed: feed}).AtomFeed()
			if feed.Image != nil {
				atomFeed.Icon = feed.Image.Url
			}
			if err := feeds.WriteXML(atomFeed, bufWriter); err != nil {
				return err
			}
			w.Header().Set("Content-Type", "application/rss+xml")
		case FormatJSON:
			jsonFeed := (&feeds.JSON{Feed: feed}).JSONFeed()
			if feed.Image != nil {
				jsonFeed.Icon = feed.Image.Url
			}
			e := json.NewEncoder(bufWriter)
			e.SetIndent("", "  ")
			if err := e.Encode(jsonFeed); err != nil {
				return err
			}
			w.Header().Set("Content-Type", "application/json")
		case FormatRSS:
			if err := feed.WriteRss(bufWriter); err != nil {
				return err
			}
			w.Header().Set("Content-Type", "application/rss+xml")
		default:
			http.Error(w, "400 invalid format", http.StatusBadRequest)
			return nil
		}
	case *podcast.Podcast:
		if err := feed.Encode(bufWriter); err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/xml")
	default:
		panic("invalid feed type")
	}

	etag := `"` + hex.EncodeToString(hasher.Sum(nil)) + `"`
	w.Header().Set("Etag", etag)
	if ifNoneMatch := r.Header.Get("If-None-Match"); ifNoneMatch != "" {
		ifNoneMatch = strings.TrimPrefix(ifNoneMatch, "W/")
		if etag == ifNoneMatch {
			w.WriteHeader(http.StatusNotModified)
			return nil
		}
	}

	if _, err := io.Copy(w, &buf); err != nil {
		return err
	}
	return nil
}
