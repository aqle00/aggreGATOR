package main

import (
	"context"
	"fmt"

	"os"
	"time"

	"github.com/aqle00/aggreGATOR/internal/database"
	"github.com/google/uuid"
)

// prints all feeds to the console
func handlerListFeeds(s *State, cmd Command) error {

	// get feeds []string from the db
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get feeds: %v", err)
	}

	//for loop
	for i := range feeds {
		username := getUsernameByID(s, feeds[i].UserID)
		fmt.Printf("Name: %s", feeds[i].Name)
		fmt.Printf("URL: %s", feeds[i].Url)
		fmt.Printf("Created by: %s", username)
	}
	return nil
}

// helper for handlerListFeeds()
func getUsernameByID(s *State, id uuid.UUID) string {
	user, err := s.db.GetUserByID(context.Background(), id)
	if err != nil {
		fmt.Errorf("failed to get username: %v", err)
		return ""
	}
	username := user.Name
	return username
}

func handlerAddFeed(s *State, cmd Command) error {
	username := s.cfg.CurrentUserName
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}

	// check if name and url was provided
	if len(cmd.args) != 2 {
		fmt.Printf(`usage: %s "<name>" "<url>"\n`, cmd.name)
		os.Exit(1)
	}

	// stuff to use later when creating a feed
	user_id := user.ID
	name := cmd.args[0]
	url := cmd.args[1]

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
		return fmt.Errorf("Couldnt create feed: %v", err)
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
