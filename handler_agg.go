package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aqle00/aggreGATOR/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerAgg(s *State, cmd Command) error {
	if len(cmd.args) != 1 {
		fmt.Printf(`usage: %s "<time>"\n`, cmd.name)
		fmt.Printf(`Ex: %s 1m\n`, cmd.name)
		fmt.Printf(`"m" can be replaced with "h", "ms", "s"\n`)

		return fmt.Errorf("Syntax error")
	}

	time_between_reqs := cmd.args[0]
	reqCD, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		return fmt.Errorf("Failed to parse given time, check syntax: %v", err)
	}

	ticker := time.NewTicker(reqCD)
	fmt.Printf("Fetching feed every %v\n", time_between_reqs)
	for ; ; <-ticker.C {
		fmt.Println("start scrape loop reached-----")
		err := scrapeFeeds(s)
		if err != nil {
			return fmt.Errorf("Error scraping feeds: %v", err)
		}
	}

	return nil
}

func scrapeFeeds(s *State) error {
	//flow
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to fetch next feed info: %v", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return fmt.Errorf("Failed to mark feed as fetched: %v", err)
	}

	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("Failed to fetch feed: %v", err)
	}

	// //print down
	fmt.Printf("aggregated feed: %v\n", feed.Channel.Title)
	fmt.Printf(`feed: %v`, feed.Channel.Item)

	for _, item := range feed.Channel.Item {
		//save posts in db instead of printing them
		layout := "Mon, 02 Jan 2006 15:04:05 -0700"
		published_time, err := time.Parse(layout, item.PubDate)
		if err != nil {
			return fmt.Errorf("Failed to parse time: %v", err)
		}
		postParams := database.CreatePostsParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: published_time,
			FeedID:      nextFeed.ID,
		}

		post, err := s.db.CreatePosts(context.Background(), postParams)
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) {
				if pqErr.Code == "23505" && pqErr.Constraint == "posts_url_key" {
					fmt.Printf("Post url already exists")
					continue
				}
			}
			return fmt.Errorf("Couldnt create post: %v", err)
		}
		fmt.Printf("id: %s\n", post.ID)
		fmt.Printf("title: %s\n", post.Title)
		fmt.Printf("description: %s\n", post.Description.String)
		fmt.Printf("Published date: %s\n", post.PublishedAt.UTC())

	}
	return nil
}

func handlerBrowse(s *State, cmd Command, user database.User) error {
	var numberOfPosts int32 = 2
	if len(cmd.args) > 0 {
		num, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf(`Usage: %s "<name>" "<number of posts to show. Ex: 5>"\n`, cmd.name)
		}
		numberOfPosts = int32(num)
	}

	getPostParams := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  numberOfPosts,
	}

	posts, err := s.db.GetPostsForUser(context.Background(), getPostParams)
	if err != nil {
		return fmt.Errorf("Failed to get posts: %v", err)
	}

	for _, post := range posts {
		fmt.Printf("id: %s\n", post.ID)
		fmt.Printf("title: %s\n", post.Title)
		fmt.Printf("description: %s\n", post.Description.String)
		fmt.Printf("Published date: %s\n", post.PublishedAt.UTC())
	}

	return nil
}
