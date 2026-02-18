package main

import (
	"fmt"
	"github.com/google/uuid"
	"time"
	"context"
	"github.com/MrBushido-002/Gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Name not in database")
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Println("User switched successfully")
	return nil
}



func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args)  != 1 {
		return fmt.Errorf("not enough Args")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		},	
	)
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}
	
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldnt set username: %s", err)
	}

	fmt.Printf("Registration complete User: %v\n ID: %v\n", user.Name, user.ID)

	return nil
}

func handlerListUsers(s *state, cmd command) error{
	current := s.cfg.CurrentUserName

	users, err := s.db.GetUsers(context.Background()) 
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name != current {
			fmt.Printf("* %s\n", user.Name)
		} else {
			fmt.Printf("* %s (current)\n", user.Name)
		}
	}
	return nil
}