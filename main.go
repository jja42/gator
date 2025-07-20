package main

import (
	"errors"
	"fmt"
	"os"

	config "github.com/jja42/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func main() {
	//Read our configFile
	configFile, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	//create our state and commands
	s := state{cfg: &configFile}

	c := commands{cmds: make(map[string]func(*state, command) error)}

	c.register("login", handlerLogin)

	//get command line arguments
	arguments := os.Args

	if len(arguments) < 2 {
		err = errors.New("command name required")
		fmt.Println(err)
		os.Exit(1)
	}

	//get command name and arguments
	command_name := arguments[1]
	params := make([]string, 0)
	if len(arguments) > 2 {
		params = arguments[2:]
	}

	//setup the command
	cmd := command{name: command_name, arguments: params}

	//run the command
	err = c.run(&s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
