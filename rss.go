package main

import (
	"context"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"time"
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
