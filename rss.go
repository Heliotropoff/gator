package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
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
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return &RSSFeed{}, nil
	}
	result := RSSFeed{}
	feedData, err := io.ReadAll(response.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	xml.Unmarshal(feedData, &result)
	result.Channel.Title = html.UnescapeString(result.Channel.Title)
	result.Channel.Description = html.UnescapeString(result.Channel.Description)
	for _, value := range result.Channel.Item {
		value.Title = html.UnescapeString(value.Title)
		value.Description = html.UnescapeString(value.Description)
	}
	return &result, nil

}
