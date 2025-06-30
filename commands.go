package main

import (
	"errors"
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	registry map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if s == nil {
		return errors.New("error: state is nil")
	}

	cmdToRun, ok := c.registry[cmd.Name]
	if !ok {
		return errors.New("error: command does not exist")
	}

	err := cmdToRun(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	_, ok := c.registry[name]
	if ok {
		fmt.Println("error: command already exists")
	}

	c.registry[name] = f
}
