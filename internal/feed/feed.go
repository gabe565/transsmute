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

var ErrContextFormat = errors.New("context format is invalid")

func WriteFeed(w http.ResponseWriter, r *http.Request, f *feeds.Feed) error {
	format, ok := FormatFromContext(r.Context())
	if !ok {
		return ErrContextFormat
	}

	var buf bytes.Buffer
	hasher := sha1.New() //nolint:gosec
	bufWriter := io.MultiWriter(&buf, hasher)
	var lastModified time.Time
	switch format {
	case FormatAtom, FormatUnknown:
		atomFeed := (&feeds.Atom{Feed: f}).AtomFeed()
		if f.Image != nil {
			atomFeed.Icon = f.Image.Url
		}
		if err := feeds.WriteXML(atomFeed, bufWriter); err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/rss+xml")
	case FormatJSON:
		jsonFeed := (&feeds.JSON{Feed: f}).JSONFeed()
		if f.Image != nil {
			jsonFeed.Icon = f.Image.Url
		}
		e := json.NewEncoder(bufWriter)
		e.SetIndent("", "  ")
		if err := e.Encode(jsonFeed); err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
	case FormatRSS:
		if err := f.WriteRss(bufWriter); err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/rss+xml")
	default:
		http.Error(w, "400 invalid format", http.StatusBadRequest)
		return nil
	}
	if !f.Updated.IsZero() {
		lastModified = f.Updated
	} else if !f.Created.IsZero() {
		lastModified = f.Created
	}

	w.Header().Set("ETag", `"`+hex.EncodeToString(hasher.Sum(nil))+`"`)
	http.ServeContent(w, r, "", lastModified, bytes.NewReader(buf.Bytes()))
	return nil
}

func WritePodcast(w http.ResponseWriter, r *http.Request, f *podcast.Podcast) error {
	var buf bytes.Buffer
	hasher := sha1.New() //nolint:gosec
	bufWriter := io.MultiWriter(&buf, hasher)
	var lastModified time.Time
	if err := f.Encode(bufWriter); err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/xml")
	lastModified, _ = time.Parse(time.RFC1123, f.PubDate)

	w.Header().Set("ETag", `"`+hex.EncodeToString(hasher.Sum(nil))+`"`)
	http.ServeContent(w, r, "", lastModified, bytes.NewReader(buf.Bytes()))
	return nil
}
