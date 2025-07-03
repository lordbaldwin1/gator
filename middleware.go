package main

import (
	"context"

	"github.com/lordbaldwin1/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
		if err != nil {
			return err
		}

		err = handler(s, cmd, user)
		if err != nil {
			return err
		}
		return nil
	}
}
