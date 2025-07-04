package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lordbaldwin1/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return errors.New("error: feed not found")
	}

	createdFeed, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	printFeedFollow(createdFeed)
	return nil
}

func handlerFollowing(s *state, _ command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return errors.New("error: no feeds")
	}

	if len(feeds) == 0 {
		return errors.New("error: no feeds")
	}

	printFollowing(feeds)
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]
	_, err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    url,
	})
	if err != nil {
		return errors.New("error: failed to unfollow feed")
	}

	return nil
}

func printFollowing(feedFollows []database.GetFeedFollowsForUserRow) {
	fmt.Printf("Feeds followed for %s: \n", feedFollows[0].Name)

	for _, feed := range feedFollows {
		fmt.Println(" - ", feed.Name_2)
	}
}

func printFeedFollow(feedFollow database.CreateFeedFollowRow) {
	fmt.Println("Created Feed Follow:")
	fmt.Println("Feed Name: ", feedFollow.FeedName)
	fmt.Println("Username: ", feedFollow.UserName)
}
