package main

import (
	"fmt"
)

// define what a command is
type Command struct {
	name string
	args []string
}

// holds a map that maps command names to their handler functions
type Commands struct {
	registeredCommands map[string]func(*State, Command) error
}

// run command executes a command by calling the function with the State if it exists in the handlerCommands map. If the command does not exist, it returns an error.
func (c *Commands) run(s *State, cmd Command) error {
	handler, ok := c.registeredCommands[cmd.name]
	if !ok {
		return fmt.Errorf("command not found: %s", cmd.name)
	}
	return handler(s, cmd)

}

// register a new handler command for a command name
func (c *Commands) register(name string, f func(*State, Command) error) error {
	c.registeredCommands[name] = f
	return nil
}
