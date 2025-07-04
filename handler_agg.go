package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lordbaldwin1/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %s <inverval> (e.g., 1h, 1m, 1g)", cmd.Name)
	}

	interval, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return errors.New("error: failed to parse interval")
	}

	ticker := time.NewTicker(interval)
	log.Printf("Collecting feeds every %s...", interval)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}
	log.Println("Found a feed to fetch!")

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
		ID:        nextFeed.ID,
	})
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", nextFeed.Name, err)
		return
	}

	// FETCH FEED USING URL, ITERATE AND PRINT TITLES TO CONSOLE
	feedData, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", nextFeed.Name, err)
		return
	}

	// print titles
	for _, item := range feedData.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			UpdatedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			Title: item.Title,
			Url:   item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: publishedAt,
			FeedID:      nextFeed.ID,
		})
		if err != nil {
			log.Printf("Failed to create post: %s", err)
			continue
		}
	}

	log.Printf("Feed %s collected, %v posts found", nextFeed.Name, len(feedData.Channel.Item))
}
