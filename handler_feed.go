package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/aqle00/aggreGATOR/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

func handlerAddFeed(s *State, cmd Command, user database.User) error {
	// check if name and url was provided
	if len(cmd.args) != 2 {
		return fmt.Errorf(`usage: %s "<name>" "<url>"\n`, cmd.name)
	}

	// stuff to use later when creating a feed
	name := cmd.args[0]
	url := cmd.args[1]

	//make userParam strcut to use in CreateUser()
	feedParams := database.CreateFeedParams{
		ID:            uuid.New(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		LastFetchedAt: sql.NullTime{},
		Name:          name,
		Url:           url,
		UserID:        user.ID,
	}

	//call CreateFeed()
	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" && pqErr.Constraint == "feeds_url_key" {
				return fmt.Errorf("feed url already exists: %v", err)
			}
		}
		return fmt.Errorf("Couldnt create feed: %v", err)
	}
	fmt.Printf("Created feed: %+v\n", feed)
	//--------automatically follow created feed----------

	autoFollowArg := []string{feed.Url}
	autoFollowCmd := Command{
		name: "follow",
		args: autoFollowArg,
	}

	err = handlerFeedFollow(s, autoFollowCmd, user)
	if err != nil {
		return err
	}

	//------------------------------------------------------------------------------------------

	fmt.Println("Feed created and followed!")
	fmt.Printf("ID: %s\n", feed.ID)
	fmt.Printf("Created At: %v\n", feed.CreatedAt)
	fmt.Printf("Updated At: %v\n", feed.UpdatedAt)
	fmt.Printf("Name: %s\n", feed.Name)
	fmt.Printf("Url: %s\n", feed.Url)
	fmt.Printf("UserID: %s\n", feed.UserID)

	return nil
}
