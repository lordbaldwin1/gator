package main

import (
	"context"
	"errors"
	"fmt"
)

func handlerReset(s *state, _ command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return errors.New("error: failed to annihilate database")
	}
	fmt.Println("All user data deleted")
	return nil
}
