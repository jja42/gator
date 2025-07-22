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

func handlerFollowing(s *state, cmd command) error {
	db := *s.db

	user, err := db.GetUser(context.Background(), s.cfg.UserName)
	if err != nil {
		return errors.New("current user could not be obtained")
	}

	follows, err := db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return errors.New("feed follows for user could not be obtained")
	}

	for _, follow := range follows {
		feed, err := db.GetFeed(context.Background(), follow.FeedID)
		if err != nil {
			return errors.New("feed could not be obtained from feed follow")
		}
		fmt.Printf("Name: %s\tURL: %s\n", feed.Name, feed.Url)
	}

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("missing name and url")
	}

	if len(cmd.arguments) == 1 {
		return errors.New("missing url")
	}

	//Connect to Database and Reset
	db := *s.db

	user, err := db.GetUser(context.Background(), s.cfg.UserName)
	if err != nil {
		return errors.New("current user could not be obtained")
	}

	//Setup Feed Params
	id := uuid.New()

	created_at := time.Now()

	updated_at := created_at

	name := cmd.arguments[0]

	url := cmd.arguments[1]

	args := database.CreateFeedParams{ID: id, CreatedAt: created_at, UpdatedAt: updated_at, Name: name, Url: url, UserID: user.ID}

	feed, err := db.CreateFeed(context.Background(), args)

	if err != nil {
		return errors.New("could not create feed")
	}

	fmt.Printf("Feed Successfully Created.\n Name: %s\tUrl: %s\tCreated At:%v\n", feed.Name, feed.Url, feed.CreatedAt)

	id = uuid.New()

	created_at = time.Now()

	updated_at = created_at

	follow_args := database.CreateFeedFollowParams{ID: id, CreatedAt: created_at, UpdatedAt: updated_at, UserID: user.ID, FeedID: feed.ID}

	_, err = db.CreateFeedFollow(context.Background(), follow_args)
	if err != nil {
		return errors.New("could not follow feed")
	}

	fmt.Printf("Feed Successfully Followed.\n")

	return nil
}

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
	//Connect to Database
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

func handlerFeeds(s *state, cmd command) error {
	//Connect to Database and Reset
	database := *s.db

	feeds, err := database.GetFeeds(context.Background())
	if err != nil {
		return errors.New("could not fetch feeds from database")
	}

	for _, feed := range feeds {
		fmt.Printf("Name: %s\tURL: %s\n", feed.Name, feed.Url)
		user, err := database.GetUserFromID(context.Background(), feed.UserID)
		if err != nil {
			return errors.New("unable to get associated user for field")
		}
		fmt.Printf("Created By: %s\n", user.Name)
		fmt.Printf("\n")
	}

	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("missing url argument")
	}

	url := cmd.arguments[0]

	feed, err := s.db.GetFeedFromURL(context.Background(), url)
	if err != nil {
		return errors.New("unable to get feed to follow")
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.UserName)
	if err != nil {
		return errors.New("unable to get user info")
	}

	//Setup Feed Follow Params
	id := uuid.New()

	created_at := time.Now()

	updated_at := created_at

	args := database.CreateFeedFollowParams{ID: id, CreatedAt: created_at, UpdatedAt: updated_at, UserID: user.ID, FeedID: feed.ID}

	follow, err := s.db.CreateFeedFollow(context.Background(), args)
	if err != nil {
		return errors.New("unable to create feed follow")
	}

	fmt.Printf("Feed Follow Created Successfully.\n")

	fmt.Printf("Feed Follow ID: %v\tUser: %s\tFeed: %s\n", follow.ID, follow.UserName, follow.FeedName)

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
	db := *s.db

	user, err := db.CreateUser(context.Background(), args)
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
