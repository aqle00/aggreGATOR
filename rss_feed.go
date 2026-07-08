package main

import (
	"context"
	"encoding/xml"
	"html"
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
	//1. http request
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	//2. set header
	req.Header.Set("User-Agent", "gator")
	//3. client do
	newClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := newClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//4. read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	rssFeed := &RSSFeed{}
	//5. unmarshal xml to RSSFeed struct
	err = xml.Unmarshal(body, rssFeed)
	if err != nil {
		return nil, err
	}
	//5. clean rssfeed
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	// and clean rss feed contents
	for i := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(rssFeed.Channel.Item[i].Title)
		rssFeed.Channel.Item[i].Description = html.UnescapeString(rssFeed.Channel.Item[i].Description)
	}

	//6. return RSSFeed struct
	return rssFeed, nil
}
