package service

import (
	"bandicute-server/internal/storage/repository/post"
	"bandicute-server/pkg/logger"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

const FeedUrlSuffix = "/rss"

type PostParser struct {
	feedParser *gofeed.Parser
}

func NewRssParser() *PostParser {
	return &PostParser{
		feedParser: gofeed.NewParser(),
	}
}

func (p *PostParser) Parse(ctx context.Context, blogUrl string) ([]*post.Model, error) {
	feedUrl := getFeedUrl(blogUrl)
	return p.parseFeed(ctx, feedUrl)
}

func getFeedUrl(blogUrl string) string {
	if strings.HasSuffix(blogUrl, FeedUrlSuffix) {
		return blogUrl
	}
	return strings.TrimSuffix(blogUrl, "/") + FeedUrlSuffix
}

func (p *PostParser) parseFeed(ctx context.Context, rssFeedUrl string) ([]*post.Model, error) {
	feed, err := p.feedParser.ParseURLWithContext(rssFeedUrl, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parseFeed post feed: %w", err)
	}

	var posts []*post.Model
	for _, item := range feed.Items {
		publishedAt, parseErr := parsePublishedAt(item)
		if parseErr != nil {
			logParseError(rssFeedUrl, item, err)
			return nil, fmt.Errorf("failed to parse published at %s: %w", publishedAt, parseErr)
		}

		posts = append(posts, createPost(item, publishedAt))
	}

	logParseSuccess(rssFeedUrl, posts)
	return posts, nil
}

func parsePublishedAt(item *gofeed.Item) (*time.Time, error) {
	if publishedAt := item.PublishedParsed; publishedAt != nil {
		return publishedAt, nil
	}

	return parsePublishedDate(item.Published)
}

// parsePublishedDate attempts to parseFeed a date string in various formats
func parsePublishedDate(date string) (*time.Time, error) {
	formats := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"Mon, 02 Jan 2006 15:04:05 -0700",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, date); err == nil {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("unable to parseFeed date: %s", date)
}

func createPost(item *gofeed.Item, publishedAt *time.Time) *post.Model {
	return &post.Model{
		Title:       item.Title,
		URL:         item.Link,
		Content:     item.Content,
		GUID:        item.GUID,
		PublishedAt: *publishedAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func logParseError(rssFeedUrl string, item *gofeed.Item, err error) {
	logger.Warn("Failed to parseFeed published date", logger.Fields{
		"rssFeedUrl": rssFeedUrl,
		"title":      item.Title,
		"date":       item.Published,
		"error":      err.Error(),
	})
}

func logParseSuccess(rssFeedUrl string, posts []*post.Model) {
	logger.Info("Successfully parsed post", logger.Fields{
		"rssFeedUrl": rssFeedUrl,
		"posts":      len(posts),
	})
}
