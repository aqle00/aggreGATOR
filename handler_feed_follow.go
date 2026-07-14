package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aqle00/aggreGATOR/internal/database"
	"github.com/lib/pq"
)

func handlerFeedFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.args) != 1 {
		fmt.Printf(`usage: %s "<url>"\n`, cmd.name)
		os.Exit(1)
	}
	url := cmd.args[0]

	// check url by get feed with url
	fmt.Printf("Looking up feed with URL: %q\n", url)
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Could not fetch feed with url %s\n %v", url, err)
	}

	//make params
	feedParams := database.CreateFeedFollowParams{
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	//call CreateFeedFollow() with params
	follow, err := s.db.CreateFeedFollow(context.Background(), feedParams)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" && pqErr.Constraint == "feeds_follows_pkey" {
				return fmt.Errorf("Feed already followed, cannot follow again")
			}
		}
		return fmt.Errorf("Could not follow: %v", err)
	}

	fmt.Printf("Followed %s\n", follow.FeedName)
	fmt.Printf("Current user: %s\n", follow.UserName)

	return nil
}

func handlerFollowingFeeds(s *State, cmd Command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Error getting followed feeds: %v", err)
	}

	for _, feed := range feeds {
		fmt.Println(feed.FeedName)
	}
	return nil
}

func handlerFeedUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.args) != 1 {
		fmt.Printf(`usage: %s "<url>"\n`, cmd.name)
		os.Exit(1)
	}
	url := cmd.args[0]

	// check url to see if feed exists
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Could not fetch feed with url %s\n %v", url, err)
	}

	deleteFollowParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.DeleteFeedFollow(context.Background(), deleteFollowParams)
	if err != nil {
		return fmt.Errorf("Failed to unfollow %s", feed.Name)
	}

	fmt.Printf("Unfollowed %s\n", feed.Name)
	fmt.Printf("Current user: %s\n", user.Name)
	return nil
}
