package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jja42/gator/internal/config"
	"github.com/jja42/gator/internal/database"
)

func handlerReset(s *state, cmd command) error {
	//Connect to Database and Reset
	database := *s.db

	err := database.DeleteUsers(context.Background())
	if err != nil {
		return errors.New("users table could not be deleted")
	}

	println("Database Was Successfully Reset.")

	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("%v", feed)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	//Connect to Database and Reset
	database := *s.db

	users, err := database.GetUsers(context.Background())
	if err != nil {
		return errors.New("could not fetch users from database")
	}

	for _, user := range users {
		fmt.Printf("%s", user.Name)
		if user.Name == s.cfg.UserName {
			fmt.Printf(" (current)")
		}
		fmt.Printf("\n")
	}

	return nil
}

func handlerLogin(s *state, cmd command) error {

	//Check for Empty Arguments
	if len(cmd.arguments) == 0 {
		return errors.New("missing username argument")
	}

	//Get username from arguments
	username := cmd.arguments[0]

	//Connect to Database and Check for User
	database := *s.db

	_, err := database.GetUser(context.Background(), username)
	if err != nil {
		return errors.New("user could not be retrieved from database")
	}

	//Set username
	err = config.SetUser(username, *s.cfg)
	if err != nil {
		return errors.New("user could not be set")
	}

	fmt.Println("User has been Logged In and Set")

	return nil
}

func handlerRegister(s *state, cmd command) error {

	//Check for Empty Arguments
	if len(cmd.arguments) == 0 {
		return errors.New("missing name argument")
	}

	//Get username from arguments
	username := cmd.arguments[0]

	//Set Up User Args
	id := uuid.New()

	created_at := time.Now()

	updated_at := created_at

	args := database.CreateUserParams{Name: username, ID: id, CreatedAt: created_at, UpdatedAt: updated_at}

	//Connect to Database and Create User
	database := *s.db

	user, err := database.CreateUser(context.Background(), args)
	if err != nil {
		return errors.New("user could not be created")
	}

	//Set our User
	config.SetUser(username, *s.cfg)

	fmt.Println("User was successfully created.")

	fmt.Printf("User Name: %s\tUser ID: %v\tCreated At: %v\n", user.Name, user.ID, user.CreatedAt)

	return nil
}

func (c *commands) run(s *state, cmd command) error {
	//check if command exists
	if v, ok := c.cmds[cmd.name]; ok {
		//run command
		err := v(s, cmd)
		//return error if present
		if err != nil {
			return err
		}
	} else {
		return errors.New("requested command does not exist in commands")
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}
