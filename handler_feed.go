package main

import (
	"fmt"
	"github.com/google/uuid"
	"time"
	"context"
	"github.com/MrBushido-002/Gator/internal/database"
)

func handlerRegisterFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args)  != 2 {
		return fmt.Errorf("not enough Args")
	}

	feed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:   user.ID,
		},	
	)
	if err != nil {
		return fmt.Errorf("couldn't add feed: %w", err)
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("error following feed: %v", err)
	}
	
	fmt.Printf("%+v", feed)
	return nil
}
