package twitter

import (
	"context"
	"net/http"
	"time"

	"github.com/gabe565/transsmute/internal/feed"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/feeds"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	scraper := twitterscraper.New()

	profile, err := scraper.GetProfile(username)
	if err != nil {
		panic(err)
	}

	f := &feeds.Feed{
		Title:       "Twitter - @" + profile.Username,
		Link:        &feeds.Link{Href: profile.URL},
		Description: profile.Biography,
		Author:      &feeds.Author{Name: profile.Name},
		Created:     time.Now(),
	}

	for tweet := range scraper.GetTweets(r.Context(), username, 200) {
		if tweet.Error != nil {
			panic(tweet.Error)
		}

		item := &feeds.Item{
			Title:       tweet.Text,
			Link:        &feeds.Link{Href: tweet.PermanentURL},
			Author:      &feeds.Author{Name: profile.Name},
			Description: tweet.HTML,
			Id:          tweet.ID,
			Created:     time.Unix(tweet.Timestamp, 0),
		}

		f.Items = append(f.Items, item)
	}

	r = r.WithContext(context.WithValue(r.Context(), feed.FeedKey, f))

	if err := feed.WriteFeed(w, r); err != nil {
		panic(err)
	}
}
