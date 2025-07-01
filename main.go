package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/lordbaldwin1/gator/internal/config"
	"github.com/lordbaldwin1/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("error: failed to open connection to database")
	}
	dbQueries := database.New(db)

	progState := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	cmds := commands{
		registry: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	if len(os.Args) < 2 {
		log.Fatal("Usage: <command> [args...]")
	}

	commandArgs := os.Args[2:]
	command := command{
		Name: os.Args[1],
		Args: commandArgs,
	}
	err = cmds.run(progState, command)
	if err != nil {
		log.Fatal(err)
	}
}
