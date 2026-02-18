package main
import(
	"fmt"
	"context"
)

func Reset(s *state, cmd command) error{
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error resetting database: %w", err)
	}
	return nil
}