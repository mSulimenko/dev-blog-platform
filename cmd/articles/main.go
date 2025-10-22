package main

import (
	"context"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/config"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/database"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/repository"
	service2 "github.com/mSulimenko/dev-blog-platform/internal/articles/service"
	httphandler "github.com/mSulimenko/dev-blog-platform/internal/articles/transport/http"
	"github.com/mSulimenko/dev-blog-platform/internal/shared/logger"
	"net/http"
	"os"
)

func main() {

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

	// repo
	usersRepo := repository.NewUsersRepository(dbpool)

	// services
	userService := service2.NewUsersService(usersRepo, log)

	// router
	handler := httphandler.NewHandler(userService, log)
	router := handler.InitRouter()

	srv := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	log.Infof("Starting server on %s:%s", cfg.HTTP.Host, cfg.HTTP.Port)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
