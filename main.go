package main

import (
	"fmt"
	"gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Print(fmt.Errorf("reading config file failed with %w", err))
		return

	}
	if err = cfg.SetUser("greg"); err != nil {
		fmt.Print(fmt.Errorf("writing config file failed with %w", err))
		return
	}
	cfg, err = config.Read()
	if err != nil {
		fmt.Print(fmt.Errorf("reading config file failed with %w", err))
		return
	}
	fmt.Println(cfg)

}
