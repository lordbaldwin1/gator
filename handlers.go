package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lordbaldwin1/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("error: username is required")
	}

	user, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return errors.New("error: user not found")
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to %s\n", s.cfg.CurrentUsername)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("error: username is required")
	}

	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}

	returnedUser, err := s.db.CreateUser(context.Background(), user)
	if err != nil {
		return errors.New("error: user already exists")
	}

	s.cfg.SetUser(returnedUser.Name)
	fmt.Printf("User has been set to %s\n", s.cfg.CurrentUsername)

	return nil
}

func handlerReset(s *state, _ command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return errors.New("error: failed to annihilate database")
	}
	fmt.Println("All user data deleted")
	return nil
}

func handlerGetUsers(s *state, _ command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return errors.New("error: failed to get users")
	}

	for _, user := range users {
		printString := fmt.Sprintf("* %s ", user.Name)
		if user.Name == s.cfg.CurrentUsername {
			printString += "(current)"
		}
		fmt.Println(printString)
	}
	return nil
}

func handlerAgg(s *state, _ command) error {
	const URL = "https://www.wagslane.dev/index.xml"

	feed, err := fetchFeed(context.Background(), URL)
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 {
		return errors.New("error: feed requires name and url")
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		return err
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

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
		return err //errors.New("error: failed to create feed follow in feed creation")
	}

	fmt.Println("Feed created successfully:")
	printFeed(returnedFeed)
	fmt.Println()
	fmt.Println("=====================================")

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

func handlerFeeds(s *state, _ command) error {
	feeds, err := s.db.GetFeedsWithUsername(context.Background())
	if err != nil {
		return err
	}

	printFeeds(feeds)

	return nil
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

func handlerFollow(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		return errors.New("error: current user not found")
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

func printFeedFollow(feedFollow database.CreateFeedFollowRow) {
	fmt.Println("Created Feed Follow:")
	fmt.Println("Feed Name: ", feedFollow.FeedName)
	fmt.Println("Username: ", feedFollow.UserName)
}

func handlerFollowing(s *state, _ command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		return errors.New("error: user not found")
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return errors.New("error: no feeds")
	}

	printFollowing(feeds)
	return nil
}

func printFollowing(feedFollows []database.GetFeedFollowsForUserRow) {
	fmt.Printf("Feeds followed for %s: \n", feedFollows[0].Name)

	for _, feed := range feedFollows {
		fmt.Println(" - ", feed.Name_2)
	}

}
