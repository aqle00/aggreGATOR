package main

import (
	"context"
	"fmt"

	"github.com/aqle00/aggreGATOR/internal/database"
)

func middlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		username := s.cfg.CurrentUserName
		user, err := s.db.GetUser(context.Background(), username)
		if err != nil {
			return fmt.Errorf("User not logged in\n%v", err)
		}
		return handler(s, cmd, user)
	}
}
