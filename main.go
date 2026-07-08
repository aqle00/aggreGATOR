package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/aqle00/aggreGATOR/internal/database"

	"github.com/aqle00/aggreGATOR/internal/config"
	_ "github.com/lib/pq"
)

// State struct to hold config
type State struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	// call config.Read() from internal/config/config.go package
	// to read the config file
	// returns a config.Config struct and error

	cfg, err := config.Read()
	if err != nil {
		panic(fmt.Errorf("failed to read config: %v", err))
	}

	// open database connection using cfg.DBURL
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		panic(fmt.Errorf("failed to open database: %v", err))
	}
	defer db.Close()

	dbQueries := database.New(db)

	programState := &State{
		cfg: &cfg,
		db:  dbQueries,
	}

	cliCommands := Commands{
		registeredCommands: make(map[string]func(*State, Command) error),
	}

	// register commands with their handlers
	cliCommands.register("login", handlerLogin)
	cliCommands.register("register", handlerRegister)
	cliCommands.register("reset", handlerReset)
	cliCommands.register("users", handlerGetUsers)
	cliCommands.register("agg", handlerAgg)
	cliCommands.register("addfeed", handlerAddFeed)

	// if < 2 that means there's only 1 argument
	// which is defaulted to the program name
	// os.Args[0] = program name
	// which means no Command was provided
	raw := os.Args[1:]
	if len(raw) < 1 {
		fmt.Println("No Command provided")
		os.Exit(1)
	}

	userInput := Command{
		name: raw[0],
		args: raw[1:],
	}

	if err := cliCommands.run(programState, userInput); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

}
