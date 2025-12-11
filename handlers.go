package main

import (
	"context"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"os"
	"time"

	"github.com/google/uuid"
)

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

func hanlderGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Println("error retrieving users from DB")
		return err
	}
	currentUser := s.config.CurrentUsername
	for _, user := range users {
		if user == currentUser {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	userName := s.config.CurrentUsername
	userData, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return err
	}
	name := cmd.args[0]
	url := cmd.args[1]
	new_feed := database.CreateFeedParams{
		Name:   name,
		Url:    url,
		UserID: userData.ID,
	}
	f, err := s.db.CreateFeed(context.Background(), new_feed)
	if err != nil {
		return err
	}
	fmt.Println(f)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	data, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range data {
		fmt.Println(feed)
	}
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
