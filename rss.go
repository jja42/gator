package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jja42/gator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, errors.New("could generate fetch request")
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.New("could not fetch feed")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("unable to read response body")
	}

	var feed RSSFeed

	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, errors.New(("unable to unmarshal xml"))
	}

	return &feed, nil
}

func scrapeFeeds(s *state) error {
	next_feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return errors.New("unable to get next feed")
	}

	current := time.Now()

	fetched_at := sql.NullTime{Time: current, Valid: true}

	s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{ID: next_feed.ID, LastFetchedAt: fetched_at})

	feed, err := fetchFeed(context.Background(), next_feed.Url)
	if err != nil {
		return err
	}

	channel := feed.Channel

	items := channel.Item

	for _, item := range items {
		id := uuid.New()

		created_at := time.Now()
		updated_at := created_at

		layout := time.RFC1123Z
		published_at, err := time.Parse(layout, item.PubDate)
		if err != nil {
			return err
		}

		args := database.CreatePostParams{ID: id, CreatedAt: created_at, UpdatedAt: updated_at, Title: item.Title,
			Url: item.Link, Description: item.Description, PublishedAt: published_at, FeedID: next_feed.ID}

		_, err = s.db.CreatePost(context.Background(), args)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			return errors.New("couldn't create post")
		}
	}

	return nil
}
