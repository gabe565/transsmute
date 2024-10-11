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
	"time"

	"github.com/eduncan911/podcast"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/feeds"
)

//go:generate go run github.com/dmarkham/enumer -type Format -trimprefix Format -transform lower -text

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
		output, _ := FormatString(strings.TrimPrefix(ext, "."))
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
	var lastModified time.Time
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
		if !feed.Updated.IsZero() {
			lastModified = feed.Updated
		} else if !feed.Created.IsZero() {
			lastModified = feed.Created
		}
	case *podcast.Podcast:
		if err := feed.Encode(bufWriter); err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/xml")
		lastModified, _ = time.Parse(time.RFC1123, feed.PubDate)
	default:
		panic("invalid feed type")
	}

	w.Header().Set("ETag", `"`+hex.EncodeToString(hasher.Sum(nil))+`"`)
	http.ServeContent(w, r, "", lastModified, bytes.NewReader(buf.Bytes()))
	return nil
}
