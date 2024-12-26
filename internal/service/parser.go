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

type Parser interface {
	ParseBlog(ctx context.Context, url string) ([]*post.Model, error)
}

type TistoryParser struct {
	parser *gofeed.Parser
}

func NewTistoryParser() Parser {
	return &TistoryParser{
		parser: gofeed.NewParser(),
	}
}

func (p *TistoryParser) ParseBlog(ctx context.Context, url string) ([]*post.Model, error) {
	// Add /rss to the Repository if it's not already there
	if !strings.HasSuffix(url, "/rss") {
		url = strings.TrimSuffix(url, "/") + "/rss"
	}

	logger.Info("Parsing post", logger.Fields{
		"url": url,
	})

	feed, err := p.parser.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse post feed: %w", err)
	}

	var posts []*post.Model
	for _, item := range feed.Items {
		// Parse the published date
		publishedAt := item.PublishedParsed
		if publishedAt == nil {
			// Try to parse the published date manually
			publishedAt, err = parseDate(item.Published)
			if err != nil {
				logger.Warn("Failed to parse published date", logger.Fields{
					"url":   url,
					"title": item.Title,
					"date":  item.Published,
					"error": err.Error(),
				})
				// Use current time as fallback
				now := time.Now()
				publishedAt = &now
			}
		}

		post := &post.Model{
			Title:       item.Title,
			URL:         item.Link,
			Content:     item.Content,
			GUID:        item.GUID,
			PublishedAt: *publishedAt,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		posts = append(posts, post)
	}

	logger.Info("Successfully parsed post", logger.Fields{
		"url":   url,
		"posts": len(posts),
	})

	return posts, nil
}

// parseDate attempts to parse a date string in various formats
func parseDate(date string) (*time.Time, error) {
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

	return nil, fmt.Errorf("unable to parse date: %s", date)
}
