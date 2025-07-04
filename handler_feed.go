package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lordbaldwin1/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	fmt.Println(name, url)

	returnedFeed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        uuid.New(),
		Name:      name,
		Url:       url,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    returnedFeed.UserID,
		FeedID:    returnedFeed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Println("Feed created successfully:")
	printFeed(returnedFeed)
	fmt.Println()
	fmt.Println("=====================================")

	return nil
}

func handlerFeeds(s *state, _ command) error {
	feeds, err := s.db.GetFeedsWithUsername(context.Background())
	if err != nil {
		return err
	}

	printFeeds(feeds)

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}

func printFeeds(feeds []database.GetFeedsWithUsernameRow) {
	fmt.Println("List of Feeds:")
	fmt.Println()

	for _, feed := range feeds {
		fmt.Println("Feed Name: ", feed.Name)
		fmt.Println("URL: ", feed.Url)
		fmt.Println("Username: ", feed.Name_2)
		fmt.Println()
	}
}
