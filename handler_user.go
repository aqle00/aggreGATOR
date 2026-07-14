package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aqle00/aggreGATOR/internal/database"
	"github.com/google/uuid"
)

// helper function to check if command was given
func getUsername(cmd Command) string {
	if len(cmd.args) != 1 {
		fmt.Printf("usage: %s <name>\n", cmd.name)
		os.Exit(1)
	}
	return cmd.args[0]
}

// helper function to check if user exists in database
// takes s *State, username string to query the db
// return true if user exists, false if not

//used like a type
// if userExists(s, <username>) {do stuff}

func userExists(s *State, username string) (bool, error) {
	_, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return false, err
}

// login command
func handlerLogin(s *State, cmd Command) error {
	username := getUsername(cmd)

	exists, err := userExists(s, username)
	if err != nil {
		return fmt.Errorf("failed to verify user: %v", err)
	}
	if !exists {
		return fmt.Errorf("User %s does not exist. Please register first.\n", username)
	}

	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("failed to set user: %v", err)
	}
	fmt.Printf("User set to: %s\n", s.cfg.CurrentUserName)
	return nil
}

// register new userthrough s.db
func handlerRegister(s *State, cmd Command) error {
	username := getUsername(cmd)

	exists, err := userExists(s, username)
	if err != nil {
		return fmt.Errorf("failed to verify user: %v", err)
	}
	if exists {
		return fmt.Errorf("User %s already exists. Please login instead.\n", username)
	}

	//make userParam strcut to use in CreateUser()
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	//call CreateUser()
	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	if err = s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("failed to set user: %v", err)
	}

	fmt.Println("User created!")
	fmt.Printf("ID: %s\n", user.ID)
	fmt.Printf("Created At: %v\n", user.CreatedAt)
	fmt.Printf("Updated At: %v\n", user.UpdatedAt)
	fmt.Printf("Name: %s\n", user.Name)
	//test block ends

	return nil
}

func handlerReset(s *State, cmd Command) error {
	return s.db.Reset(context.Background())
}

// speudo code here work on this asap
func handlerGetUsers(s *State, cmd Command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users: %v", err)
	}

	for _, user := range users {
		if user == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user)
			continue
		}
		fmt.Printf("* %s\n", user)
	}
	return nil
}
