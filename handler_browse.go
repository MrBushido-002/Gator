package main

import(
	"fmt"
	"context"
	"strconv"
	"github.com/MrBushido-002/Gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error{
	limit := 3
	if len(cmd.Args) != 0 {
		val, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = val
	}


	params := database.GetPostForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
		}

	 posts, err := s.db.GetPostForUser(context.Background(), params)
	 if err != nil {
		return fmt.Errorf("could not get posts: %w", err)
	 }

	for _, post := range posts {
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("Published: %v\n", post.PublishedAt.Time.Format("Mon Jan 2, 2006"))
		fmt.Printf("Link:      %s\n", post.Url)
		if post.Description.Valid {
			fmt.Printf("\n%s\n", post.Description.String)
		}
		fmt.Println("========================================")
	}
	 return nil 
}