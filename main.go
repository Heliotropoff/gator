package main

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"os"
	"time"

	"github.com/google/uuid"
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

type state struct {
	db     *database.Queries
	config *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no username was provided")
	}
	userNameProvided := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), userNameProvided)
	if err != nil {
		errMsg := fmt.Errorf("problem retrieving user by the name %s, error: %w", userNameProvided, err)
		return errMsg
	}
	if err := s.config.SetUser(user.Name); err != nil {
		return err
	}
	fmt.Println(user.Name, " was set as current user name")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no username was provided")
	}
	uName := cmd.args[0]
	DbArgs := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      uName,
	}
	usr, err := s.db.CreateUser(context.Background(), DbArgs)
	if err != nil {
		errMSG := fmt.Errorf("db error for adding user %s. error: %w", DbArgs.Name, err)
		fmt.Println(errMSG)
		os.Exit(1)
	}
	s.config.SetUser(usr.Name)
	fmt.Println("New user was created")
	fmt.Print(usr)
	return nil

}

func hanlderReset(s *state, cmd command) error {
	if err := s.db.DeleteUsers(context.Background()); err != nil {
		fmt.Println("reset was not successful")
		return err
	}
	fmt.Println("reset was successful")
	return nil
}

type commands struct {
	supported map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if fn, ok := c.supported[cmd.name]; !ok {
		return fmt.Errorf("command %s not supported", cmd.name)
	} else {
		if err := fn(s, cmd); err != nil {
			return err
		}
	}
	return nil

}

func (c *commands) register(name string, f func(*state, command) error) {
	c.supported[name] = f
}
