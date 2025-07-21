package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	config "github.com/jja42/gator/internal/config"
	"github.com/jja42/gator/internal/database"
)

type state struct {
	db  *database.Queries
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

	dbURL := configFile.DB_URL

	//Open a connection to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	//create our state and commands
	s := state{cfg: &configFile, db: dbQueries}

	c := commands{cmds: make(map[string]func(*state, command) error)}

	c.register("login", handlerLogin)

	c.register("register", handlerRegister)

	c.register("reset", handlerReset)

	c.register("users", handlerUsers)

	c.register("agg", handlerAgg)

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
