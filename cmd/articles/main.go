package main

import (
	"context"
	"fmt"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/config"
	"github.com/mSulimenko/dev-blog-platform/internal/shared/database"
	"github.com/mSulimenko/dev-blog-platform/internal/shared/logger"
	"os"
)

func main() {
	fmt.Println("Starting articles")
	cfg := config.Load()

	log := logger.New(cfg.Env)
	defer log.Sync()

	// db
	dbpool, err := database.NewPool(context.Background(), cfg.DB.Dsn)
	if err != nil {
		log.Error("cannot connect to database: ", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	log.Info("Connected to database")

	err = database.RunMigrations(dbpool, cfg.DB.MigrationsDir)
	if err != nil {
		log.Error("migrations failed: ", err)
		os.Exit(1)
	}
	log.Info("Migrations applied successfully")

	// service
}
