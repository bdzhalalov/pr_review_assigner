package main

import (
	"github.com/bdzhalalov/pr-review-assigner/config"
	"github.com/bdzhalalov/pr-review-assigner/pkg/database"
	"github.com/bdzhalalov/pr-review-assigner/pkg/logger"
)

func main() {
	cfg := config.InitConfig()

	log := logger.Logger(&cfg)

	db, err := database.ConnectToDB(&cfg, log)

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error closing database connection: %v", err)
		}
	}(db)
}
