package main

import (
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Print(fmt.Errorf("reading config file failed with %w", err))
		return

	}
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		newErr := fmt.Errorf("error while connecting to DB: %w", err)
		fmt.Println(newErr)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	currentState := state{
		db:     dbQueries,
		config: &cfg,
	}

	currentCommands := commands{
		supported: make(map[string]func(*state, command) error),
	}
	currentCommands.register("login", handlerLogin)
	currentCommands.register("register", handlerRegister)
	currentCommands.register("reset", hanlderReset)
	currentCommands.register("users", hanlderGetUsers)
	currentCommands.register("agg", handlerAgg)
	currentCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	currentCommands.register("feeds", handlerFeeds)
	currentCommands.register("follow", middlewareLoggedIn(handlerFollow))
	currentCommands.register("following", middlewareLoggedIn(handlerFollowing))
	currentCommands.register("unfollow", handlerUnfollow)
	currentCommands.register("browse", handlerBrowse)
	passed_args := os.Args
	if len(passed_args) < 2 {
		err := fmt.Errorf("no arguments were provided")
		fmt.Println(err)
		os.Exit(1)
	}
	passed_command := command{
		name: passed_args[1],
		args: passed_args[2:],
	}
	if err = currentCommands.run(&currentState, passed_command); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
