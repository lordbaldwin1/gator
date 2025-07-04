package main

import "fmt"

func handlerHelp(s *state, cmd command) error {
	fmt.Println("Available commands:")
	fmt.Println()

	fmt.Println("login: login to existing user")
	fmt.Println("usage: login [username]")

	fmt.Println("register: create new user")
	fmt.Println("usage: register [username]")

	fmt.Println("reset: reset all data")
	fmt.Println("usage: reset")

	fmt.Println("users: lists users")
	fmt.Println("usage: users")

	fmt.Println("agg: start aggregating posts for feeds")
	fmt.Println("usage: agg")

	fmt.Println("addfeed: adds a feed via url to aggregate posts from")
	fmt.Println("usage: addfeed [username] [feed url]")

	fmt.Println("feeds: lists all feeds")
	fmt.Println("usage: feeds")

	fmt.Println("follow: follow a feed via url")
	fmt.Println("usage: follow [feed url]")

	fmt.Println("following: lists feeds you are currently following")
	fmt.Println("usage: following")

	fmt.Println("unfollow: unfollows a feed via url")
	fmt.Println("usage: unfollow [feed url]")

	fmt.Println("browse: look at posts")
	fmt.Println("usage: browse [optional # of posts]")
	return nil
}
