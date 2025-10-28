package main

import (
	"context"
	"fmt"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/config"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/repository"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/service"
	grpcclient "github.com/mSulimenko/dev-blog-platform/internal/articles/transport/grpc"
	httphandler "github.com/mSulimenko/dev-blog-platform/internal/articles/transport/http"
	"github.com/mSulimenko/dev-blog-platform/internal/shared/database"
	"github.com/mSulimenko/dev-blog-platform/internal/shared/logger"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Starting articles")
	cfg := config.Load()
	fmt.Println(cfg)
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
	articlesRepo := repository.NewArticlesRepository(dbpool)

	// service
	articleService := service.NewArticlesService(log, articlesRepo)

	// grpc client
	grpcAuthClient, err := grpcclient.NewAuthClient(context.Background(),
		log,
		cfg.GRPC.Addr,
		cfg.GRPC.RetryTimeout,
		cfg.GRPC.MaxRetries,
	)
	if err != nil {
		log.Error("failed to initialise grpcAuthClient: ", err)
		os.Exit(1)
	}

	// router
	handler := httphandler.NewHandler(articleService, log, grpcAuthClient)
	router := handler.InitRouter()

	srv := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	log.Infof("Starting server on %s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	if err = srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %w", err)
	}

}
