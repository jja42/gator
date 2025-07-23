package main

import (
	"context"
	"errors"

	"github.com/jja42/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.UserName)
		if err != nil {
			return errors.New("unable to get user info")
		}

		return handler(s, cmd, user)
	}
}
