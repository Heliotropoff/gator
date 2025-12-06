package main

import (
	"fmt"
	"gator/internal/config"
	"os"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Print(fmt.Errorf("reading config file failed with %w", err))
		return

	}
	currentState := state{
		Config: &cfg,
	}
	currentCommands := commands{
		supported: make(map[string]func(*state, command) error),
	}
	currentCommands.register("login", handlerLogin)
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
	Config *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no username was provided")
	}
	userName := cmd.args[0]
	if err := s.Config.SetUser(userName); err != nil {
		return err
	}
	fmt.Println(userName, " was set as current user name")
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
