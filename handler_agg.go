package main

import (
	"context"
	"fmt"
)

// test functino to fetch 1 specific feed

// will edit later to become main aggregator function
func handlerAgg(s *State, cmd Command) error {
	feedURL := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	fmt.Printf("%v", feed)
	return nil
}
