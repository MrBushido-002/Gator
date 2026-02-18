package main
import(
	"fmt"
	"log"
	"os"
	"database/sql"

	"github.com/MrBushido-002/Gator/internal/database"
	"github.com/MrBushido-002/Gator/internal/config"
)

import _ "github.com/lib/pq"

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error opening database")
	}

	dbQueries := database.New(db)


	programState := &state{
		cfg: &cfg,
		db: dbQueries,
	}

	cmds := commands{
    	handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", Reset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerRegisterFeed))
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", middlewareLoggedIn(follow))
	cmds.register("following", middlewareLoggedIn(handlerListFeedFollows))
	cmds.register("unfollow", middlewareLoggedIn(handlerFeedUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "not enough args")
		os.Exit(1)
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := command{Name: cmdName, Args: cmdArgs}

	err = cmds.run(programState, cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}