package main
import(
	"fmt"
	"context"
	"time"
	"strings"
	"log"
	"database/sql"
	"github.com/google/uuid"
	"github.com/MrBushido-002/Gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1{
		return fmt.Errorf("Not enough or too many args")
	}
	duration, err := time.ParseDuration(cmd.Args[0])
	if err != nil{
		return fmt.Errorf("could not get duration: %w", err)
	}
	ticker := time.NewTicker(duration)

	fmt.Printf("Collecting feeds every %v\n", duration)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			fmt.Printf("error scraping feeds: %v\n", err)
		}
	}
	return nil
}



func handlerListFeeds(s *state, cmd command) error{
	feeds, err := s.db.GetFeeds(context.Background()) 
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Couldn't get user: %v", err)
		}
		fmt.Printf("Name: %s \nURL: %s\nUser: %s\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}

func scrapeFeeds(s *state) error {
	next, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil{
		return fmt.Errorf("could not get next feed: %w", err)
	}
	err = s.db.MarkFeedFetched(context.Background(), next.ID)
	if err != nil{
		return fmt.Errorf("could not mark feed fetched: %w", err)
	}
	rssFeed, err := fetchFeed(context.Background(), next.Url)
	if err != nil{
		return fmt.Errorf("could not fetch feed by url: %w", err)
	}
	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("* %s\n", item.Title)

		pubAt := sql.NullTime{}

		description := sql.NullString{
			String: item.Description,
			Valid:  true, 
		}

		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			pubAt.Time = t
			pubAt.Valid = true
		}


		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: description,
			PublishedAt: pubAt,
			FeedID:		 next.ID,
			},
		)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	return nil
}