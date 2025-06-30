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
