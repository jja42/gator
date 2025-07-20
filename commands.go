package main

import (
	"errors"
	"fmt"

	"github.com/jja42/gator/internal/config"
)

func handlerLogin(s *state, cmd command) error {

	//Check for Empty Arguments
	if len(cmd.arguments) == 0 {
		return errors.New("missing username argument")
	}

	//Get username from arguments
	username := cmd.arguments[0]

	//Set username
	err := config.SetUser(username, *s.cfg)
	if err != nil {
		return err
	}

	fmt.Println("User has been Logged In and Set")

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
