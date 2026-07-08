package main

import (
	"context"
	"fmt"

	"os"
	"time"

	"github.com/aqle00/aggreGATOR/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *State, cmd Command) error {
	//connect this to CreateFeed()

	// check if name and url was provided
	if len(cmd.args) != 2 {
		fmt.Printf(`usage: %s "<name>" "<url>"\n`, cmd.name)
		os.Exit(1)
	}

	// stuff to use later when creating a feed
	username := s.cfg.CurrentUserName
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}

	user_id := user.ID

	name := cmd.args[0]
	url := cmd.args[1]

	// check if url is unique( url exists in the feed or no )

	//make userParam strcut to use in CreateUser()
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user_id,
	}

	//call CreateFeed()
	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	fmt.Println("Feed created!")
	fmt.Printf("ID: %s\n", feed.ID)
	fmt.Printf("Created At: %v\n", feed.CreatedAt)
	fmt.Printf("Updated At: %v\n", feed.UpdatedAt)
	fmt.Printf("Name: %s\n", feed.Name)
	fmt.Printf("Url: %s\n", feed.Url)
	fmt.Printf("UserID: %s\n", feed.UserID)

	return nil
}
