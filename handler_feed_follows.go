package main

import (
	"fmt"
	"github.com/google/uuid"
	"time"
	"context"
	"github.com/MrBushido-002/Gator/internal/database"
)

func follow(s *state, cmd command, user database.User) error {
	if len(cmd.Args)  != 1 {
		return fmt.Errorf("not enough Args")
	}

	url := cmd.Args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("could not get feed: %v", err)
	}

	feed_follows, err := s.db.CreateFeedFollow(context.Background(),database.CreateFeedFollowParams{
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		UserID:		user.ID,
		FeedID:		feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("could not create feed follows: %v", err)
	}
	fmt.Printf("Name: %v\nUser: %v", feed_follows.FeedName, feed_follows.UserName)
	return nil
}

func handlerListFeedFollows(s *state, cmd command, user database.User) error{
	feed_follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not get feed follows: %v", err)
	}

	for _, feed_follow := range feed_follows {
		fmt.Printf("Feed: %v\nUser: %v", feed_follow.FeedName, user.Name)
	}
	return nil
}

func handlerFeedUnfollow(s *state, cmd command, user database.User) error{
	if len(cmd.Args)  != 1 {
		return fmt.Errorf("not enough Args")
	}
	url := cmd.Args[0]
	err := s.db.Unfollow(context.Background(), database.UnfollowParams{
		UserID: user.ID,
		Url: url,
		},
	)
	if err != nil {
		return fmt.Errorf("could not unfollow: %w", err)
	}
	return nil
}
